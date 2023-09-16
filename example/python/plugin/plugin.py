from plugin.parser import Parser
from plugin.plugin_manager import PluginManager
import plugin.protocol_pb2
from plugin import protocol_pb2_grpc
import grpc
import inspect
from concurrent import futures


class Plugin:
    pitem = None
    parsed = None
    def __init__(self, pitem) -> None:
        self.pitem = pitem
        pn = Parser(pitem)
        result = pn.Parse()
        self.parsed = result.to_dict()
        
    def Start(self):
        server_options = [
            ("grpc.keepalive_time_ms", 20000),
            ("grpc.keepalive_timeout_ms", 10000),
            ("grpc.http2.min_ping_interval_without_data_ms", 5000),
            ("grpc.max_connection_idle_ms", 10000),
            ("grpc.max_connection_age_ms", 30000),
            ("grpc.max_connection_age_grace_ms", 5000),
            ("grpc.http2.max_pings_without_data", 5),
            ("grpc.keepalive_permit_without_calls", 1),
        ]
        port = "50051"
        server = grpc.server(
            thread_pool=futures.ThreadPoolExecutor(max_workers=10),
             options=server_options,

        )
        protocol_pb2_grpc.add_PluginProtocolServicer_to_server(PluginManager(self.parsed,self.pitem), server)
        server.add_insecure_port("[::]:" + port)
        server.start()
        print("Server started, listening on " + port)
        server.wait_for_termination()