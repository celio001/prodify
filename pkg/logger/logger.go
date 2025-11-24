package logger

import (
	"strings"

	"go.uber.org/zap"
)

var Log *zap.Logger

func Init(env string) {

	configLogger := configLogger(env)

	logger, err := configLogger.Build()
	if err != nil {
		panic(err)
	}

	Log = logger

	Log = Log.With(
		zap.String("env", strings.ToLower(env)),
	)
}
