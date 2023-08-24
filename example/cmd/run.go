package main

import (
	"fmt"

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

	item := plugins.NewTestItem()
	fmt.Println(err)
	fmt.Println(item.AddSlice("test"))
	fmt.Println(item.AddSlice("test"))
	fmt.Println(item.AddSlice("test"))
	fmt.Println(item.AddSlice("test"))
	item.AddSlice("test")
	item.AddSlice("test")
	fmt.Println(item.GetSlice())
	fmt.Println(item.Name())

}
