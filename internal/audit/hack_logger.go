package audit

import (
	"fmt"
	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func (zl ZapLogger) Log(keyvals ...interface{}) error {
	if len(keyvals)%2 != 0 {
		return fmt.Errorf("odd number of arguments")
	}
	var loggerFunc func(args ...interface{})
	logger := zl.logger.WithOptions(zap.AddCallerSkip(2))
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			return fmt.Errorf("non-string key: %v", keyvals[i])
		}
		if key == "level" {
			switch keyvals[i+1] {
			case "debug":
				loggerFunc = logger.Debug
			case "warn":
				loggerFunc = logger.Warn
			case "error":
				loggerFunc = logger.Error
			case "info":
				loggerFunc = logger.Info
			default:
				loggerFunc = logger.Info
			}
		}
		if key == "msg" {
			loggerFunc(keyvals[i+1])
			return nil
		}
	}
	return nil
}
