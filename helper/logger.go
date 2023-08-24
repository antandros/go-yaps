package helper

import (
	"os"
	"path"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func createLogger(fileName string) *zap.Logger {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    4, // megabytes
		MaxBackups: 3,
		MaxAge:     2, // days
	})

	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	return zap.New(core)
}
func ReplaceGoPath(environ []string, newGoPath string) []string {
	newEnviron := []string{}
	for _, v := range environ {
		if !strings.HasPrefix(v, "GOPATH=") {
			newEnviron = append(newEnviron, v)
		}
	}
	return append(newEnviron, "GOPATH="+newGoPath)
}
func LoggerNamed(fname string) *zap.Logger {
	basePath := path.Dir(fname)
	_, err := os.Stat(basePath)
	if err != nil {
		os.MkdirAll(basePath, 0700)
	}
	return createLogger(fname)
}
func Logger() *zap.Logger {
	filePath := os.Getenv("LOG_PATH")
	if filePath == "" {
		filePath = "log/go-yaps.log"
	}
	return createLogger(filePath)
}
