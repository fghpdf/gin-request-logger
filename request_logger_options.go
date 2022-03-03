package gin_request_logger

import "go.uber.org/zap"

type Options struct {
	LogResponse bool
	Logger      *zap.Logger
}
