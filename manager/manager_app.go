package manager

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/antandros/go-yaps/helper"
	grmon "github.com/bcicen/grmon/agent"
	"github.com/google/uuid"
)

func (pm *PluginManager) registerPlugin(pconfig *PluginConfig) *Plugin {
	fmt.Println(pconfig.File)
	var plg Plugin
	if !pconfig.Binary {
		plg = Plugin{
			name:     pconfig.Impl.Name(),
			plg:      pconfig.Impl,
			socket:   pconfig.File,
			addr:     pconfig.Addr,
			port:     pconfig.Port,
			isRemote: pconfig.RemotePlugin,
			language: pconfig.Language,
		}
		plg.ParseStruct()
	} else {
		for i := range pm.Plugins {
			plugn := pm.Plugins[i]
			if strings.EqualFold(plugn.Name(), pconfig.Name) {
				return nil
			}

		}

		plg = Plugin{
			name:     pconfig.Name,
			socket:   pconfig.File,
			bin:      pconfig.Binary,
			addr:     pconfig.Addr,
			port:     pconfig.Port,
			isRemote: pconfig.RemotePlugin,
		}
		plg.SetManager(pm)
		err := plg.CreateClient()
		if err != nil {
			panic(err)
		}
		for {
			if strings.EqualFold(plg.ClientStatus(), "ready") {
				break
			} else {
				fmt.Println(plg.ClientStatus())
			}
			time.Sleep(100 * time.Millisecond)
		}

	}

	return &plg
}

var managerB *PluginManager

func GetManager() (*PluginManager, error) {
	if managerB == nil {
		return nil, errors.New("please register manager before call items")
	}
	return managerB, nil
}
func RegisterManager(cnf *PluginManagerConfig) *PluginManager {
	grmon.Start()

	fmt.Println("RegisterManager")
	pth, err := os.MkdirTemp("", "pluginManager")
	fmt.Println("TEMP", pth, err)
	if err != nil {
		panic(err)
	}
	var manId string
	if cnf != nil {
		manId = cnf.ManagerId
	} else {
		manId = uuid.New().String()
	}
	if manId == "" {
		manId = uuid.New().String()
	}
	logger := helper.Logger()
	if cnf.LogFile == "" {
		cnf.LogFile = fmt.Sprintf("./logs/manager/%s.log", manId)
		logger = helper.LoggerNamed(cnf.LogFile)
	}
	pthLog := path.Dir(cnf.LogFile)
	_, pthErr := os.Stat(pthLog)
	if pthErr != nil {
		errMk := os.MkdirAll(pthLog, 0700)
		fmt.Println("MkdirAll", pth, errMk)
		if errMk != nil {
			panic(errMk)
		}
	}

	pId := fmt.Sprintf("%s.sokcet", manId)
	fmt.Println("pId", pId, logger)

	defer os.RemoveAll(pth)
	socket := filepath.Join(pth, "manager", pId)

	managerB = &PluginManager{
		socket: socket,
		logger: logger,
		config: cnf,
	}

	return managerB
}
