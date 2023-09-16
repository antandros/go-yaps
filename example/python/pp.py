

from typing import List
from plugin.plugin import Plugin
from dd.den import DenemeItemX
from datetime import datetime
from bson import objectid
class DenemeItem:
    Ben:str
    
    def SetStr(self,elem):
        self.Ben = elem

class TestItem:
    DenemeItemx : list[DenemeItem]
    DenemeItemx2 : list[DenemeItem]
    DenemeItemx3 : list[DenemeItem]
    X:str
    def __init__(self) -> None:
        self.DenemeItemx = []
        
    def Init2(self) -> str:
        return ""
    def Init(self)  -> str:
        return "x"
    def AddSlic2exxx(self, elem:list[DenemeItemX], elem1:DenemeItemX, elem2:bool, elem3:float, elem4:int, elem5:datetime, elem6:objectid.ObjectId)  -> list[DenemeItemX]:
        return self.DenemeItemx
    
    def AddSlic2e(self, elem:str)  -> list[DenemeItem]:
        return self.DenemeItemx
    def AddSlic2e(self, elem:str)  -> (DenemeItem, DenemeItemX):
        return self.DenemeItemx
    
    def AddSlic3e(self, elem1:list[DenemeItem], elem2:list[DenemeItem])  -> list[DenemeItem]:
        return self.DenemeItemx
    
    def AddSlic3e4(self, elem1:list[DenemeItem], elem3:str)  -> list[DenemeItem]:
        return self.DenemeItemx
    
    def SetLogger(self, elem:any)  -> None:
        pass
    def Name(self)  -> str:
        return "pytest"
    
    def GetSlice(self)  -> list[DenemeItem]:
        return self.DenemeItemx
    
    def AddSlice(self, elem:str)  -> None:
        d = DenemeItem()
        d.SetStr(elem + "--- added")
        self.DenemeItemx.append(d)

if __name__ == "__main__":
    p = TestItem()
    pp = Plugin(p)
    pp.Start()