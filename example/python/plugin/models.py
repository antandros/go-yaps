
from typing import Any


class Param:
    PackagePath :str
    Type        :str
    TypeKind    :int
    IsStruct    :bool
    IsPtr       :bool
    IsSlice     :bool
    Index       :int
    def __init__(self) -> None:
        self.PackagePath = ""
        self.Type        =  ""
        self.TypeKind    = 0
        self.IsStruct    = False
        self.IsPtr       = False
        self.IsSlice     = False
        self.Index       =  0

class Method:
    Name:str
    NamedString: str
    InParams: list[Param]
    OutParams: list[Param]
    Index: int
    def __init__(self) -> None:
        self.InParams = []
        self.OutParams = []
    def appendInParam(self,item:Param) -> None :
        self.InParams.append(item)
        
    def appendOutParam(self,item:Param) -> None :
        self.OutParams.append(item)
   
    def to_dict(self) -> dict:
        initems = []
        outitems = []
        for i in self.InParams:
            initems.append(i.__dict__)
        for i in self.OutParams:
            outitems.append(i.__dict__)
        return {
            "Name":self.Name,
            "Index":self.Index,
            "InParams":initems,
            "OutParams":outitems,
        }
    
class Field:
    Name        :str
    Type        :str
    Index       :int
    TypeKind    :str
    IsStruct    :bool
    IsPtr       :bool
    IsSlice     :bool
    IsSame      :bool
    PackagePath :str
    
class StructItem:
    Name:str
    Fields: list[Field]
    Methods: list[Method]
    PackagePath: str
    Imports: list[str]
    TypeName: str
    BasePacket: str
    def __init__(self) -> None:
        self.Fields = []
        self.Methods = []
        self.Imports = []
        self.PackagePath = ""
        self.TypeName = ""
        self.BasePacket = ""
        self.Name = ""
    def to_dict(self) -> dict:
        fls = []
        mtd = []
        for i in self.Fields:
            fls.append(i.__dict__)
        for i in self.Methods:
            mtd.append(i.to_dict())
        return {
            "Name":self.Name,
            "Fields":fls,
            "Methods":mtd,
            "PackagePath":self.PackagePath,
            "Imports":[],
            "TypeName":self.TypeName,
            "BasePacket":self.BasePacket,
        }
        
        
class BaseStruct:
    Imports: list[str]
    TypeName: str
    Name:str
    Item: StructItem
    RelatedItems: list[StructItem]
    def to_dict(self):
        return {
            "Imports":[],
            "TypeName":self.TypeName,
            "Name":self.Name,
            "Item":self.Item.to_dict(),
            "RelatedItems":[],
        }
