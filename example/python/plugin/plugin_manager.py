from plugin.protocol_pb2_grpc import PluginProtocolServicer
from plugin.protocol_pb2 import ConfigResponse , FunctionResponse
from concurrent import futures
import json
class PluginManager(PluginProtocolServicer):
    parsed = None
    pluginItem = None
    def __init__(self, pars, pluginItem) -> None:
        self.parsed = pars
        self.pluginItem = pluginItem
        super().__init__()
    
    def RequestConfig(self, request, context):
        data = json.dumps(self.parsed).encode('utf-8')
        return ConfigResponse(success=True,data=data)
    
    def CallFunction(self, request, context):
        attrs = None
        baseAttr = []
        
        if hasattr(request, "in"):
            attrs = getattr(request, "in")
        if attrs:
            for attr in attrs:
                if hasattr(attr, "in"):
                    inData = getattr(attr,"in")
                    converted = json.loads(inData)
                    baseAttr.append(converted)
        print(request)
        if hasattr(self.pluginItem,request.function):
            fnc = getattr(self.pluginItem,request.function)
            resp = fnc(*baseAttr)
            respData = b""
            if type(resp) is list:
                respItem = []
                for it in resp:
                    respItem.append(it.__dict__)
                respData = json.dumps(respItem).encode('utf-8')
            else:
                if resp is not None:
                    respData = json.dumps(resp.__dict__).encode('utf-8')
            
            return FunctionResponse(data=respData,success=True)
        return FunctionResponse(data=None,success=False)
            
                