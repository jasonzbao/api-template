package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Define a custom type for the context key to avoid collisions.
type loggerKey struct{}

// Singleton logger
var logger *zap.Logger

func Configure(env string) error {
	var zapConfig zap.Config
	if env == "local" {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
		// Overwrite sampling for now
		zapConfig.Sampling.Initial = 1
		zapConfig.Sampling.Thereafter = 1
	}

	var err error
	logger, err = zapConfig.Build()
	if err != nil {
		return err
	}

	// Set it as default global logger
	zap.ReplaceGlobals(logger)
	return nil
}

func Sync() error {
	return logger.Sync()
}

func G() *zap.Logger {
	if logger == nil {
		err := Configure("prod")
		if err != nil {
			return nil
		}
	}
	return logger
}

func NewContext(ctx context.Context, fields ...zapcore.Field) context.Context {
	//nolint:staticcheck
	return context.WithValue(ctx, loggerKey{}, *WithContext(ctx).With(fields...))
}

/*
*
/* WithContext returns a zap logger stored in context.
/* Usage:
/* logger.WithContext(c).Info("Hello world")
*/
func WithContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return G()
	}

	if ctxLogger, ok := ctx.Value(loggerKey{}).(zap.Logger); ok {
		return &ctxLogger
	} else {
		return G()
	}
}
