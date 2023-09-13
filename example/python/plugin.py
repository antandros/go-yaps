import protocol_pb2
import protocol_pb2_grpc
import grpc
from concurrent import futures


class PluginManager(protocol_pb2_grpc.PluginProtocolServicer):
    def RequestConfig(self, request, context):
        return super().RequestConfig(request, context)
    def CallFunction(self, request, context):
        return super().CallFunction(request, context)

    def __init__(self) -> None:
        pass
    

class Plugin:
    def __init__(self) -> None:
        pass
    
    def Start(self):
        port = "50051"
        server = grpc.server(
            thread_pool=futures.ThreadPoolExecutor(max_workers=10),
        )
        protocol_pb2_grpc.add_GreeterServicer_to_server(PluginManager(), server)
        server.add_insecure_port("[::]:" + port)
        server.start()
        print("Server started, listening on " + port)
        server.wait_for_termination()