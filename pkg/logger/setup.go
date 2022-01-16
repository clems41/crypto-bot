package logger

import (
	"crypto-bot/pkg/utils/envUtils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// init will create new zap logger (https://github.com/uber-go/zap).
// Please set DEBUG_LOGs env variable to false if you don't want to use debug log level.
func init() {
	// Define if debug should be logged
	debugMode := envUtils.GetFromEnvOrDefault(envDebug, defaultDebug) == "true"
	var cfg zap.Config
	if debugMode {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	// Logger configuration
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	cfg.Encoding = "console"
	log, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	// Instantiate logger
	zapLogger := log.Sugar()
	customLogger = &CustomLogger{
		logger: zapLogger,
	}
	zapLogger.Infof("Logger has been successfully instantiated")
}
