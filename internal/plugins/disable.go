package plugins

import (
	"errors"
	"log"
	"lualoader/internal/golua"
)

func DisablePlugins() error {
	for k, v := range golua.LuaStatePool {
		log.Println("Disabling " + k)

		err := golua.LuaPluginRunDisableFunc(v)
		if err != nil {
			return errors.New("failed to run function \"disable()\"")
		}

		golua.LuaCloseState(v)

		delete(golua.LuaStatePool, k)

		log.Println("Disabled " + k)
	}
	return nil
}
