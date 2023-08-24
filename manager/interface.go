package manager

import "go.uber.org/zap"

type PluginInterface interface {
	Name() string
	SetLogger(*zap.Logger)
	Init() error
}
