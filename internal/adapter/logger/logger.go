package logger

import "log"

type Logger interface {
	Println(v ...interface{})
}

type LogAdapter struct {
	logger *log.Logger
}

func New(logger *log.Logger) *LogAdapter {
	return &LogAdapter{logger: logger}
}

func (l *LogAdapter) Println(v ...interface{}) {
	l.logger.Println(v...)
}
