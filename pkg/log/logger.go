package log

import (
	"go.uber.org/zap"
)

var l *zap.Logger

// Log Mode
const (
	DEVELOPMENT = "develop"
	PRODUCE     = "produce"
)

func InitLog(mode string) error {
	var err error
	switch mode {
	case DEVELOPMENT:
		l, err = zap.NewDevelopment()
	case PRODUCE:
		l, err = zap.NewProduction()
	default:
		l = zap.NewExample()
	}
	return err
}

func Logger() *zap.Logger {
	return l
}
