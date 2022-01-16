package logger

import (
	"go.uber.org/zap"
)

// CustomLogger allow to log in your project without caring about which logger have been chosen (Zap actually https://github.com/uber-go/zap).
// Custom methods are already implemented, using different level of logging (error, info, debug, etc.).
// Other like Print and Printf are defined in order to be able to use this logger with httpServer package define in this project.
type CustomLogger struct{
	// https://github.com/uber-go/zap
	logger *zap.SugaredLogger
}

var (
	customLogger *CustomLogger
)

func (customLogger *CustomLogger) Debugf(template string, args ...interface{}) {
	customLogger.logger.Debugf(template, args...)
}

func (customLogger * CustomLogger) Infof(template string, args ...interface{}) {
	customLogger.logger.Infof(template, args...)
}

func (customLogger * CustomLogger) Warnf(template string, args ...interface{}) {
	customLogger.logger.Warnf(template, args...)
}

func (customLogger * CustomLogger) Errorf(template string, args ...interface{}) {
	customLogger.logger.Errorf(template, args...)
}

func (customLogger * CustomLogger) Fatalf(template string, args ...interface{}) {
	customLogger.logger.Fatalf(template, args...)
}

func (customLogger * CustomLogger) Debug(args ...interface{}) {
	customLogger.logger.Debug(args...)
}

func (customLogger * CustomLogger) Info(args ...interface{}) {
	customLogger.logger.Info(args...)
}

func (customLogger * CustomLogger) Warn(args ...interface{}) {
	customLogger.logger.Warn(args...)
}

func (customLogger * CustomLogger) Error(args ...interface{}) {
	customLogger.logger.Error(args...)
}

func (customLogger * CustomLogger) Fatal(args ...interface{}) {
	customLogger.logger.Fatal(args...)
}

func (customLogger * CustomLogger) Print(args ...interface{}) {
	customLogger.logger.Info(args...)
}

func (customLogger * CustomLogger) Printf(template string, args ...interface{}) {
	customLogger.logger.Infof(template, args...)
}
