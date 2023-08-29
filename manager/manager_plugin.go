package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/antandros/go-yaps/parser"
	"github.com/antandros/go-yaps/protocol"
	"github.com/antandros/go-yaps/yaperror"
	"github.com/fatih/structs"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/connectivity"
)

type Plugin struct {
	name           string
	logfile        string
	socket         string
	addr           string
	token          string
	port           int64
	isRemote       bool
	bin            bool
	configResponse bool

	language   Language
	execOut    bytes.Buffer
	execErr    bytes.Buffer
	logger     *zap.Logger
	plg        PluginInterface
	server     protocol.Server
	client     *protocol.Client
	manager    *PluginManager
	process    *os.Process
	exec       *exec.Cmd
	StructData parser.Struct
	StructMap  map[string]interface{}
}

func (p *Plugin) KillProcess() {
	if p.process != nil {
		p.process.Kill()
	}
}
func (p *Plugin) SetManager(manager *PluginManager) {

	p.manager = manager
}
func (p *Plugin) Call(structName string, function string, args []interface{}) (interface{}, error) {

	p.GetLogger().Info("Function call from plugin client", zap.String("struct", structName), zap.Any("args", args), zap.String("function", function))
	isError, err := p.ValidateFunction(structName, function, args)
	if !isError {
		p.GetLogger().Error("Validation error function", zap.Error(err), zap.Any("err_data", err), zap.String("struct", structName), zap.Any("args", args), zap.String("function", function))
		return []interface{}{}, err
	}
	params := p.FunctionParamPop(structName, function, args)

	return p.client.Call(function, structName, params)
}
func (p *Plugin) FunctionParamPop(structName, function string, args []interface{}) []protocol.InItem {
	var inItems []protocol.InItem
	item := p.StructMap["Item"].(map[string]interface{})
	if !strings.EqualFold(item["Name"].(string), structName) {
		item = nil
		refItems := p.StructMap["RelatedItems"].([]interface{})
		for _, refIt := range refItems {
			xitem := refIt.(map[string]interface{})
			if strings.EqualFold(xitem["Name"].(string), structName) {
				item = xitem
			}
		}
	}
	for _, fnc := range item["Methods"].([]interface{}) {
		method := fnc.(map[string]interface{})
		methodName := method["Name"].(string)
		if strings.EqualFold(methodName, function) {
			inParams := method["InParams"].([]interface{})
			for _, inParam := range inParams {
				param := inParam.(map[string]interface{})
				itemIndex := param["Index"].(float64) - 1

				inItems = append(inItems, protocol.InItem{
					Index:    int(param["Index"].(float64)),
					Type:     param["Type"].(string),
					BaseData: args[int(itemIndex)],
				})
			}
		}
	}
	return inItems
}
func (p *Plugin) ValidateStruct(item map[string]interface{}, function string, args []interface{}) []map[string]string {
	var err []map[string]string
	var found bool
	for _, fnc := range item["Methods"].([]interface{}) {
		method := fnc.(map[string]interface{})
		methodName := method["Name"].(string)
		if strings.EqualFold(methodName, function) {
			found = true
			inParams := method["InParams"].([]interface{})
			if len(args) != len(inParams) {
				errText := map[string]string{
					"paramLength": fmt.Sprint(len(inParams)),
					"argLenth":    fmt.Sprint(len(inParams)),
					"error":       "param length error",
				}
				return []map[string]string{
					errText,
				}
			}
			var errText []string
			for _, inParam := range inParams {
				param := inParam.(map[string]interface{})
				paramStruct := param["IsStruct"].(bool)
				paramSlice := param["IsSlice"].(bool)
				paramPtr := param["IsPtr"].(bool)
				itemIndex := param["Index"].(float64) - 1
				if len(args) <= int(itemIndex) {
					tdata := fmt.Sprintf("param %f indexed is not found", itemIndex)
					errText = append(errText, tdata)
					continue
				}
				argData := args[int(param["Index"].(float64)-1)]
				argTypeOf := reflect.TypeOf(argData)
				isPtr := false
				isSlice := false
				if argTypeOf.Kind() == reflect.Slice {
					isSlice = true
					argTypeOf = argTypeOf.Elem()
				}

				if argTypeOf.Kind() == reflect.Ptr {
					argTypeOf = argTypeOf.Elem()
					isPtr = true
				}
				if paramPtr != isPtr {
					if paramPtr {
						errText = append(errText, "param is ptr your value is not ptr ")
					} else {
						errText = append(errText, "value is ptr function param is not ptr")
					}
				}
				if paramSlice != isSlice {
					if paramSlice {
						errText = append(errText, "param is slice your value is not slice ")
					} else {
						errText = append(errText, "value is slice function param is not slice")
					}
				}

				if paramStruct != (argTypeOf.Kind() == reflect.Struct) {
					if paramStruct {
						errText = append(errText, "param is struct your value is not struct ")
					} else {
						errText = append(errText, "value is struct function param is not struct")
					}
				}
				if len(errText) > 0 {
					err = append(err, map[string]string{
						"error":    "validation",
						"argIndex": fmt.Sprint(itemIndex),
						"detail":   strings.Join(errText, ","),
					})
				}

			}
		}
	}
	if found {
		return err
	}
	errorText := map[string]string{"error": "function not found"}
	return []map[string]string{
		errorText,
	}
}
func (p *Plugin) ValidateFunction(structName string, function string, args []interface{}) (bool, error) {
	itemName := p.StructMap["Item"].(map[string]interface{})
	if strings.EqualFold(itemName["Name"].(string), structName) {
		enErr := p.ValidateStruct(itemName, function, args)
		if len(enErr) > 0 {
			err := yaperror.Error(yaperror.VALIDATE_ITEM, nil, yaperror.WithExra(map[string]interface{}{
				"errors": enErr,
			}))
			p.GetLogger().Error("Validation error struct", zap.Error(err), zap.Any("err_data", err))
			return false, err
		}
		return true, nil
	}
	refItems := p.StructMap["RelatedItems"].([]interface{})
	for _, refIt := range refItems {
		item := refIt.(map[string]interface{})
		if strings.EqualFold(item["Name"].(string), structName) {
			enErr := p.ValidateStruct(item, function, args)
			if len(enErr) > 0 {
				err := yaperror.Error(yaperror.VALIDATE_ITEM, nil, yaperror.WithExra(map[string]interface{}{
					"errors": enErr,
				}))
				p.GetLogger().Error("Valodation error related items", zap.Error(err))
				return false, err
			}
			return true, nil
		}
	}

	return false, yaperror.Error(yaperror.VALIDATE_ITEM, nil, yaperror.WithExra(map[string]interface{}{
		"baseStruct": itemName,
		"callStruct": structName,
	}))
}
func (p *Plugin) Connected() bool {
	if p.client == nil {
		fmt.Println("Client not init")
		return false
	}
	return p.client.ConnectionStatus() == connectivity.Ready
}
func (p *Plugin) ClientStatus() string {
	if p.client == nil {
		return "CLIENT_NOT_INIT"
	}
	if !p.configResponse {
		return "config_wait"
	}
	return p.client.ConnectionStatus().String()
}
func (p *Plugin) PluginCaller(fu string, str string, ins []*protocol.InTypes) any {
	p.GetLogger().Info("Request caller", zap.Any("out", fu))
	if p.manager != nil {
		return p.manager.StructCall(fu, str, ins)
	} else {
		p.GetLogger().Info("Manager Not Set")
		panic("manager cant found")
	}
}
func (p *Plugin) CreateToken() string {
	p.token = uuid.NewString()
	return p.token
}
func (p *Plugin) ProcessStatus() {
	fmt.Println("ERR", p.exec.Err)
	fmt.Println("ProcessState", p.exec.ProcessState)
	fmt.Println("Process", p.exec.Process)
	fmt.Println("Pid", p.exec.Process.Pid)
	fmt.Println(p.execOut.String())
	fmt.Println(p.client.ConnectionStatus().String())
	fmt.Println("=======================")
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
}
func (p *Plugin) CreateProcess() error {
	if p.manager == nil {
		return yaperror.Error(yaperror.MANAGER_NOT_FOUND, nil)
	}
	binItems, err := os.ReadDir(p.manager.config.BinaryPath)
	if err != nil {
		return err
	}
	var itemPath string
	for _, binItem := range binItems {
		if strings.EqualFold(binItem.Name(), p.name) {
			itemPath = binItem.Name()
		}
	}
	if itemPath == "" {
		return yaperror.Error(yaperror.PLUGIN_BINARY_NOT_FOUND, nil)
	}
	itemPath = path.Join("./", p.manager.config.BinaryPath, itemPath)
	server := []string{"-server", p.manager.socket}
	token := []string{"-token", p.CreateToken()}
	register := []string{"-register", p.Socket()}
	cmdArgs := append(server, token...)
	cmdArgs = append(cmdArgs, register...)
	p.exec = exec.Command(itemPath, cmdArgs...)
	p.addr = p.socket
	p.exec.Stdout = &p.execOut
	p.exec.Stderr = &p.execErr
	err = p.exec.Start()
	fmt.Println("RUNERR", err)
	fmt.Println(p.exec)
	if err != nil {
		return yaperror.Error(yaperror.RUN_BINARY, err, yaperror.WithExra(map[string]interface{}{
			"output": p.execOut.String(),
			"err":    p.execErr.String(),
		}))
	}
	p.process = p.exec.Process

	return nil
}
func (p *Plugin) CreateClient() error {
	err := p.CreateProcess()
	if err != nil {
		return err
	}
	time.Sleep(time.Second)

	go p.ConnectClient()

	return nil
}
func (p *Plugin) ConnectClient() {
	p.client = protocol.NewClient(p.addr, p.isRemote, context.Background(), p.GetLogger())
	err := p.client.Connect()
	if err != nil {
		panic(err)
	}
	p.client.WaitConnect()
	data, err := p.client.GetConfig()
	if err != nil {
		fmt.Println("error", err)
		fmt.Println("error", err)
		fmt.Println("error", err)
		fmt.Println("error", err)
		fmt.Println("error", err)
		fmt.Println("error", err)
		fmt.Println("error", err)
		fmt.Println("error", err)
	}
	p.configResponse = true
	var intf map[string]interface{}
	json.Unmarshal(data, &intf)
	p.StructMap = intf
}
func (p *Plugin) ParseStruct() {
	parser := parser.NewStructParser(p.plg)
	p.StructData = parser
	structItem := structs.New(parser)
	p.StructMap = structItem.Map()
}
func (p *Plugin) GetLogger() *zap.Logger {
	if p.logger == nil {
		p.logger = p.manager.GetPluginLogger(p.name)
	}
	return p.logger
}

func (p *Plugin) GenConfig() map[string]interface{} {

	return p.StructMap
}
func (p *Plugin) Serve() error {
	if p.language != Go {
		panic("not implemented language")
	}

	p.logger = p.manager.GetPluginLogger(p.name)

	var err error
	fmt.Println("p.socket", p.socket)
	p.socket = p.Socket()
	var uri string
	var lType string
	if p.isRemote {
		uri = fmt.Sprintf("%s:%d", p.addr, p.port)
		lType = "tcp"
	} else {
		uri = p.socket
		lType = "unix"
	}
	p.server, err = protocol.NewServer(lType, uri, p)
	if err != nil {
		p.GetLogger().Fatal("cant create a protocol server", zap.Error(err), zap.String("sokcet_uri", p.socket), zap.String("name", p.Name()))
	}

	p.server.ConfigFunction = p.GenConfig
	return p.server.Run()
}
func (p *Plugin) Name() string {
	return p.name
}
func (p *Plugin) SetSocket(sck string) {
	p.socket = sck
}
func (p *Plugin) Socket() string {
	if p.socket == "" {
		pth, err := os.MkdirTemp("/tmp/", p.Name())
		if err != nil {
			p.GetLogger().Fatal("cant create a tempdir", zap.Error(err))
		}

		manId := uuid.New().String()

		pId := fmt.Sprintf("%s.socket", manId)
		socket := filepath.Join(pth, pId)
		p.socket = socket

		//defer os.RemoveAll(pth)
	}
	return p.socket
}
func (p *Plugin) Language() Language {
	return p.language
}

type PluginConfig struct {
	Multiple     bool
	Language     Language
	Binary       bool
	File         string
	Name         string
	RemotePlugin bool
	Addr         string
	Port         int64
	Impl         PluginInterface
}
