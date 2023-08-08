package logger

import "log"

type ILogger interface {
	Println(v ...interface{})
}

type logAdapter struct {
	logger *log.Logger
}

func New(logger *log.Logger) *logAdapter {
	return &logAdapter{logger: logger}
}

func (l *logAdapter) Println(v ...interface{}) {
	l.logger.Println(v...)
}
