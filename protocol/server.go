package protocol

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"time"

	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	PluginProtocolServer
	socket   string
	server   *grpc.Server
	listener net.Listener
	logger   *zap.Logger
	UnimplementedPluginProtocolServer
	ConfigFunction func() map[string]interface{}
	Caller         func(string, string, []*InTypes) any
}

func (s Server) Run() error {
	s.ConfigFunction()
	if err := s.server.Serve(s.listener); err != nil {
		return err
	}
	return nil
}
func (s Server) CallFunction(ctx context.Context, req *FunctionRequest) (*FunctionResponse, error) {
	s.logger.Info("Call function on server", zap.Any("req", req), zap.Any("params", req.In))
	resp := s.Caller(req.GetFunction(), req.GetStruct(), req.GetIn())
	data, err := json.Marshal(resp)
	fncResp := &FunctionResponse{
		Data: data,
	}
	return fncResp, err
}
func (s Server) Stat(ctx context.Context, req *Empty) (*StatResponse, error) {
	return nil, errors.New("not imp yet")
}
func (s Server) mustEmbedUnimplementedPluginProtocolServer() {}
func (s Server) HeartBeat(ctx context.Context, req *Empty) (*Empty, error) {
	return nil, nil
}
func (s Server) RequestConfig(ctx context.Context, req *Empty) (*ConfigResponse, error) {
	var resp *ConfigResponse
	if s.ConfigFunction != nil {
		respFnc := s.ConfigFunction()

		respData, err := json.Marshal(respFnc)
		if err != nil {
			return nil, err
		}
		resp = &ConfigResponse{
			Data:    respData,
			Success: true,
		}

	} else {
		resp = &ConfigResponse{
			Data:    []byte(""),
			Success: false,
		}
	}
	return resp, nil
}

type PluginManagerInterface interface {
	GenConfig() map[string]interface{}
	GetLogger() *zap.Logger
	PluginCaller(string, string, []*InTypes) any
}

func NewServer(ltype string, socket string, p PluginManagerInterface) (Server, error) {
	lis, err := net.Listen(ltype, socket)
	if err != nil {
		return Server{}, err
	}
	handler := &Handler{
		logger: p.GetLogger(),
	}

	var kaep = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}

	var kasp = keepalive.ServerParameters{
		MaxConnectionIdle:     60 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
		MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
		MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  1 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout:               15 * time.Second, // Wait 1 second for the ping ack before assuming the connection is dead
	}

	gServer := grpc.NewServer(grpc.StatsHandler(handler), grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))
	srv := Server{
		socket:         socket,
		server:         gServer,
		listener:       lis,
		logger:         p.GetLogger(),
		Caller:         p.PluginCaller,
		ConfigFunction: p.GenConfig,
	}
	RegisterPluginProtocolServer(srv.server, srv)
	return srv, nil

}
