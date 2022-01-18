package sqlConnection

import "fmt"

func errSqlDriverNotSupported(sqlDriver string, allowedSqlDrivers []string) error {
	return fmt.Errorf("SQL Driver %s is not supported. Please use one of %v", sqlDriver, allowedSqlDrivers)
}
