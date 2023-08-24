package main

import (
	"github.com/antandros/go-yaps/example/plugin_cmd/myplugin"
	"github.com/antandros/go-yaps/manager"
	"github.com/antandros/go-yaps/plugin"
)

func main() {

	cnf := &manager.PluginManagerConfig{
		LogFile: "logs/plgrunner.log",
	}
	plg := new(myplugin.TestItem)
	plugin.RegisterPlugin(plg, cnf)
}
