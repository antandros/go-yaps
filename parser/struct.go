package parser

import (
	"fmt"
	"reflect"
	"strings"
)

type Param struct {
	Index       int
	PackagePath string
	Type        string
	TypeKind    uint
	IsStruct    bool
	IsPtr       bool
	IsSlice     bool
	Obj         reflect.Type `json:"-"`
}

type Method struct {
	Name        string
	PackagePath string
	InParams    []*Param
	OutParams   []*Param
	NamedString string
	Index       int
}

func (mt *Method) GenerateEmptyReturn() string {
	var rets []string
	for _, fld := range mt.OutParams {
		switch strings.ToLower(fld.Type) {
		case "string":
			if fld.IsSlice {
				rets = append(rets, `nil`)
			} else {

				rets = append(rets, `""`)
			}
		case "error":
			continue
		case "int":
		case "uint":
		case "uint8":
		case "int8":
		case "uint16":
		case "int16":
		case "uint32":
		case "int32":
		case "uint64":
		case "int64":
			if fld.IsSlice {
				rets = append(rets, "nil")
			} else {

				rets = append(rets, "0")
			}
		default:
			rets = append(rets, "nil")
		}
	}

	return strings.Join(rets, ",")
}
func (mt *Method) HasReturnError() bool {
	for _, fld := range mt.OutParams {
		if strings.EqualFold(fld.Type, "error") {
			return true
		}
	}
	return false
}
func (mt *Method) LenOut() int {
	return len(mt.OutParams)
}
func (mt *Method) LenOutIfEq(ln int) bool {
	return len(mt.OutParams) == ln
}
func (mt *Method) LenOutIfGt(ln int) bool {
	return len(mt.OutParams) > ln
}
func (mt *Method) ParseParams(method reflect.Method) {
	for i := 1; i < method.Type.NumIn(); i++ {

		inParam := method.Type.In(i)
		inParamOrg := method.Type.In(i)
		var isPtr bool
		var isSlice bool
		if inParam.Kind() == reflect.Slice {
			isSlice = true
		}
		if inParam.Kind() == reflect.Ptr {
			inParam = inParam.Elem()
			isPtr = true
		}

		pr := &Param{
			Obj:         inParam,
			Type:        inParamOrg.String(),
			TypeKind:    uint(inParam.Kind()),
			PackagePath: inParamOrg.PkgPath(),
			IsStruct:    inParam.Kind() == reflect.Struct,
			IsPtr:       isPtr,
			IsSlice:     isSlice,
			Index:       i,
		}
		mt.InParams = append(mt.InParams, pr)

	}
	for i := 0; i < method.Type.NumOut(); i++ {
		inParam := method.Type.Out(i)
		inParamOrg := method.Type.Out(i)
		var isPtr bool
		var isSlice bool
		if inParam.Kind() == reflect.Slice {
			isSlice = true
			inParam = inParam.Elem()
		}
		if inParam.Kind() == reflect.Ptr {
			inParam = inParam.Elem()
			isPtr = true
		}

		if inParamOrg.Kind() == reflect.Slice {
			inParamOrg = inParamOrg.Elem()
		}
		if inParamOrg.Kind() == reflect.Ptr {
			inParamOrg = inParamOrg.Elem()
		}
		fmt.Println()
		pr := &Param{
			Obj:         inParam,
			Type:        inParamOrg.Name(),
			PackagePath: inParamOrg.PkgPath(),
			IsStruct:    inParam.Kind() == reflect.Struct,
			IsPtr:       isPtr,
			IsSlice:     isSlice,
			Index:       i,
		}
		fmt.Println(pr)

		mt.OutParams = append(mt.OutParams, pr)

	}
}
func (mt *Method) Parse(method reflect.Method) {
	mt.Name = method.Name
	mt.PackagePath = method.PkgPath
	mt.Index = method.Index
	mt.NamedString = method.Type.String()
	mt.ParseParams(method)
}

type Field struct {
	Name        string
	Type        string
	Index       int
	TypeKind    string
	IsStruct    bool
	IsPtr       bool
	IsSlice     bool
	IsSame      bool
	PackagePath string
	Obj         reflect.Type `json:"-"`
}

func (fi *Field) Parse(field reflect.StructField, pkg string) bool {
	var isPtr bool
	var isStruct bool
	var slice bool
	var IsSame bool
	if field.Type.Kind() == reflect.Ptr {
		isPtr = true
		noPtr := field.Type.Elem()
		if noPtr.Kind() == reflect.Struct {
			isStruct = true
		}
		if noPtr.Kind() == reflect.Slice {
			slice = true
		}

	} else {
		if field.Type.Kind() == reflect.Struct {
			isStruct = true
		}
		if field.Type.Kind() == reflect.Slice {
			slice = true
		}
	}
	if slice {
		sliceElem := field.Type.Elem()
		if sliceElem.Kind() == reflect.Ptr {
			sliceElem = sliceElem.Elem()
		}
		if sliceElem.Kind() == reflect.Struct {
			isStruct = true
		}
	}
	if strings.EqualFold(field.Type.PkgPath(), pkg) || field.Type.PkgPath() == "" {
		IsSame = true
	}
	fi.Name = field.Name
	fi.PackagePath = field.Type.PkgPath()
	fi.IsPtr = isPtr
	fi.IsSlice = slice
	fi.IsSame = IsSame
	fi.IsStruct = isStruct
	typeFn := field.Type
	if typeFn.Kind() == reflect.Ptr {
		typeFn = typeFn.Elem()
	}
	if typeFn.Kind() == reflect.Slice {
		typeFn = typeFn.Elem()
		if typeFn.Kind() == reflect.Ptr {
			typeFn = typeFn.Elem()
		}
	}

	fi.Type = typeFn.Name()

	if !IsSame {
		fi.Type = fmt.Sprint(field.Type)
	}
	fi.TypeKind = field.Type.Kind().String()
	return isStruct
}

type StructItem struct {
	Name        string
	Fields      []Field
	Methods     []*Method
	PackagePath string
	Imports     []string
	TypeName    string
	BasePacket  string
	Object      reflect.Type `json:"-"`
}

func (si *StructItem) GetBaseItem(item reflect.Type) reflect.Type {
	itemType := item
	if itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}
	if itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}
	return itemType
}
func (si *StructItem) AddImport(imp string) {
	have := false
	for _, name := range si.Imports {
		if strings.EqualFold(name, imp) {
			have = true
		}
	}
	if !have {
		si.Imports = append(si.Imports, imp)
	}
}
func (si *StructItem) Parse(parser func(reflect.Type)) {
	outItem := si.GetBaseItem(si.Object)
	si.Name = outItem.Name()

	si.PackagePath = outItem.PkgPath()
	if si.PackagePath != si.BasePacket {
		si.AddImport(si.PackagePath)
	}
	si.TypeName = outItem.Kind().String()
	if si.TypeName == "" {
		fmt.Println(outItem.String())
	}
	itemNew := reflect.New(outItem)
	for m := 0; m < itemNew.NumMethod(); m++ {

		mth := itemNew.Type().Method(m)
		meth := &Method{}
		meth.Parse(mth)
		si.Methods = append(si.Methods, meth)

	}
	for i := 0; i < outItem.NumField(); i++ {
		field := outItem.Field(i)
		fieldObj := Field{
			Index: i,
			Obj:   field.Type,
		}
		fieldObj.Parse(field, si.BasePacket)
		si.Fields = append(si.Fields, fieldObj)

		if field.Type.PkgPath() != si.BasePacket {
			if field.Type.PkgPath() != "" {

				si.AddImport(field.Type.PkgPath())
			}
		}
	}
}

type Struct struct {
	Name         string        `json:"name"`
	BasePacket   string        `json:"basepacket"`
	Imports      []string      `json:"imports"`
	Item         *StructItem   `json:"item"`
	RelatedItems []*StructItem `json:"related"`
	Object       any           `json:"-"`
}

func (s *Struct) ParseMethods(relatedStrck *StructItem) {
	for _, mth := range relatedStrck.Methods {
		for _, fld := range mth.InParams {
			if fld.IsStruct {
				obj := s.ConvertItem(fld.Obj)
				s.ParseStructItem(obj)
			}
		}
		for _, fld := range mth.OutParams {
			if fld.IsStruct {
				obj := s.ConvertItem(fld.Obj)
				s.ParseStructItem(obj)
			}
		}

	}
}
func (s *Struct) ConvertItem(obj reflect.Type) reflect.Type {
	if obj.Kind() == reflect.Slice {
		obj = obj.Elem()
	}
	if obj.Kind() == reflect.Ptr {
		obj = obj.Elem()
	}
	return obj
}
func (s *Struct) ParseFields(relatedStrck *StructItem) {
	for _, fld := range relatedStrck.Fields {
		if fld.IsStruct {
			obj := s.ConvertItem(fld.Obj)
			s.ParseStructItem(obj)
		}
	}
}
func (s *Struct) ParseStructItem(obj reflect.Type) {
	refObject := obj
	excludeImports := []string{"go.uber.org/zap"}
	itemPackage := s.Item.GetBaseItem(obj).PkgPath()
	if !strings.EqualFold(itemPackage, s.BasePacket) {

		have := false
		for _, n := range s.Imports {
			if strings.EqualFold(n, itemPackage) {
				have = true
			}
		}
		for _, n := range excludeImports {
			if strings.EqualFold(n, itemPackage) {
				have = true
			}
		}
		if !have {
			s.Imports = append(s.Imports, itemPackage)
		}
		return
	}
	if obj.Kind() == reflect.Ptr {
		refObject = obj.Elem()
	}
	for _, alreadyRelated := range s.RelatedItems {

		if strings.EqualFold(alreadyRelated.Name, refObject.Name()) {

			return
		}
	}
	relatedStrck := &StructItem{
		Object:     obj,
		BasePacket: s.BasePacket,
	}
	relatedStrck.Parse(s.ParseStructItem)
	sIndex := len(s.RelatedItems)
	s.RelatedItems = append(s.RelatedItems, relatedStrck)
	s.ParseFields(s.RelatedItems[sIndex])
	s.ParseMethods(s.RelatedItems[sIndex])

}

func (s *Struct) Parse() {
	s.Item = &StructItem{
		Object: reflect.TypeOf(s.Object),
	}
	s.Item.BasePacket = s.Item.GetBaseItem(s.Item.Object).PkgPath()
	s.BasePacket = s.Item.GetBaseItem(s.Item.Object).PkgPath()
	s.Item.Parse(s.ParseStructItem)
	s.ParseFields(s.Item)
	s.ParseMethods(s.Item)
}

func NewStructParser(item any) Struct {
	stck := Struct{Object: item}
	stck.Parse()
	return stck
}
