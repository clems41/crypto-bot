package sqlTest

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
)

// deleteDatabase will delete database that is named dbName using database pointer
func deleteDatabase(DB *gorm.DB, dbName string, sqlDriver string) (err error) {
	err = execSqlCommand(DB, sqlDriver, disconnectUsersFromDatabase, dbName)
	if err != nil {
		return
	}
	err = execSqlCommand(DB, sqlDriver, dropDatabase, gorm.Expr(dbName))
	return errors.WithStack(err)
}

func execSqlCommand(DB *gorm.DB, sqlDriver string, command string, args ...interface{}) (err error) {
	// If command doesn't exist on map, return error
	if sqlCommandsByDriver[command] == nil {
		return fmt.Errorf("command %s is not defined in sqlCommandsByDriver", command)
	}
	// If command is not defined for one driver, skip it
	sqlCommand := sqlCommandsByDriver[command][sqlDriver]
	if sqlCommand == "" {
		return
	}
	return DB.
		Exec(sqlCommand, args...).
		Error
}

// getTemplateDatabaseName return name of template that will be used to test
func getTemplateDatabaseName(applicationName string) (name string) {
	name = fmt.Sprintf("%s_%s", templateDatabaseNamePrefix, applicationName)
	// We remove all - from database name, neither it will return errors
	name = strings.ReplaceAll(name, "-", "_")
	return
}

// getTemplateDatabaseName return name of new unique testing database that can be used to test
func getNewTestingDatabaseName(applicationName string) (name string) {
	// We create a new unique database tp avoid collision when tests will run in parallel
	uniqueID := uuid.New().ID()
	name = fmt.Sprintf("%s_%s_%d", testingDatabaseNamePrefix, applicationName, uniqueID)
	// We remove all - from database name, neither it will return errors
	name = strings.ReplaceAll(name, "-", "_")
	return
}