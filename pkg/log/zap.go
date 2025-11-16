package log

import (
	"os"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger zap.Logger

func Init(cfg *koanf.Koanf) {
	lvlErr := false
	fmtErr := false

	cfgLvl, err := zapcore.ParseLevel(cfg.String("level"))
	if err != nil {
		cfgLvl = zapcore.DebugLevel
		lvlErr = true
	}

	lvlEnablerFn := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= cfgLvl
	})

	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if cfg.Bool("color") {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	var encoder zapcore.Encoder

	switch cfg.String("format") {
	case "text":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	case "json", "":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
		fmtErr = true
	}

	output := zapcore.Lock(os.Stdout)
	lg := zap.New(zapcore.NewCore(encoder, output, lvlEnablerFn)).Named(cfg.String("title"))

	defaultLogger = lg
	_ = zap.ReplaceGlobals(lg)

	if lvlErr {
		defaultLogger.Error("log level specified in config is invalid, error is assumed")
	}

	if fmtErr {
		defaultLogger.Error("formatter specified in config is invalid, json is assumed")
	}
}
