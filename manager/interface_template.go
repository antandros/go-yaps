package manager

func ManagerTemplate() string {
	template := `
	`
	return template
}
func GetTemplate() string {
	template := `
package {{ package_name }}

import (
	"errors"
	"go.uber.org/zap"
	"github.com/antandros/go-yaps/manager"
	{% for import in item.Imports %}"{{import}}"{% endfor %}
)

type {{ item.Name }} struct {
{% for field in item.Fields  %}	{{ field.Name }} {% if field.IsSlice %}[]{% endif %}{% if field.IsPtr %}*{% endif %}{{field.Type}}
{% endfor %}	Caller func(string, string,string,[]interface{}) (interface{}, error)
	{{item.Name}}Logger *zap.Logger
}
{% for fnc in item.Methods %}
func (pi *{{ item.Name }}) {{fnc.Name}} ({%-for ins in fnc.InParams %}param{{ins.Index}} {{ins.Type}} 
{%- endfor %}) {% if fnc.LenOutIfGt(0) and fnc.HasReturnError() == False %}(
	{%-for ins in fnc.OutParams %}{% if ins.IsSlice %}[]{% endif %}{% if ins.IsPtr %}*{% endif %}{{ins.Type}} {%- if not loop.last %},{% endif %}
{%- endfor %}{% if not fnc.HasReturnError() %},error{% endif %}){% else %}error{% endif %} {
	if pi.Caller  == nil {
		return {% if fnc.LenOut() > 0  and fnc.HasReturnError() == False %}{{fnc.GenerateEmptyReturn()}},{% endif %} errors.New("caller not init")
	}
	var params []interface{}
	{%for ins in fnc.InParams %}params = append(params,param{{ins.Index}}){% endfor %}
	{%- if fnc.LenOut() >  0  and not fnc.HasReturnError()  %}
	response, err := pi.Caller("{{plugName}}","{{item.Name}}","{{fnc.Name}}",params)
	{%- else %}
	_, err := pi.Caller("{{plugName}}","{{item.Name}}","{{fnc.Name}}",params)
	{%- endif %}
	if err != nil {
		if pi.{{item.Name}}Logger != nil {
			pi.{{item.Name}}Logger.Error("function response error", zap.Error(err), zap.String("function","{{fnc.Name}}"), zap.String("interface","{{item.Name}}"))
		}
		return {% if fnc.LenOut() > 0 and not fnc.HasReturnError() %}{{fnc.GenerateEmptyReturn()}},{% endif %} err
	}
	{%- if fnc.LenOut() > 0 and not fnc.HasReturnError() %}
	pReturn := response.([]interface{})
	{%- for oparam in fnc.OutParams%}
	{%- if oparam.IsSlice %}
	var pReturn_{{loop.index}} []{%- if oparam.IsPtr%}*{%- endif %}{{oparam.Type}}
	if repeatItem, ok := pReturn[{{loop.index-1}}].([]interface{}); ok {
		for _ , val := range repeatItem {
			{%- if oparam.IsStruct   %}
			if mItem, ok := val.(map[string]interface{}); ok {
				nItem := New{{oparam.Type}}()
				nItem.FromMap(mItem)
				pReturn_{{loop.index}} = append(pReturn_{{loop.index}}, {%- if not oparam.IsPtr%}*{%- endif %}nItem)
			}
			{%- endif %}
		}
	}
	{%- else %}
	{%- if oparam.IsStruct and oparam.IsSame  %}
	var pReturn_{{loop.index}} {%- if oparam.IsPtr%}*{%- endif %} {{oparam.Type}}
	if mItem, ok := pReturn[{{ loop.index - 1 }}].(map[string]interface{}); ok {
		nItem := New{{oparam.Type}}()
		nItem.FromMap(mItem)
		pReturn_{{loop.index}} = {%- if oparam.IsPtr%}*{%- endif %}nItem
	}
	{%- else %}
	var pReturn_{{loop.index}} {%- if oparam.IsPtr%}*{%- endif %} {{oparam.Type}}
	if mItem, ok := pReturn[{{ loop.index - 1 }}].({%- if oparam.IsPtr%}*{%- endif %}{{oparam.Type}}); ok {
		pReturn_{{loop.index}} = mItem
	}
	{%- endif %}
	{%- endif %}
	{%- endfor %}
	return {%for oparam in fnc.OutParams%}pReturn_{{loop.index}}{% if not loop.last %},{% endif %}{% endfor %} {% if not fnc.HasReturnError() %}, nil{% endif %}

	{%- else %}
	return nil
	{%- endif %}
}
{% endfor %}
func (pi *{{ item.Name }}) FromMap(data map[string]interface{}) {
	
	{%- for field in item.Fields  %}
	if dataItem, ok :=  data["{{field.Name}}"]; ok {
	{%- if field.IsStruct and field.IsSame %}
	{%- if field.IsSlice %}
		var  _{{field.Name|lower}} []{% if field.IsPtr %}*{% endif %}{{field.Type}}
		repeatItem , isOk :=dataItem.([]interface{})
		if isOk {
			for _ , repItem := range repeatItem {
				if mapItem, ok :=  repItem.(map[string]interface{}); ok {
					iItem := New{{field.Type}}()
					iItem.FromMap(mapItem)
					_{{field.Name|lower}} = append(_{{field.Name|lower}}, {% if not field.IsPtr %}*{% endif %}iItem)
				}
			}
		}
	{%- else %}
		_{{field.Name|lower}} := New{{field.Type}}()
		if mapItem, ok :=  dataItem.(map[string]interface{}); ok {
			_{{field.Name|lower}}.FromMap(mapItem)
		}
	{%- endif %}
		pi.{{ field.Name }} = _{{field.Name|lower}} 
	{%- else %}
	{%- if field.IsSlice %}
		var  _{{field.Name|lower}} []{% if field.IsPtr %}*{% endif %}{{field.Type}}
		repeatItem , isOk :=dataItem.([]interface{})
		if isOk {
			for _, kItem := range repeatItem {
				if iData, ok :=  kItem.({% if field.IsPtr %}*{% endif %}{{field.Type}}); ok {
					_{{field.Name|lower}} = append(_{{field.Name|lower}}, {% if not field.IsPtr %}*{% endif %}iData)
				}
			}
		}
	{%- else %}
		if iData, ok :=  dataItem.({% if field.IsPtr %}*{% endif %}{{field.Type}}); ok {
			pi.{{ field.Name }} = iData
		}
	{%- endif %}
	{%- endif %}
	}
	{% endfor %}
}
{% if inititem%}

func New{{ item.Name }}() (*{{ item.Name }})  {
	mng ,err := manager.GetManager()
	if err != nil {
		panic(err)
	}
	
	plgItem := &{{ item.Name }}{
		{{ item.Name }}Logger:mng.GetPluginLogger("{{plugName}}"),
		Caller:         mng.CallFunction,
	}
	mng.RegisterPlugin(&manager.PluginConfig{
		Name:   "{{ plugName }}",
		Binary: true,
	})
	return plgItem

}
{% else %}
func New{{ item.Name }}() (*{{ item.Name }})  {
	
	mng ,err := manager.GetManager()
	if err != nil {
		panic(err)
	}
	
	plgItem := &{{ item.Name }}{
		{{ item.Name }}Logger:mng.GetPluginLogger("{{plugName}}"),
		Caller:         mng.CallFunction,
	}
	
	return plgItem

}
{% endif %}

`

	return template
}
