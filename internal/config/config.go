package config

import (
	"encoding/json"
	"lualoader/internal/utils"
	"os"
)

type Config struct {
	Port string `json:"port"`
}

func CheckConfig() error {
	exi, err := utils.Exists("config.json")
	if err != nil {
		return err
	}
	if !exi {
		f, err := os.Create("config.json")
		if err != nil {
			return err
		}
		defer f.Close()

		// 写入默认配置
		var defaultCfg Config
		defaultCfg.Port = ":19130"
		bts, err := json.MarshalIndent(defaultCfg, "", "    ")
		if err != nil {
			return err
		}
		_, err = f.Write(bts)
		if err != nil {
			return err
		}
	}
	return nil
}

// 加载配置文件
func ReadConfigTo(c *Config) error {
	b, err := os.ReadFile("config.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, c)
	if err != nil {
		return err
	}
	return nil
}
