from plugin.parser import Parser
import plugin.protocol_pb2
import grpc
import inspect


class Plugin:
    pitem = None
    parsed = None
    def __init__(self, pitem) -> None:
        self.pitem = pitem
        pn = Parser(pitem)
        pn.Parse()
        
    def Start(self):
        return
        port = "50051"
        server = grpc.server(
            thread_pool=futures.ThreadPoolExecutor(max_workers=10),
        )
        protocol_pb2_grpc.add_PluginProtocolServicer_to_server(PluginManager(), server)
        server.add_insecure_port("[::]:" + port)
        server.start()
        print("Server started, listening on " + port)
        server.wait_for_termination()