package plugins

import (
	"errors"
	"log"
	"lualoader/internal/golua"
	"os"
	"path"
	"runtime"
	"strings"
)

func EnablePlugins(pluginsPath string) error {
	var wd, err = os.Getwd()
	if err != nil {
		return err
	}
	var ppath = path.Join(wd, pluginsPath)

	dir, err := os.ReadDir(ppath)
	if err != nil {
		return err
	}

	golua.InitLuaStatePool()

	for _, fi := range dir {
		if fi.IsDir() {
			info, err := fi.Info()
			if err != nil {
				continue
			}
			err = enablePlugin(path.Join(ppath, info.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func enablePlugin(pluginPath string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(pluginPath)
	defer os.Chdir(currentDir)
	if err != nil {
		return err
	}

	var manifest Manifest
	err = readManifestTo(pluginPath, &manifest)
	if err != nil {
		return err
	}

	log.Println("Enabling plugin " + manifest.Name)

	if manifest.ManifestVersion != CURRENT_MANIFEST_VERSION {
		return errors.New("bad manifest version")
	}

	if len(manifest.PluginVersion) != 3 {
		return errors.New("bad plugin version")
	}

	flag := false
	if len(manifest.Platform.Os) != 0 {
		for _, os := range manifest.Platform.Os {
			flag = (os == runtime.GOOS)
			if flag {
				break
			}
		}
		if !flag {
			return errors.New("unsupported os")
		}
	}
	if len(manifest.Platform.Arch) != 0 {
		for _, arch := range manifest.Platform.Arch {
			flag = (arch == runtime.GOARCH)
			if flag {
				break
			}
		}
		if !flag {
			return errors.New("unsupported arch")
		}
	}

	log.Printf("Author: %s; v%d.%d.%d\n", strings.Join(manifest.Author, ", "), manifest.PluginVersion[0], manifest.PluginVersion[1], manifest.PluginVersion[2])

	if _, ok := golua.LuaStatePool[manifest.Name]; ok {
		return errors.New("find another plugin with the same name: " + manifest.Name)
	}
	golua.LuaStatePool[manifest.Name] = golua.LuaNewState()
	err = golua.LuaDoFile(golua.LuaStatePool[manifest.Name], path.Join(pluginPath, manifest.Entry))
	if err != nil {
		return err
	}

	err = golua.LuaPluginRunEnableFunc(golua.LuaStatePool[manifest.Name])
	if err != nil {
		return err
	}

	err = golua.LuaPluginSetHandler(golua.LuaStatePool[manifest.Name])
	if err != nil {
		return err
	}

	log.Println("Enabled plugin " + manifest.Name)

	return nil
}
