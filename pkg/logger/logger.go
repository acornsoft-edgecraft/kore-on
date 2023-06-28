// Package logger - Common Logging (inspired https://www.mountedthoughts.com/golang-logger-interface/)
package logger

import (
	"errors"
	"os"
	"path/filepath"
)

// ===== [ Constants and Variables ] =====

const (
	// LevelDebug has verbose message
	LevelDebug = "debug"
	// LevelInfo is default log level
	LevelInfo = "info"
	// LevelWarn is for logging messages about possible issues
	LevelWarn = "warn"
	// LevelError is for logging errors
	LevelError = "error"
	// LevelFatal is for logging fatal messages. The system shutdown after logging the message
	LevelFatal = "fatal"

	// InstanceZapLogger indicates that using the instance of Zap logger
	InstanceZapLogger int = iota
	// InstanceLogrusLogger indicates that using the instance of Logrus logger
	InstanceLogrusLogger
)

var (
	// A global variable so that log functions can be directly acccessed
	log                      Logger
	errInvalidLoggerInstance = errors.New("Invalid logger instance")
)

// ===== [ Types ] =====

// Fields - Type to pass when we want to call WithFields for structrued logging
type Fields map[string]interface{}

// Logger - Interface for logging
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	WithFields(keyValues Fields) Logger
	WithError(err error) Logger
}

// Config - Stores the config for the logging. For some loggers there can only be one level across writers,
// for such the level of Console is picked by default
type Config struct {
	EnableConsole     bool
	ConsoleJSONFormat bool
	ConsoleLevel      string
	EnableFile        bool
	FileJSONFormat    bool
	FileLevel         string
	FileLocation      string
}

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// New - Create default logger
func New() error {
	executablePath, err := os.Executable()
	if err != nil {
		return err
	}
	// 실행 파일이 있는 디렉토리 경로 추출
	executableDir := filepath.Dir(executablePath)

	// 실행 파일의 상위 경로 추출
	parentDir := filepath.Dir(executableDir)

	// Sets the Common logging
	logConf := Config{
		EnableConsole: true,
		ConsoleLevel:  LevelDebug,
		//ConsoleJSONFormat: true,
		EnableFile:     true,
		FileLevel:      LevelInfo,
		FileJSONFormat: true,
		FileLocation:   parentDir + "/logs/koreonctl.log",
	}

	return NewLogger(logConf, InstanceZapLogger)
}

// NewLogger - Returns an instance of logger
func NewLogger(conf Config, loggerInstance int) error {
	switch loggerInstance {
	case InstanceZapLogger:
		logger, err := newZapLogger(conf)

		if err != nil {
			return err
		}
		log = logger
		return nil
	case InstanceLogrusLogger:
		logger, err := newLogrusLogger(conf)

		if err != nil {
			return err
		}
		log = logger
		return nil
	default:
		return errInvalidLoggerInstance
	}
}

// Debug - Writes debug level information.
func Debug(args ...interface{}) {
	log.Debug(args)
}

// Debugf - Writes debug level information through the specified format.
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Info - Writes info level information.
func Info(args ...interface{}) {
	log.Info(args)
}

// Infof - Writes info level information through the specified format.
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warn - Writes warn level information.
func Warn(args ...interface{}) {
	log.Warn(args)
}

// Warnf - Writes warn level information through the specified format.
func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Error - Writes error level information.
func Error(args ...interface{}) {
	log.Error(args)
}

// Errorf - Writes error level information through the specified format.
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatal - Writes fatal level information.
func Fatal(args ...interface{}) {
	log.Fatal(args)
}

// Fatalf - Writes fatal level information through the specified format.
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Panic - Writes panic information.
func Panic(args ...interface{}) {
	log.Panic(args)
}

// Panicf - Writes panic information through the specified format.
func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

// WithFields - Returns a logger that contains the specified fields information.
func WithFields(keyValues Fields) Logger {
	return log.WithFields(keyValues)
}

// WithError - Returns a logger that contains the fields with specified error
func WithError(err error) Logger {
	return log.WithError(err)
}
