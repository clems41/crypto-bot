package sqlTest

import "cine-circle-api/pkg/sql/sqlConnection"

const (
	testingDatabaseNamePrefix  = "testing"
	templateDatabaseNamePrefix = "testing_template"
)

// Define which sqlDriver can be used with OpenCleanDatabaseFromTemplate method
var (
	allowedDriversWithOpenCleanDatabaseFromTemplate = []string{
		sqlConnection.PostgresqlDriver,
	}
)

// Define which sqlDriver can be used with CreateTestTemplateDatabase method
var (
	allowedDriversWithCreateTestTemplateDatabase = []string{
		sqlConnection.PostgresqlDriver,
	}
)

// Define which sqlDriver can be used with OpenCleanDatabase method
var (
	allowedDriversWithOpenCleanDatabase = []string{
		sqlConnection.PostgresqlDriver,
		sqlConnection.MysqlDriver,
	}
)

// Define here all SQL commands that can be used test template
const (
	createDatabase                    = "CREATE_DATABASE"
	limitConnectionToTemplateDatabase = "LIMIT_CONNECTION_TO_TEMPLATE_DATABASE"
	listDatabasesLike                 = "LIST_DATABASES_LIKE"
	disconnectUsersFromDatabase       = "DISCONNECT_USERS_FROM_DATABASE"
	dropDatabase                      = "DROP_DATABASE"
	createDatabaseFromTemplate        = "CREATE_DATABASE_FROM_TEMPLATE"
)

// sqlCommandsByDriver map all SQL commands to each SQL driver
var sqlCommandsByDriver = map[string]map[string]string{
	createDatabase: {
		sqlConnection.PostgresqlDriver: "CREATE DATABASE ?",
		sqlConnection.MysqlDriver:      "CREATE DATABASE ?",
	},
	limitConnectionToTemplateDatabase: {
		sqlConnection.PostgresqlDriver: "ALTER DATABASE ? WITH ALLOW_CONNECTIONS 0",
		sqlConnection.MysqlDriver:      "",
	},
	listDatabasesLike: {
		sqlConnection.PostgresqlDriver: "SELECT datname as database_name FROM pg_database WHERE datname LIKE ?",
		sqlConnection.MysqlDriver:      "SELECT schema_name as database_name FROM information_schema.schemata WHERE schema_name LIKE ?",
	},
	disconnectUsersFromDatabase: {
		sqlConnection.PostgresqlDriver: "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = ?",
		sqlConnection.MysqlDriver:      "",
	},
	dropDatabase: {
		sqlConnection.PostgresqlDriver: "DROP DATABASE IF EXISTS ?",
		sqlConnection.MysqlDriver:      "DROP DATABASE IF EXISTS ?",
	},
	createDatabaseFromTemplate: {
		sqlConnection.PostgresqlDriver: "CREATE DATABASE ? WITH TEMPLATE ?",
		sqlConnection.MysqlDriver:      "DOESNT EXIST ON MYSQL",
	},
}
