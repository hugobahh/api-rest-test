package logger

type Logger interface {
	Debug(moduleName, functionName, msg string)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warning(moduleName, functionName, msg string)
	Warningf(moduleName, functionName, format string, args ...any)
	Error(moduleName, functionName string, err error)
	Errorf(moduleName, functionName, format string, args ...any)
	FatalIfError(moduleName, functionName string, errs ...error)
	Fatal(moduleName, functionName string, err error)
	Fatalf(moduleName, functionName, format string, args ...any)
}
