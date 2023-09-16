package plugins

import (
	"errors"
	"fmt"
	"time"

	"github.com/antandros/go-yaps/manager"
	"go.uber.org/zap"
)

type TestItem struct {
	DenemeItem      *DenemeItem
	DenemeItemSlice []DenemeItem
	Zaman           time.Time
	Caller          func(string, string, string, []interface{}) (interface{}, error)
	TestItemLogger  *zap.Logger
}

func (pi *TestItem) AddSlice(param1 string) error {
	if pi.Caller == nil {
		return errors.New("caller not init")
	}
	var params []interface{}
	params = append(params, param1)
	_, err := pi.Caller("Test", "TestItem", "AddSlice", params)
	if err != nil {
		if pi.TestItemLogger != nil {
			pi.TestItemLogger.Error("function response error", zap.Error(err), zap.String("function", "AddSlice"), zap.String("interface", "TestItem"))
		}
		return err
	}
	return nil
}

func (pi *TestItem) CreateDn(param1 string) error {
	if pi.Caller == nil {
		return errors.New("caller not init")
	}
	var params []interface{}
	params = append(params, param1)
	_, err := pi.Caller("Test", "TestItem", "CreateDn", params)
	if err != nil {
		if pi.TestItemLogger != nil {
			pi.TestItemLogger.Error("function response error", zap.Error(err), zap.String("function", "CreateDn"), zap.String("interface", "TestItem"))
		}
		return err
	}
	return nil
}

func (pi *TestItem) CreateDnX(param1 string) (string, error) {
	if pi.Caller == nil {
		return "", errors.New("caller not init")
	}
	var params []interface{}
	params = append(params, param1)
	response, err := pi.Caller("Test", "TestItem", "CreateDnX", params)
	if err != nil {
		if pi.TestItemLogger != nil {
			pi.TestItemLogger.Error("function response error", zap.Error(err), zap.String("function", "CreateDnX"), zap.String("interface", "TestItem"))
		}
		return "", err
	}
	pReturn := response.([]interface{})
	var pReturn_1 string
	if mItem, ok := pReturn[0].(string); ok {
		pReturn_1 = mItem
	}
	return pReturn_1, nil
}

func (pi *TestItem) CreateDnXS(param1 string) (string, string, error) {
	if pi.Caller == nil {
		return "", "", errors.New("caller not init")
	}
	var params []interface{}
	params = append(params, param1)
	response, err := pi.Caller("Test", "TestItem", "CreateDnXS", params)
	if err != nil {
		if pi.TestItemLogger != nil {
			pi.TestItemLogger.Error("function response error", zap.Error(err), zap.String("function", "CreateDnXS"), zap.String("interface", "TestItem"))
		}
		return "", "", err
	}
	pReturn := response.([]interface{})
	var pReturn_1 string
	if mItem, ok := pReturn[0].(string); ok {
		pReturn_1 = mItem
	}
	var pReturn_2 string
	if mItem, ok := pReturn[1].(string); ok {
		pReturn_2 = mItem
	}
	return pReturn_1, pReturn_2, nil
}

func (pi *TestItem) CreateDnXSX(param1 string) error {
	if pi.Caller == nil {
		return errors.New("caller not init")
	}
	var params []interface{}
	params = append(params, param1)
	_, err := pi.Caller("Test", "TestItem", "CreateDnXSX", params)
	if err != nil {
		if pi.TestItemLogger != nil {
			pi.TestItemLogger.Error("function response error", zap.Error(err), zap.String("function", "CreateDnXSX"), zap.String("interface", "TestItem"))
		}
		return err
	}
	return nil
}

func (pi *TestItem) CreateDnXSXX(param1 string) error {
	if pi.Caller == nil {
		return errors.New("caller not init")
	}
	var params []interface{}
	params = append(params, param1)
	_, err := pi.Caller("Test", "TestItem", "CreateDnXSXX", params)
	if err != nil {
		if pi.TestItemLogger != nil {
			pi.TestItemLogger.Error("function response error", zap.Error(err), zap.String("function", "CreateDnXSXX"), zap.String("interface", "TestItem"))
		}
		return err
	}
	return nil
}

func (pi *TestItem) GetSlice() ([]*DenemeItem, error) {
	if pi.Caller == nil {
		return nil, errors.New("caller not init")
	}
	var params []interface{}

	response, err := pi.Caller("Test", "TestItem", "GetSlice", params)
	if err != nil {
		if pi.TestItemLogger != nil {
			pi.TestItemLogger.Error("function response error", zap.Error(err), zap.String("function", "GetSlice"), zap.String("interface", "TestItem"))
		}
		return nil, err
	}
	fmt.Println("GetSlice", response)
	pReturn := response.([]interface{})
	var pReturn_1 []*DenemeItem
	if repeatItem, ok := pReturn[0].([]interface{}); ok {
		for _, val := range repeatItem {
			if mItem, ok := val.(map[string]interface{}); ok {
				nItem := NewDenemeItem()
				nItem.FromMap(mItem)
				pReturn_1 = append(pReturn_1, nItem)
			}
		}
	}
	return pReturn_1, nil
}

func (pi *TestItem) Init() error {
	if pi.Caller == nil {
		return errors.New("caller not init")
	}
	var params []interface{}

	_, err := pi.Caller("Test", "TestItem", "Init", params)
	if err != nil {
		if pi.TestItemLogger != nil {
			pi.TestItemLogger.Error("function response error", zap.Error(err), zap.String("function", "Init"), zap.String("interface", "TestItem"))
		}
		return err
	}
	return nil
}

func (pi *TestItem) Name() (string, error) {
	if pi.Caller == nil {
		return "", errors.New("caller not init")
	}
	var params []interface{}

	response, err := pi.Caller("Test", "TestItem", "Name", params)
	if err != nil {
		if pi.TestItemLogger != nil {
			pi.TestItemLogger.Error("function response error", zap.Error(err), zap.String("function", "Name"), zap.String("interface", "TestItem"))
		}
		return "", err
	}
	pReturn := response.([]interface{})
	var pReturn_1 string
	if mItem, ok := pReturn[0].(string); ok {
		pReturn_1 = mItem
	}
	return pReturn_1, nil
}

func (pi *TestItem) SetLogger(param1 *zap.Logger) error {
	if pi.Caller == nil {
		return errors.New("caller not init")
	}
	var params []interface{}
	params = append(params, param1)
	_, err := pi.Caller("Test", "TestItem", "SetLogger", params)
	if err != nil {
		if pi.TestItemLogger != nil {
			pi.TestItemLogger.Error("function response error", zap.Error(err), zap.String("function", "SetLogger"), zap.String("interface", "TestItem"))
		}
		return err
	}
	return nil
}

func (pi *TestItem) FromMap(data map[string]interface{}) {
	if dataItem, ok := data["DenemeItem"]; ok {
		_denemeitem := NewDenemeItem()
		if mapItem, ok := dataItem.(map[string]interface{}); ok {
			_denemeitem.FromMap(mapItem)
		}
		pi.DenemeItem = _denemeitem
	}

	if dataItem, ok := data["DenemeItemSlice"]; ok {
		var _denemeitemslice []DenemeItem
		repeatItem, isOk := dataItem.([]interface{})
		if isOk {
			for _, repItem := range repeatItem {
				if mapItem, ok := repItem.(map[string]interface{}); ok {
					iItem := NewDenemeItem()
					iItem.FromMap(mapItem)
					_denemeitemslice = append(_denemeitemslice, *iItem)
				}
			}
		}
		pi.DenemeItemSlice = _denemeitemslice
	}

	if dataItem, ok := data["Zaman"]; ok {
		if iData, ok := dataItem.(time.Time); ok {
			pi.Zaman = iData
		}
	}

}

func NewTestItem() *TestItem {
	mng, err := manager.GetManager()
	if err != nil {
		panic(err)
	}
	mng.RegisterPlugin(&manager.PluginConfig{
		Name:         "Test",
		Binary:       false,
		RemotePlugin: true,
		Addr:         "localhost",
		Port:         50051,
		NoValidate:   true,
	})
	plgItem := &TestItem{
		TestItemLogger: mng.GetPluginLogger("Test"),
		Caller:         mng.CallFunction,
	}

	return plgItem

}
