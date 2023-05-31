package logger

import (
	"go.uber.org/zap"
)

var cfg zap.Config

func LoggerInitializer(level string) (*zap.SugaredLogger, error) {
	cfg = zap.NewProductionConfig()
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	cfg.Level = lvl
	logger := zap.Must(cfg.Build())
	defer logger.Sync()
	return logger.Sugar(), nil
}
