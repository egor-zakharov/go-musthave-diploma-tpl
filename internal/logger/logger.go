package logger

import "go.uber.org/zap"

func Log() *zap.Logger {
	return zap.L()
}

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(zl)
	return nil
}
