package main

import (
	"github.com/antandros/go-yaps/example/plugin_cmd/myplugin"
	"github.com/antandros/go-yaps/manager"
)

func main() {
	cnf := &manager.PluginManagerConfig{
		KillUnattendedClients: true,
		ManagerId:             "base",
		BinaryPath:            "./bin/",
		CreateBinary:          true,
		CreateInterface:       true,
		InterfacePath:         "plugins",
	}
	pmanager := manager.RegisterManager(cnf)
	prs := &myplugin.TestItem{}
	pmanager.RegisterPlugin(&manager.PluginConfig{
		Language: manager.Go,
		Impl:     prs,
	})
	pmanager.CreateInterfaces()

}
