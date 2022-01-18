package sqlConnection

const (
	envApplicationName = "APPLICATION_NAME"
	envSqlDriver       = "DB_SQL_DRIVER"
	envHost            = "DB_HOST"
	envPort            = "DB_PORT"
	envUser            = "DB_USER"
	envDbName          = "DB_NAME"
	envPassword        = "DB_PASSWORD"
	envExtraConfigs    = "DB_EXTRA_CONFIGS"
	envDebug           = "DB_DEBUG"
	envDetailedLogs    = "DB_LOG"
)

const (
	defaultSqlDriver    = "postgresql"
	defaultDebug        = "true"
	defaultDetailedLogs = "true"
)

// Default values with PostgreSQL as SqlDriver
const (
	defaultPostgresqlHost         = "localhost"
	defaultPostgresqlPort         = "5432"
	defaultPostgresqlUser         = "postgres"
	defaultPostgresqlPassword     = "postgres"
	defaultPostgresqlDbName       = "postgres"
	defaultPostgresqlExtraConfigs = "sslmode=disable TimeZone=Pacific/Noumea"
)

// Default values with MySQL as SqlDriver
const (
	defaultMysqlHost         = "localhost"
	defaultMysqlPort         = "3307"
	defaultMysqlUser         = "root"
	defaultMysqlPassword     = "mysql"
	defaultMysqlDbName       = "test"
	defaultMysqlExtraConfigs = "charset=utf8mb4&parseTime=True&loc=Local"
)

// SQL drivers list
const (
	PostgresqlDriver = "postgresql"
	MysqlDriver      = "mysql"
)

var (
	allowedSqlDrivers = []string{
		PostgresqlDriver,
		MysqlDriver,
	}
)
