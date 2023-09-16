package plugin

import (
	"flag"
	"fmt"
	"runtime/debug"

	"github.com/antandros/go-yaps/manager"
)

func Version() string {
	info, _ := debug.ReadBuildInfo()

	return fmt.Sprintf("-%s", info.GoVersion)
}
func RegisterPlugin(plg manager.PluginInterface, cnf *manager.PluginManagerConfig) {
	var serverURI string
	flag.StringVar(&serverURI, "server", "", "a string var")

	var registerURI string
	flag.StringVar(&registerURI, "register", "", "register unix file")

	var token string
	flag.StringVar(&token, "token", "", "service register token")

	var remoteAddress string
	flag.StringVar(&remoteAddress, "listen", "", "remote listen address")

	var remotePort int64
	flag.Int64Var(&remotePort, "port", 9991, "remote liten port ")

	flag.Parse()
	pmanager := manager.RegisterManager(cnf)
	pmanager.RunPlugin(&manager.PluginConfig{
		Language:     manager.Go,
		Impl:         plg,
		IsPlugin:     true,
		RemotePlugin: remoteAddress != "",
		Addr:         remoteAddress,
		Port:         remotePort,
		File:         registerURI,
	})
}
