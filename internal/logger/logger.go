package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var cfg zap.Config

func LoggerInitializer(level string) (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("during logger initialize error ocurred - %v", err)
	}
	defer logger.Sync()

	cfg = zap.NewProductionConfig()
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	cfg.Level = lvl
	return logger.Sugar(), nil
}
