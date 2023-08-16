package logger

type logger interface {
	Println(v ...interface{})
}

type LogAdapter struct {
	log logger
}

func New(l logger) *LogAdapter {
	return &LogAdapter{log: l}
}

func (l *LogAdapter) Println(v ...interface{}) {
	l.log.Println(v...)
}
