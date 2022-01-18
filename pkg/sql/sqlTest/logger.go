package sqlTest

// Logger interface specify which functions that are needed to use a logger with this sqlTest package
type Logger interface {
	Errorf(template string, args ...interface{})
	Infof(template string, args ...interface{})
}
