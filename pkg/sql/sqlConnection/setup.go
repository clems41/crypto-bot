package sqlConnection

import (
	"cine-circle-api/pkg/logger"
	"cine-circle-api/pkg/utils/envUtils"
	"cine-circle-api/pkg/utils/sliceUtils"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Database struct {
	db *gorm.DB
}

type Config struct {
	SqlDriver       string
	Host            string
	User            string
	Password        string
	Port            string
	DbName          string
	ExtraConfigs    string
	Debug           bool
	DetailedLogs    bool
	ApplicationName string
}

// DataSourceName return connection string based on Config and SQL driver
func (cfg Config) DataSourceName() (source string) {

	switch cfg.SqlDriver {
	case MysqlDriver:
		source = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DbName,
			cfg.ExtraConfigs)
	default:
		source = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s application_name='%s' %s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DbName,
			cfg.ApplicationName,
			cfg.ExtraConfigs)
	}

	return
}

// Open create connection to database using GORM from Config
// If customConfig is nil, we will use default config using env variables with GetConfigFromEnv
func Open(customConfig *Config) (db *gorm.DB, err error) {

	var dbConfig Config
	if customConfig == nil {
		dbConfig, err = GetConfigFromEnv()
		if err != nil {
			return
		}
	} else {
		dbConfig = *customConfig
	}

	logLevel := gormLogger.Silent

	if dbConfig.Debug {
		logLevel = gormLogger.Info
	}

	// Use gorm logger if logs are requested
	gormCfg := gorm.Config{}
	if dbConfig.DetailedLogs {
		newLogger := gormLogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			gormLogger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logLevel,    // Log level
				Colorful:      true,        // Disable color
			},
		)
		gormCfg.Logger = newLogger
	}

	// Open connection with database depending on SQL driver specified (default: postgresql)
	switch dbConfig.SqlDriver {
	case MysqlDriver:
		db, err = gorm.Open(mysql.Open(dbConfig.DataSourceName()), &gormCfg)
	default:
		db, err = gorm.Open(postgres.Open(dbConfig.DataSourceName()), &gormCfg)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if db == nil {
		return nil, fmt.Errorf("database pointer is nil")
	}

	logger.Infof("Database connection has been successfully established to host %s with user %s",
		dbConfig.Host, dbConfig.User)

	return
}

// GetConfigFromEnv return Config using environment variables or default values based on SqlDriver
func GetConfigFromEnv() (config Config, err error) {
	sqlDriver := envUtils.GetFromEnvOrDefault(envSqlDriver, defaultSqlDriver)
	if !sliceUtils.SliceContainsStr(allowedSqlDrivers, sqlDriver) {
		return config, errSqlDriverNotSupported(sqlDriver, allowedSqlDrivers)
	}
	// Default config is based on environment variables or if not defined default values (depending on sqlDriver)
	switch sqlDriver {
	case PostgresqlDriver:
		config = Config{
			Host:         envUtils.GetFromEnvOrDefault(envHost, defaultPostgresqlHost),
			User:         envUtils.GetFromEnvOrDefault(envUser, defaultPostgresqlUser),
			Password:     envUtils.GetFromEnvOrDefault(envPassword, defaultPostgresqlPassword),
			Port:         envUtils.GetFromEnvOrDefault(envPort, defaultPostgresqlPort),
			DbName:       envUtils.GetFromEnvOrDefault(envDbName, defaultPostgresqlDbName),
			ExtraConfigs: envUtils.GetFromEnvOrDefault(envExtraConfigs, defaultPostgresqlExtraConfigs),
		}
	case MysqlDriver:
		config = Config{
			Host:         envUtils.GetFromEnvOrDefault(envHost, defaultMysqlHost),
			User:         envUtils.GetFromEnvOrDefault(envUser, defaultMysqlUser),
			Password:     envUtils.GetFromEnvOrDefault(envPassword, defaultMysqlPassword),
			Port:         envUtils.GetFromEnvOrDefault(envPort, defaultMysqlPort),
			DbName:       envUtils.GetFromEnvOrDefault(envDbName, defaultMysqlDbName),
			ExtraConfigs: envUtils.GetFromEnvOrDefault(envExtraConfigs, defaultMysqlExtraConfigs),
		}
	default:
		return config, errSqlDriverNotSupported(sqlDriver, allowedSqlDrivers)
	}
	config.SqlDriver = sqlDriver
	config.Debug = envUtils.GetFromEnvOrDefault(envDebug, defaultDebug) == "true"
	config.DetailedLogs = envUtils.GetFromEnvOrDefault(envDetailedLogs, defaultDetailedLogs) == "true"
	config.ApplicationName, err = envUtils.GetFromEnvOrError(envApplicationName)
	return
}

func Close(DB *gorm.DB) (err error) {
	sqlDB, err := DB.DB()
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(sqlDB.Close())
}
