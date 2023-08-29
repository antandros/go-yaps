package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/antandros/go-yaps/helper"
	"github.com/antandros/go-yaps/parser"
	"github.com/antandros/go-yaps/protocol"
	"github.com/antandros/go-yaps/yaperror"
	"github.com/modern-go/reflect2"
	"github.com/noirbizarre/gonja"
	"go.uber.org/zap"
)

type Language uint

const (
	Go Language = iota
	Python
	Java
	JavaScript
)

type PluginManagerConfig struct {
	KillUnattendedClients bool
	ManagerId             string
	BinaryPath            string
	CreateBinary          bool
	LogFile               string
	CreateInterface       bool
	InterfacePath         string
}

type PluginManager struct {
	config       *PluginManagerConfig
	socket       string
	Plugins      []Plugin
	Plugin       Plugin
	logger       *zap.Logger
	plogger      map[string]*zap.Logger
	PluginConfig *PluginConfig
	mainmanager  bool
}

func (pm *PluginManager) Defer() {
	if err := recover(); err != nil {
		pm.logger.Error("Error recovery", zap.Any("error", err), zap.Any("Stack", debug.Stack()))
	}
}
func (pm *PluginManager) StructCall(fu string, str string, ins []*protocol.InTypes) any {
	defer pm.Defer()
	unItems := make([]reflect.Value, len(ins))
	pm.logger.Info("Request StructCall", zap.Any("out", fu), zap.Any("ins", len(ins)))
	outItems := ins

	for _, itm := range outItems {
		nType := itm.GetType()
		refType := reflect2.TypeByName(nType)
		if refType == nil {
			refType = reflect2.Type2(reflect.TypeOf(""))
		}
		var itemVal interface{}
		json.Unmarshal(itm.GetIn(), &itemVal)
		valn := reflect.ValueOf(itemVal)
		valn = valn.Convert(refType.Type1())
		unItems[itm.Index-1] = valn

	}

	method := reflect.ValueOf(pm.PluginConfig.Impl).MethodByName(fu)
	valn := []reflect.Value{}
	if len(ins) > 0 {
		valn = unItems
	}
	resp := method.Call(valn)
	var respData []interface{}
	if len(resp) > 0 {
		for _, d := range resp {
			respData = append(respData, d.Interface())
		}
		pm.logger.Info("response", zap.Any("data", respData))
		return respData
	}
	return resp[0].Interface()
}
func (pm *PluginManager) GetPluginLogger(plugin string) *zap.Logger {
	if pm.PluginConfig != nil {
		return pm.logger
	}
	if pm.plogger == nil {
		pm.plogger = make(map[string]*zap.Logger)
	}
	if logger, ok := pm.plogger[plugin]; ok {
		return logger
	}
	logPath := fmt.Sprintf("logs/plugin/%s.log", plugin)
	for _, plg := range pm.Plugins {
		if strings.EqualFold(plg.Name(), plugin) {
			if plg.logfile != "" {
				logPath = plg.logfile
			}
		}
	}
	pm.plogger[plugin] = helper.LoggerNamed(logPath)
	return pm.plogger[plugin]
}
func (pm *PluginManager) GetErrorLogger() *zap.Logger {
	return pm.logger
}
func (pm *PluginManager) RegisterPlugin(pconfig *PluginConfig) {

	pm.logger.Info("Register plugin", zap.String("name", pconfig.Name))
	plugin := pm.registerPlugin(pconfig)
	if plugin == nil {
		panic("plugin already registered")
	}
	pm.Plugins = append(pm.Plugins, *plugin)
}

func (pm *PluginManager) RunPlugin(pconfig *PluginConfig) {

	pm.logger.Info("Register plugin", zap.String("name", pconfig.Impl.Name()), zap.Any("pconfig", pconfig))
	pm.Plugin = *pm.registerPlugin(pconfig)
	pm.PluginConfig = pconfig
	pm.Plugin.SetManager(pm)
	pm.Plugin.Serve()
}

func (pm *PluginManager) CreateDirectory(dir string) {
	_, err := os.Stat(dir)
	if err != nil {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			pm.logger.Fatal("cant create path", zap.String("path", dir))
		}
	}
}
func (pm *PluginManager) Wait() {
	time.Sleep(time.Second * 4)
	for {

		for i := range pm.Plugins {
			pm.Plugins[i].ProcessStatus()
		}
		time.Sleep(time.Second)
	}

}
func (pm *PluginManager) Create() error {
	for i := range pm.Plugins {
		pm.Plugins[i].SetManager(pm)
		err := pm.Plugins[i].CreateClient()
		if err != nil {
			panic(err)
		}
	}
	return nil
}
func (pm *PluginManager) CreateBinary() error {
	pm.logger.Info("Generate binaries")
	goPath := path.Join(runtime.GOROOT(), "bin")
	if len(pm.Plugins) == 0 {
		return yaperror.Error(yaperror.ZERO_REGISTERED_PLUGIN, nil)
	}
	binPath, _ := filepath.Abs(pm.config.BinaryPath)
	for _, plugin := range pm.Plugins {
		cmdArgs := append([]string{}, "build")
		binFile := path.Join(binPath, strings.ToLower(plugin.name))
		fmt.Println(binFile)
		cmdArgs = append(cmdArgs, "-o")
		cmdArgs = append(cmdArgs, binFile)
		cmdArgs = append(cmdArgs, fmt.Sprintf("%s.go", binFile))

		build := exec.Command("go", cmdArgs...)
		fmt.Println(cmdArgs)
		build.Env = helper.ReplaceGoPath(os.Environ(), goPath)
		output, err := build.CombinedOutput()
		fmt.Println(string(output))
		fmt.Println(output, err)
	}

	return nil
}
func (pm *PluginManager) Generate() {
	pm.mainmanager = true
	pm.logger.Info("Generate start", zap.Any("config", pm.config))
	if pm.config.CreateBinary {
		if pm.config.BinaryPath == "" {
			pm.config.BinaryPath = "./bin/"
		}
		pm.CreateDirectory(pm.config.BinaryPath)
		pm.CreateBinary()
	}
	if pm.config.CreateInterface {
		pm.logger.Info("Creating a interfaces")
		if pm.config.InterfacePath == "" {
			pm.config.InterfacePath = "./plugins/"
		}
		pm.CreateDirectory(pm.config.InterfacePath)

		pm.CreateInterfaces()
	}
}
func (pm *PluginManager) CreateInterfaceFile(pname string, packagePath string, item *parser.StructItem, allitem []*parser.StructItem, pkgName string, name string, inititem bool) error {
	tpl, err := gonja.FromString(GetTemplate())
	if err != nil {
		return err
	}
	out, err := tpl.Execute(gonja.Context{"item": item, "package_name": pkgName, "inititem": inititem, "allitem": allitem, "plugName": pname})
	if err != nil {
		return err
	}
	tempName := path.Join(packagePath, fmt.Sprintf("%s.go", strings.ToLower(name)))
	fmt.Println(tempName)
	f, err := os.Create(tempName)
	if err != nil {

		return err
	}
	defer f.Close()
	f.Write([]byte(out))
	f.Close()
	return nil
}
func (pm *PluginManager) Error(errType yaperror.YapsError, err error, options ...yaperror.Options) error {

	errs := yaperror.Error(errType, err, options...)
	errs.ZapError(pm.logger)
	return errs
}
func (pm *PluginManager) CreateInterfaces() error {
	if len(pm.Plugins) == 0 {
		return yaperror.Error(yaperror.ZERO_REGISTERED_PLUGIN, nil)
	}
	pathSlice := strings.Split(pm.config.InterfacePath, "/")
	packagePath := pathSlice[len(pathSlice)-1]
	if packagePath == "" {
		packagePath = pathSlice[len(pathSlice)-2]
	}
	for _, plugin := range pm.Plugins {
		item := plugin.StructData.Item
		fmt.Println("item", item.Name)

		pkgNameSlice := strings.Split(packagePath, "/")
		pkgName := pkgNameSlice[len(pkgNameSlice)-1]
		if pkgName == "" {
			pkgName = pkgNameSlice[len(pkgNameSlice)-2]
		}
		err := pm.CreateInterfaceFile(plugin.name, packagePath, item, plugin.StructData.RelatedItems, pkgName, item.Name, true)
		if err != nil {
			errorDesc := map[string]interface{}{
				"plugin":  plugin.name,
				"pkgName": pkgName,
				"item":    item,
			}
			pm.Error(yaperror.INTERFACE_CREATE, err, yaperror.WithExra(errorDesc))
		}
		for _, rItem := range plugin.StructData.RelatedItems {
			err = pm.CreateInterfaceFile(plugin.name, packagePath, rItem, plugin.StructData.RelatedItems, pkgName, fmt.Sprintf("%s_%s", item.Name, rItem.Name), false)
			if err != nil {
				errorDesc := map[string]interface{}{
					"plugin":  plugin.name,
					"pkgName": pkgName,
					"item":    rItem,
				}
				pm.Error(yaperror.INTERFACE_CREATE, err, yaperror.WithExra(errorDesc))
			}
		}

	}
	return nil
}
func (pm *PluginManager) CallFunction(pluginName string, strucName string, function string, args []interface{}) (interface{}, error) {
	fmt.Println("Call FNC", pluginName, strucName, function)
	for i, plg := range pm.Plugins {
		if strings.EqualFold(plg.name, pluginName) {
			fmt.Println("Plugin found", pluginName, strucName, function)
			if !pm.Plugins[i].Connected() {
				err := yaperror.Error(yaperror.NOT_CONNECTED, nil, yaperror.WithExra(map[string]interface{}{
					"plugin":   pluginName,
					"struct":   strucName,
					"function": function,
					"params":   args,
					"status":   pm.Plugins[i].ClientStatus(),
				}))
				pm.logger.Error("plugin not found", zap.Error(err))
				return nil, err
			}
			pm.GetPluginLogger(pluginName).Info("Plugin connected and send call", zap.Any("pluginName", pluginName), zap.Any("args", args), zap.Any("strucName", strucName), zap.Any("function", function))

			response, err := pm.Plugins[i].Call(strucName, function, args)
			pm.GetPluginLogger(pluginName).Info("Plugin connected and send call", zap.Any("err", err), zap.Any("response", response), zap.Any("pluginName", pluginName), zap.Any("args", args), zap.Any("strucName", strucName), zap.Any("function", function))
			if response == nil {
				return []interface{}{}, err
			}
			return response, err

		}
	}
	return nil, nil
}
