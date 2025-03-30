package kafka

type (
	logger interface {
		Infof(msg string, args ...interface{})
		Errorf(msg string, args ...interface{})
		Warnf(msg string, args ...interface{})
	}
)
