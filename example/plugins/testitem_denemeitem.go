
package plugins

import (
	"errors"
	"go.uber.org/zap"
	"github.com/antandros/go-yaps/manager"
	
)

type DenemeItem struct {
	Ben string
	Caller func(string, string,string,[]interface{}) (interface{}, error)
	DenemeItemLogger *zap.Logger
}

func (pi *DenemeItem) SetStr (param1 string) error {
	if pi.Caller  == nil {
		return  errors.New("caller not init")
	}
	var params []interface{}
	params = append(params,param1)
	_, err := pi.Caller("Test","DenemeItem","SetStr",params)
	if err != nil {
		if pi.DenemeItemLogger != nil {
			pi.DenemeItemLogger.Error("function response error", zap.Error(err), zap.String("function","SetStr"), zap.String("interface","DenemeItem"))
		}
		return  err
	}
	return nil
}

func (pi *DenemeItem) FromMap(data map[string]interface{}) {
	if dataItem, ok :=  data["Ben"]; ok {
		if iData, ok :=  dataItem.(string); ok {
			pi.Ben = iData
		}
	}
	
}

func NewDenemeItem() (*DenemeItem)  {
	
	mng ,err := manager.GetManager()
	if err != nil {
		panic(err)
	}
	
	plgItem := &DenemeItem{
		DenemeItemLogger:mng.GetPluginLogger("Test"),
		Caller:         mng.CallFunction,
	}
	
	return plgItem

}

