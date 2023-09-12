package yaperror

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type YapsError int

const (
	UNKNOW YapsError = iota
	FILE_NOT_FOUND
	FILE_CREATE_INTERFACE
	CANT_PARSE_PLUGIN
	PARSE_MESSAGE_BODY
	NOT_CONNECTED
	KILL_ERROR
	DISCONNECT_CLIENT
	PARSE_MESSAGE_NAME
	PARSE_MESSAGE_GENERAL
	PARSE_INVALID_ITEM_TYPE
	PARSE_KEY_NOT_FOUND
	CANT_INIT
	CANT_CALL
	ZERO_REGISTERED_PLUGIN
	REGISTER_MANAGER
	INTERFACE_CREATE
	PROTOCOL_SERVER
	CONFIG_GET_ERROR
	CLIENT_RESPONSE_ERROR
	VALIDATE_ITEM
	RUN_BINARY
	MANAGER_NOT_FOUND
	PLUGIN_BINARY_NOT_FOUND
)

func (y YapsError) String() string {
	switch y {
	case UNKNOW:
		return "unknow error"
	case RUN_BINARY:
		return "can not run client binary"
	case DISCONNECT_CLIENT:
		return "cant disconnect client"
	case VALIDATE_ITEM:
		return "cant validate function call"
	case NOT_CONNECTED:
		return "plugin not connected"
	case KILL_ERROR:
		return "cant kill process"
	case CONFIG_GET_ERROR:
		return "plugin config response unsuccessfull"
	case CANT_CALL:
		return "cant call function"
	case FILE_NOT_FOUND:
		return "file not found"
	case PROTOCOL_SERVER:
		return "cant create a protocol server"
	case CANT_INIT:
		return "cant init a manager"
	case INTERFACE_CREATE:
		return "cant create a interface"
	case ZERO_REGISTERED_PLUGIN:
		return "cant find any registered plugin"
	case REGISTER_MANAGER:
		return "please register manager before call"
	case MANAGER_NOT_FOUND:
		return "manager not found in plugin client"
	case PARSE_INVALID_ITEM_TYPE:
		return "cant determinate a item type"
	case CANT_PARSE_PLUGIN:
		return "cant parse plugin"
	case PLUGIN_BINARY_NOT_FOUND:
		return "cant found plugin binary in binary folder"
	}
	return ""
}

type YAPSError struct {
	Code    YapsError
	Message string
	Err     error
	Extra   map[string]interface{}
}

func (e *YAPSError) ZapError(logger *zap.Logger) {
	var fields []zapcore.Field
	for key, val := range e.Extra {
		fields = append(fields, zap.Any(key, val))
	}
	fields = append(fields, zap.Error(e.Err))
	fields = append(fields, zap.Int("code", int(e.Code)))
	logger.Error(e.Code.String(), fields...)
}
func (e *YAPSError) ErrorCode() string {
	return e.Code.String()
}
func (e *YAPSError) Error() string {
	var extras []string
	for key, val := range e.Extra {
		extras = append(extras, fmt.Sprintf("%s: %v", key, val))
	}
	extraText := strings.Join(extras, ", ")
	var errText string
	if e.Err != nil {
		errText = fmt.Sprintf("Wrapped Error: %v", e.Err)
	}
	return fmt.Sprintf("ErroCode : %d, Message: %s, %s %s", e.Code, e.Code.String(), errText, extraText)
}

type ErrorOptions struct {
	Message string
	Extra   map[string]interface{}
}

func WithExra(extra map[string]interface{}) Options {
	return func(c *ErrorOptions) {
		c.Extra = extra
	}
}
func WithMessage(msg string) Options {
	return func(c *ErrorOptions) {
		c.Message = msg
	}
}

type Options func(*ErrorOptions)

func Error(errType YapsError, err error, options ...Options) *YAPSError {
	config := &ErrorOptions{}
	for _, opt := range options {
		opt(config)
	}
	yaperror := &YAPSError{
		Code:    errType,
		Message: config.Message,
		Extra:   config.Extra,
		Err:     err,
	}

	return yaperror
}
