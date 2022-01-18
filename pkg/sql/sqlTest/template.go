package sqlTest

import (
	"cine-circle-api/pkg/sql/sqlConnection"
	"cine-circle-api/pkg/utils/sliceUtils"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// CreateTestTemplateDatabase will create template database and migrate all repositories for creating database schemas, tables, columns, indexes, etc...
// This database template is used for testing. Create new database and migrate all repositories at the beginning of each test will take too much time.
// So, we are creating this template before running all testcases, and for each test, new testing database will be created from this template (without doing any migrations).
// To work, this function needs :
//  - repositoriesMigration function that define how repositories should be migrated
//  - logger used to log some info and errors from defer func
func CreateTestTemplateDatabase(repositoriesMigration func(DB *gorm.DB) (err error), logger Logger) (err error) {
	// Try to connect to repositorySQL storage using actual database
	defaultConfig, err := sqlConnection.GetConfigFromEnv()
	if err != nil {
		return
	}
	sqlDriver := defaultConfig.SqlDriver
	if !sliceUtils.SliceContainsStr(allowedDriversWithCreateTestTemplateDatabase, sqlDriver) {
		return fmt.Errorf("SqlDriver %s is not supported, only %v are supported", sqlDriver, allowedDriversWithCreateTestTemplateDatabase)
	}
	currentDB, err := sqlConnection.Open(&defaultConfig)
	if err != nil {
		return
	}
	// Disconnect at the end
	defer func() {
		deferErr := sqlConnection.Close(currentDB)
		if deferErr != nil {
			logger.Errorf(err.Error())
		}
	}()

	// Delete all old testing database (cleanup)
	err = cleanupTestingDatabases(currentDB, sqlDriver)
	if err != nil {
		return
	}
	logger.Infof("Old testing and template databases successfully dropped")

	// Create new template database
	templateDatabaseName := getTemplateDatabaseName(defaultConfig.ApplicationName)
	err = execSqlCommand(currentDB, sqlDriver, createDatabase, gorm.Expr(templateDatabaseName))
	if err != nil {
		return errors.WithStack(err)
	}

	// Open connection on freshly created template database
	templateConfig := defaultConfig
	templateConfig.DbName = templateDatabaseName
	templateDB, err := sqlConnection.Open(&templateConfig)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		deferErr := sqlConnection.Close(templateDB)
		if deferErr != nil {
			logger.Errorf(err.Error())
		}
	}()
	logger.Infof("New template database successfully created : %s", templateDatabaseName)

	// Migrate all tables into template database
	err = repositoriesMigration(templateDB)
	if err != nil {
		return errors.WithStack(err)
	}
	logger.Infof("All repositories have been successfully migrated on new template database %s", templateDatabaseName)

	// Make sure the template is not modified
	err = execSqlCommand(currentDB, sqlDriver, limitConnectionToTemplateDatabase, gorm.Expr(templateDatabaseName))
	if err != nil {
		return errors.WithStack(err)
	}
	logger.Infof("Template database %s successfully locked", templateDatabaseName)
	return
}

// cleanupTestingDatabases will delete all previous testing databases that have not been already deleted.
// It will also delete previous template database before creating new one.
func cleanupTestingDatabases(DB *gorm.DB, sqlDriver string) (err error) {
	// Delete all old testing database (cleanup)
	list := make([]struct {
		DatabaseName string
	}, 0)
	err = DB.
		Raw(sqlCommandsByDriver[listDatabasesLike][sqlDriver], testingDatabaseNamePrefix+"%").
		Scan(&list).
		Error
	if err != nil {
		return errors.WithStack(err)
	}
	for _, elem := range list {
		if elem.DatabaseName != "" {
			err = deleteDatabase(DB, elem.DatabaseName, sqlDriver)
			if err != nil {
				return err
			}
		}
	}
	return
}
