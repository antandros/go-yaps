from plugin.models import BaseStruct , StructItem
import inspect
import json
class Parser:
    item : any
    StructData: BaseStruct
    
    def __init__(self, elem) -> None:
        self.item = elem
    def parseFunction(self,name , item):
        z =inspect.signature(item)
        
        print("\t", z, z.return_annotation.__name__)
        pass
    def parseItem(self, item, isbase=False):
        insp = inspect.getmembers(item)
        for key, item in insp:
            if not key.startswith("__"):
                if inspect.ismethod(item):
                    self.parseFunction(key, item)
            elif key == "__annotations__":
                print(item)
                for field, val in item.items(): 
                    print(field, type(field), val)
            else:
                print(key, item)
        
    def Parse(self):
        clsobj = self.item.__class__
        name = clsobj.__name__
        self.StructData = BaseStruct()
        self.StructData.Name = name
        itemStruct = StructItem()
        self.parseItem(self.item)
        
        
       
