
class Param:
	Index       :int
	PackagePath :str
	Type        :str
	TypeKind    :int
	IsStruct    :bool
	IsPtr       :bool
	IsSlice     :bool

class Method:
    Name:str
    NamedString: str
    InParams: list[Param]
    OutParams: list[Param]
    Index: int
    
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

class BaseStruct:
    Imports: list[str]
    TypeName: str
    Name:str
    Item: StructItem
    RelatedItems: list[StructItem]
