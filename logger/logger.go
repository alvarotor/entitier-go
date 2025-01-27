package logger

type Logger interface {
	Info(id string, message string)
	Error(id string, message string)
	Debug(id string, message string)
}
