package main

import (
	"context"
	"io"
	"log"
	"lualoader/internal/config"
	"lualoader/internal/golua"
	"lualoader/internal/plugins"
	"lualoader/internal/utils"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

var (
	cfg  config.Config
	elog *log.Logger
)

const (
	VERSION_MAJOR uint = 1
	VERSION_MINOR uint = 0
	VERSION_PATCH uint = 0

	PLUGINS_PATH string = "plugins"
)

func init() {
	log.Printf("lualoader v%d.%d.%d\n", VERSION_MAJOR, VERSION_MINOR, VERSION_PATCH)
	log.Println("Init...")

	eLogFile, err := os.Create("errlog.txt")
	if err != nil {
		log.Fatalln("Init: cannot create log file: ", err)
	}
	elog = log.New(io.MultiWriter(eLogFile, os.Stdout), "[ERR] ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	err = config.CheckConfig()
	if err != nil {
		elog.Fatalln("Init: CheckConfig: ", err)
	}
	err = config.ReadConfigTo(&cfg)
	if err != nil {
		elog.Fatalln("Init: ReadConfigTo: ", err)
	}

	err = utils.CheckDir(PLUGINS_PATH)
	if err != nil {
		elog.Fatalln("Init: CheckDir: ", err)
	}

	// 注册路由
	golua.Mux = http.NewServeMux()

	// Enable Plugins
	err = plugins.EnablePlugins(PLUGINS_PATH)
	if err != nil {
		elog.Fatalln(err)
	}

	log.Println("Init Done.")
}

// waitgroup
var wg sync.WaitGroup

// 终止信号用channel
var (
	signalListener = make(chan os.Signal)
	stopSignal     = make(chan struct{})
	httpDone       = make(chan struct{})
)

// stop信号广播
func sl() {
	<-signalListener
	close(signalListener)
	close(stopSignal)
	log.Println("Stopping...")
	wg.Done()
}

func main() {
	// 监听终止信号
	signal.Notify(signalListener, os.Interrupt)
	wg.Add(3)
	go sl()

	var srv = new(http.Server)
	srv.Addr = cfg.Port
	srv.Handler = golua.Mux
	go func() {
		<-stopSignal
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			elog.Printf("HTTP server Shutdown: %v\n", err)
		}
		log.Println("HTTP server Shutdown.")
		close(httpDone)
		wg.Done()
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		elog.Fatalf("HTTP server ListenAndServe: %v\n", err)
	}

	go func() {
		<-stopSignal
		<-httpDone
		err := plugins.DisablePlugins()
		if err != nil {
			elog.Fatalln(err)
		}
		wg.Done()
	}()

	wg.Wait()

	log.Println("Stopped.")
}
