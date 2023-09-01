package main

import (
	"github.com/antandros/go-yaps/example/plugins"
	"github.com/antandros/go-yaps/manager"
)

func main() {
	cnf := &manager.PluginManagerConfig{
		KillUnattendedClients: true,
		ManagerId:             "base",
		BinaryPath:            "./bin/",
	}
	pmanager := manager.RegisterManager(cnf)
	err := pmanager.Create()
	if err != nil {
		panic(err)
	}

	_ = plugins.NewTestItem()

}
