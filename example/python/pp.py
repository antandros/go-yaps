

from typing import List
from plugin.plugin import Plugin

class DenemeItem:
    Ben:str
    
    def SetStr(self,elem):
        self.Ben = elem

class Pluginitem:
    DenemeItemx : List[DenemeItem]
    DenemeItemx2 : List[DenemeItem]
    DenemeItemx3 : List[DenemeItem]
    def Init2(self) -> str:
        return ""
    def Init(self)  -> str:
        return "x"
    def AddSlice(self, elem:str)  -> str:
        d = DenemeItem()
        d.SetStr(elem)
        return ""
if __name__ == "__main__":
    p = Pluginitem()
    pp = Plugin(p)
    pp.Start()