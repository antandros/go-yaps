package protocol

import (
	"context"
	"encoding/json"
	"errors"
	"net"

	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
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
	gServer := grpc.NewServer(grpc.StatsHandler(handler))
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
