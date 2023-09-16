package main

import (
	"fmt"
	"time"

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
	resp := item.AddSlice("test")
	fmt.Println(resp)
	items, err := item.GetSlice()
	fmt.Println(items, err)
	item.AddSlice("test 3")
	item.AddSlice("test 4")
	item.AddSlice("test 5")

	<-time.After(time.Second)
	items, err = item.GetSlice()
	fmt.Println(items, err)

}
