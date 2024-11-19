package mysql

type Logger interface {
	Error(moduleName, functionName string, err error)
	Fatal(moduleName, functionName string, err error)
	Infof(format string, args ...any)
}
