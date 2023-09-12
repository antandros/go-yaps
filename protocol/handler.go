package protocol

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/stats"
)

type Handler struct {
	logger *zap.Logger
}

func (h *Handler) TagRPC(c context.Context, t *stats.RPCTagInfo) context.Context {
	//h.logger.Info("TagRPC", zap.Any("t", t))
	return c
}

// HandleRPC processes the RPC stats.
func (h *Handler) HandleRPC(c context.Context, hg stats.RPCStats) {
	//h.logger.Info("HandleRPC", zap.Any("hand", hg))
}

func (h *Handler) TagConn(c context.Context, s *stats.ConnTagInfo) context.Context {

	//h.logger.Info("Tag Conn", zap.Any("s", s))
	return c
}

// HandleConn processes the Conn stats.
func (h *Handler) HandleConn(c context.Context, s stats.ConnStats) {
	h.logger.Info("CONN  con", zap.Any("stat", s))
	switch s.(type) {
	case *stats.ConnEnd:

		h.logger.Info("client  disconnected", zap.Any("stat", s))

	}
}

// init
