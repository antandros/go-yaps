from plugin.models import BaseStruct , StructItem , Method, Param, Field
import inspect
import json
class Parser:
    item : any
    StructData: BaseStruct
    
    def __init__(self, elem) -> None:
        self.item = elem
        
    def converMainTypes(self,name:str) -> str:
        if name == "<class 'str'>" or name == "str":
            return "string"
        elif name == "<class 'int'>" or name == "int":
            return "int64"
        elif name == "<class 'float'>" or name == "float":
            return "float64"
        elif name == "<class 'bool'>" or name == "bool":
            return "bool"
        else:
            return "interface{}"
        
    def convertTypeName(self,typeName:str) -> (str, bool, bool) :
        isPtr = False
        isSlice = False
        if typeName.startswith("*"):
            typeName = typeName[1:]
            isPtr = True
        if typeName.startswith("list["):
            isSlice = True
            typeName = typeName.split("[")[1].split("]")[0]
        
        if typeName.startswith("<class"):
            typeName = self.converMainTypes(typeName)
            
        elif typeName.find(".") > -1:
            typeName =  self.converMainTypes(typeName.split(".")[-1]) 
        
        return typeName, isPtr, isSlice
    
    def parseFunction(self, name:str , item:object, index:int) -> Method:
        z =inspect.signature(item)
        mt = Method()
        if type(z.return_annotation) == tuple:
            pIndex = 1
            for param in z.return_annotation:
                p = Param()
                typeName = str(param)
                typeName, p.IsPtr, p.IsSlice = self.convertTypeName(typeName)
                p.Type = typeName
                p.Index = pIndex
                p.TypeKind = typeName
                mt.appendOutParam(p)
                pIndex +=1
                
                
        else:
            p = Param()
            typeName = str(z.return_annotation)
            typeName, p.IsPtr, p.IsSlice = self.convertTypeName(typeName)
            p.Type = typeName
            p.Index = 0
            p.TypeKind = typeName
            mt.appendOutParam(p)
        
        mt.Name = name
        mt.Index = index
        pIndex=1
        
        for paramName, param in dict(z.parameters).items():
            p =Param()
            p.Index = pIndex
            typeName = str(param.annotation)
            typeName, p.IsPtr, p.IsSlice = self.convertTypeName(typeName)
            p.Type = typeName
            p.TypeKind = typeName
            mt.appendInParam(p)
            pIndex+= 1
        
        return mt
        
    def parseItem(self, item, name:str, isbase=False) -> StructItem:
        utem = StructItem()
        utem.Name = name
        utem.TypeName = "struct"
        insp = inspect.getmembers(item)
        funcIndex = 1
        for key, item in insp:
            if not key.startswith("__"):
                if inspect.ismethod(item):
                    gg = self.parseFunction(key, item, funcIndex)
                    utem.Methods.append(gg)
                    funcIndex+=1
            elif key == "__annotations__":
                fIndex = 0
                for field, val in item.items(): 
                    fi = Field()
                    typeName = str(val)
                    typeName, fi.IsPtr, fi.IsSlice = self.convertTypeName(typeName)
                    fi.Name = field
                    fi.Index = fIndex
                    fi.Type = typeName
                    fIndex+=1
                    utem.Fields.append(fi)
        return utem
                
        
        
    def Parse(self) -> BaseStruct:
        clsobj = self.item.__class__
        name = clsobj.__name__
        self.StructData = BaseStruct()
        self.StructData.Name = name
        self.StructData.TypeName = "struct"
        itemStruct = self.parseItem(self.item, name)
        self.StructData.Item = itemStruct
        return self.StructData
        
        
       
