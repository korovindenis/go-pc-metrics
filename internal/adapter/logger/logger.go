package logger

type logger interface {
	Println(v ...interface{})
}

type LogAdapter struct {
	logger logger
}

func New(l logger) *LogAdapter {
	return &LogAdapter{logger: l}
}

func (l *LogAdapter) Println(v ...interface{}) {
	l.logger.Println(v...)
}
