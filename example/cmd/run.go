package main

import (
	"fmt"
	"time"

	"github.com/antandros/go-yaps/example/plugins"
	"github.com/antandros/go-yaps/manager"
)

func main() {
	cnf := &manager.PluginManagerConfig{
		ManagerId:  "base",
		BinaryPath: "./bin/",
	}
	pmanager := manager.RegisterManager(cnf)
	err := pmanager.Create()
	if err != nil {
		panic(err)
	}
	fmt.Println("creating item")
	item := plugins.NewTestItem()
	fmt.Println("init item")
	resp := item.AddSlice("test")
	fmt.Println(resp)
	items, err := item.GetSlice()
	fmt.Println(items, err)
	item.AddSlice("test 3")
	item.AddSlice("test 4")
	item.AddSlice("test 5")
	ir := 0
	for {
		items, err = item.GetSlice()

		fmt.Println("-------------------------------------")
		fmt.Println(items, err)

		for i := range pmanager.Plugins {
			fmt.Println(pmanager.Plugins[i].Name(), pmanager.Plugins[i].Connected(), pmanager.Plugins[i].Connecting())
			if ir == 3 {
				fmt.Println("Test reconnect")
				pmanager.Plugins[i].Reconnect()
			}
			ir++
			if pmanager.Plugins[i].Connected() == false && pmanager.Plugins[i].Connecting() == false {
				pmanager.Plugins[i].Reconnect()
				<-time.After(time.Second * 10)
				if pmanager.Plugins[i].Connected() {
					item.AddSlice("test 3")
					item.AddSlice("test 4")
				}
			}
		}
		<-time.After(time.Second * 3)
	}

}
