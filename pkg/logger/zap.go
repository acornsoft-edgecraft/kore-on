// Package logger - Logger for Zap
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// zapLogger - SugaredLogger
type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

// ===== [ Implementations ] =====

// Debug - Implements of Logger's Debug
func (zl *zapLogger) Debug(args ...interface{}) {
	zl.sugaredLogger.Debug(args...)
}

// Debugf - Implements of Logger's Debugf
func (zl *zapLogger) Debugf(format string, args ...interface{}) {
	zl.sugaredLogger.Debugf(format, args...)
}

// Info - Implements of Logger's Info
func (zl *zapLogger) Info(args ...interface{}) {
	zl.sugaredLogger.Info(args...)
}

// Infof - Implements of Logger's Infof
func (zl *zapLogger) Infof(format string, args ...interface{}) {
	zl.sugaredLogger.Infof(format, args...)
}

// Warn - Implements of Logger's Warn
func (zl *zapLogger) Warn(args ...interface{}) {
	zl.sugaredLogger.Warn(args...)
}

// Warnf - Implements of Logger's Warnf
func (zl *zapLogger) Warnf(format string, args ...interface{}) {
	zl.sugaredLogger.Warnf(format, args...)
}

// Error - Implements of Logger's Error
func (zl *zapLogger) Error(args ...interface{}) {
	zl.sugaredLogger.Error(args...)
}

// Errorf - Implements of Logger's Errorf
func (zl *zapLogger) Errorf(format string, args ...interface{}) {
	zl.sugaredLogger.Errorf(format, args...)
}

// Fatal - Implements of Logger's Fatal
func (zl *zapLogger) Fatal(args ...interface{}) {
	zl.sugaredLogger.Fatal(args...)
}

// Fatalf - Implements of Logger's Fatalf
func (zl *zapLogger) Fatalf(format string, args ...interface{}) {
	zl.sugaredLogger.Fatalf(format, args...)
}

// Panic - Implements of Logger's Panic
func (zl *zapLogger) Panic(args ...interface{}) {
	zl.sugaredLogger.Panic(args...)
}

// Panicf - Implements of Logger's Panicf
func (zl *zapLogger) Panicf(format string, args ...interface{}) {
	zl.sugaredLogger.Fatalf(format, args...)
}

// WithFields - Implements of Logger's WithFields
func (zl *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	newLogger := zl.sugaredLogger.With(f...)
	return &zapLogger{newLogger}
}

// WithError - Implements of Logger's WithError
func (zl *zapLogger) WithError(err error) Logger {
	zl.sugaredLogger.Error(err)
	return &zapLogger{zl.sugaredLogger}
}

// ===== [ Private Functions ] =====

// getEncoder - Returns an encoder for Zap
func getEncoder(isJSON bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getZapLevel - Returns an log level for Zap
func getZapLevel(level string) zapcore.Level {
	switch level {
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// newZapLogger - returns an instance of Logger or Error for Zap
func newZapLogger(conf Config) (Logger, error) {
	cores := []zapcore.Core{}
	if conf.EnableConsole {
		level := getZapLevel(conf.ConsoleLevel)
		writer := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(getEncoder(conf.ConsoleJSONFormat), writer, level)
		cores = append(cores, core)
	}

	if conf.EnableFile {
		level := getZapLevel(conf.FileLevel)
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename: conf.FileLocation,
			MaxSize:  100,
			Compress: true,
			MaxAge:   28,
		})
		core := zapcore.NewCore(getEncoder(conf.FileJSONFormat), writer, level)
		cores = append(cores, core)
	}

	combinedCore := zapcore.NewTee(cores...)

	// AddCallerSkip skip 2 number of callers, this is important else the file that gets
	// logged will always be the wrapped file. In our case zap.go
	logger := zap.New(combinedCore, zap.AddCallerSkip(3), zap.AddCaller()).Sugar()

	return &zapLogger{
		sugaredLogger: logger,
	}, nil
}

// ===== [ Public Functions ] =====
