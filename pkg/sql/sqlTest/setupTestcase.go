package sqlTest

import (
	"cine-circle-api/pkg/logger"
	"cine-circle-api/pkg/sql/sqlConnection"
	"cine-circle-api/pkg/utils/sliceUtils"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"testing"
)

// OpenCleanDatabaseFromTemplate create new empty database based on template.
// All queries will be executed on this new fresh database.
// Only postgreSQL is supported as SQL Driver for this method.
//  !!! You should have created template database before calling this method (cf. CreateTestTemplateDatabase) !!!
func OpenCleanDatabaseFromTemplate(t *testing.T) (DB *gorm.DB, closeFunction func()) {
	// Open connection with real database for creating testing database
	realDBConfig, err := sqlConnection.GetConfigFromEnv()
	if err != nil {
		t.Fatalf(err.Error())
	}
	sqlDriver := realDBConfig.SqlDriver
	if !sliceUtils.SliceContainsStr(allowedDriversWithOpenCleanDatabaseFromTemplate, sqlDriver) {
		t.Fatalf("SqlDriver %s is not supported, only %v are supported", sqlDriver, allowedDriversWithOpenCleanDatabaseFromTemplate)
	}
	realDB, err := sqlConnection.Open(&realDBConfig)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Deleting testing database if already exists before creating new one from template
	testingDatabaseName := getNewTestingDatabaseName(realDBConfig.ApplicationName)
	err = deleteDatabase(realDB, testingDatabaseName, sqlDriver)
	require.NoError(t, err)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Create new testing database using template
	templateDatabaseName := getTemplateDatabaseName(realDBConfig.ApplicationName)
	err = execSqlCommand(realDB, sqlDriver, createDatabaseFromTemplate, gorm.Expr(testingDatabaseName), gorm.Expr(templateDatabaseName))
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Using new clean database for testing
	testingConfig := realDBConfig
	testingConfig.DbName = testingDatabaseName
	testingDB, err := sqlConnection.Open(&testingConfig)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// create closeFunction : to be used in testcases
	closeFunction = func() {
		// closing testing database before deleting it from real database
		err = sqlConnection.Close(testingDB)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// delete testing database
		err = deleteDatabase(realDB, testingDatabaseName, sqlDriver)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// closing connection from real database
		err = sqlConnection.Close(realDB)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	return testingDB, closeFunction
}

// OpenCleanDatabase create new empty database based on migration method.
// All queries will be executed on this new fresh database.
// Only postgreSQL adn MySQL is supported as SQL Driver for this method.
func OpenCleanDatabase(t *testing.T, repositoriesMigration func(DB *gorm.DB) (err error)) (DB *gorm.DB, closeFunction func()) {
	// Open connection with real database for creating testing database
	realDBConfig, err := sqlConnection.GetConfigFromEnv()
	if err != nil {
		t.Fatalf(err.Error())
	}
	sqlDriver := realDBConfig.SqlDriver
	if !sliceUtils.SliceContainsStr(allowedDriversWithOpenCleanDatabase, sqlDriver) {
		t.Fatalf("SqlDriver %s is not supported, only %v are", sqlDriver, allowedDriversWithOpenCleanDatabase)
	}
	realDB, err := sqlConnection.Open(&realDBConfig)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Deleting testing database if already exists before creating new empty one
	testingDatabaseName := getNewTestingDatabaseName(realDBConfig.ApplicationName)
	err = deleteDatabase(realDB, testingDatabaseName, sqlDriver)
	require.NoError(t, err)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = execSqlCommand(realDB, sqlDriver, createDatabase, gorm.Expr(testingDatabaseName))
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Open connection on freshly created testing database
	testingConfig := realDBConfig
	testingConfig.DbName = testingDatabaseName
	testingDB, err := sqlConnection.Open(&testingConfig)
	if err != nil {
		t.Fatalf(err.Error())
	}
	logger.Infof("New testing database successfully created : %s", testingDatabaseName)

	// Migrate all tables into testing database
	err = repositoriesMigration(testingDB)
	if err != nil {
		t.Fatalf(err.Error())
	}
	logger.Infof("All repositories have been successfully migrated on new testing database %s", testingDatabaseName)

	// create closeFunction : to be used in testcases
	closeFunction = func() {
		// closing testing database before deleting it from real database
		err = sqlConnection.Close(testingDB)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// delete testing database
		err = deleteDatabase(realDB, testingDatabaseName, sqlDriver)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// closing connection from real database
		err = sqlConnection.Close(realDB)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	return testingDB, closeFunction
}
