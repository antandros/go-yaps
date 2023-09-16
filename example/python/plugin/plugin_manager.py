from plugin.protocol_pb2_grpc import PluginProtocolServicer
from concurrent import futures

class PluginManager(PluginProtocolServicer):
    
    def RequestConfig(self, request, context):
        return super().RequestConfig(request, context)
    def CallFunction(self, request, context):
        return super().CallFunction(request, context)

                