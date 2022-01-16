package logger

// Logger return pointer on custom logger.
// In this way you can use this logger in other librairies that allow to set custom logger.
func Logger() *CustomLogger {
	return customLogger
}

// Implementation of methods that can be called directly without using custom logger pointer.

func Debugf(template string, args ...interface{}) {
	customLogger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	customLogger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	customLogger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	customLogger.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	customLogger.Fatalf(template, args...)
}

func Debug(args ...interface{}) {
	customLogger.Debug(args...)
}

func Info(args ...interface{}) {
	customLogger.Info(args...)
}

func Warn(args ...interface{}) {
	customLogger.Warn(args...)
}

func Error(args ...interface{}) {
	customLogger.Error(args...)
}

func Fatal(args ...interface{}) {
	customLogger.Fatal(args...)
}
