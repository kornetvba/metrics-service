package logger

import "go.uber.org/zap"

var Logger = zap.NewNop()

func Initialization() error {
	level, err := zap.ParseAtomicLevel("info")
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = level

	lg, err := cfg.Build()
	if err != nil {
		return err
	}
	Logger = lg
	return nil

}
