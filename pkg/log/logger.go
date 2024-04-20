package log

import (
	"go.uber.org/zap"
)

var sl *zap.SugaredLogger

// Log Mode
const (
	DEVELOPMENT = "develop"
	PRODUCE     = "produce"
)

func InitLog(mode string) error {
	var err error
	var logger *zap.Logger
	switch mode {
	case DEVELOPMENT:
		logger, err = zap.NewDevelopment()
	case PRODUCE:
		logger, err = zap.NewProduction()
	default:
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return err
	}
	sl = logger.Sugar()
	return nil
}

func Logger() *zap.SugaredLogger {
	return sl
}
