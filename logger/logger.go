package logger

type Logger interface {
	Debug(msg string, v ...interface{})
	Info(msg string, v ...interface{})
	Warn(msg string, v ...interface{})
	Error(msg string, v ...interface{})
}
