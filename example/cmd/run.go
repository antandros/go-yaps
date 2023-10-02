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
	for {
		fmt.Println(time.Now().Format(time.RFC1123))
		items, err = item.GetSlice()
		fmt.Println(items, err)
		<-time.After(time.Second * 5)
	}

}
