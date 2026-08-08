package logger

func (l *Logger) LogRequest(msg string, status int, attrs ...any) {
	switch {
	case status >= 500:
		l.Error(msg, attrs...)
	case status >= 400:
		l.Warn(msg, attrs...)
	default:
		l.Info(msg, attrs...)
	}
}
