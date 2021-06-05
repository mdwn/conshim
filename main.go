package main

import (
	"os"

	"github.com/meowfaceman/conshim/cmd"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	setupLogger()

	cmd.Execute()
}

func setupLogger() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	logger := zap.New(zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.InfoLevel))

	zap.ReplaceGlobals(logger)
}
