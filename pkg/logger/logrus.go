// Package logger - Logger for Logus
package logger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// logrusLogEntry - defines the log entry for logrus
type logrusLogEntry struct {
	entry *logrus.Entry
}

// logrusLogger - defines the logger for logrus
type logrusLogger struct {
	logger *logrus.Logger
}

// ===== [ Implementations ] =====

// Debug - Implements of Logger's Debug on LogrusLogger
func (ll *logrusLogger) Debug(args ...interface{}) {
	ll.logger.Debug(args...)
}

// Debugf - Implements of Logger's Debugf on LogrusLogger
func (ll *logrusLogger) Debugf(format string, args ...interface{}) {
	ll.logger.Debugf(format, args...)
}

// Info - Implements of Logger's Info on LogrusLogger
func (ll *logrusLogger) Info(args ...interface{}) {
	ll.logger.Info(args...)
}

// Infof - Implements of Logger's Infof on LogrusLogger
func (ll *logrusLogger) Infof(format string, args ...interface{}) {
	ll.logger.Infof(format, args...)
}

// Warn - Implements of Logger's Warn on LogrusLogger
func (ll *logrusLogger) Warn(args ...interface{}) {
	ll.logger.Warn(args...)
}

// Warnf - Implements of Logger's Warnf on LogrusLogger
func (ll *logrusLogger) Warnf(format string, args ...interface{}) {
	ll.logger.Warnf(format, args...)
}

// Error - Implements of Logger's Error on LogrusLogger
func (ll *logrusLogger) Error(args ...interface{}) {
	ll.logger.Error(args...)
}

// Errorf - Implements of Logger's Errorf on LogrusLogger
func (ll *logrusLogger) Errorf(format string, args ...interface{}) {
	ll.logger.Errorf(format, args...)
}

// Fatal - Implements of Logger's Fatal on LogrusLogger
func (ll *logrusLogger) Fatal(args ...interface{}) {
	ll.logger.Fatal(args...)
}

// Fatalf - Implements of Logger's Fatalf on LogrusLogger
func (ll *logrusLogger) Fatalf(format string, args ...interface{}) {
	ll.logger.Fatalf(format, args...)
}

// Panic - Implements of Logger's Panic on LogrusLogger
func (ll *logrusLogger) Panic(args ...interface{}) {
	ll.logger.Fatal(args...)
}

// Panicf - Implements of Logger's Panicf on LogrusLogger
func (ll *logrusLogger) Panicf(format string, args ...interface{}) {
	ll.logger.Fatalf(format, args...)
}

// WithFields - Implements of Logger's WithFields on LogrusLogger
func (ll *logrusLogger) WithFields(fields Fields) Logger {
	return &logrusLogEntry{
		entry: ll.logger.WithFields(convertToLogrusFields(fields)),
	}
}

// WithError - Implements of Loggers WithError on LogrusLogger
func (ll *logrusLogger) WithError(err error) Logger {
	return &logrusLogEntry{
		entry: ll.logger.WithError(err),
	}
}

// Debug - Implements of Logger's Debug on LogrusLogEntry
func (le *logrusLogEntry) Debug(args ...interface{}) {
	le.entry.Debug(args...)
}

// Debugf - Implements of Logger's Debugf on LogrusLogEntry
func (le *logrusLogEntry) Debugf(format string, args ...interface{}) {
	le.entry.Debugf(format, args...)
}

// Info - Implements of Logger's Info on LogrusLogEntry
func (le *logrusLogEntry) Info(args ...interface{}) {
	le.entry.Info(args...)
}

// Infof - Implements of Logger's Infof on LogrusLogEntry
func (le *logrusLogEntry) Infof(format string, args ...interface{}) {
	le.entry.Infof(format, args...)
}

// Warn - Implements of Logger's Warn on LogrusLogEntry
func (le *logrusLogEntry) Warn(args ...interface{}) {
	le.entry.Warn(args...)
}

// Warnf - Implements of Logger's Warnf on LogrusLogEntry
func (le *logrusLogEntry) Warnf(format string, args ...interface{}) {
	le.entry.Warnf(format, args...)
}

// Error - Implements of Logger's Error on LogrusLogEntry
func (le *logrusLogEntry) Error(args ...interface{}) {
	le.entry.Error(args...)
}

// Errorf - Implements of Logger's Errorf on LogrusLogEntry
func (le *logrusLogEntry) Errorf(format string, args ...interface{}) {
	le.entry.Errorf(format, args...)
}

// Fatal - Implements of Logger's Fatal on LogrusLogEntry
func (le *logrusLogEntry) Fatal(args ...interface{}) {
	le.entry.Fatal(args...)
}

// Fatalf - Implements of Logger's Fatalf on LogrusLogEntry
func (le *logrusLogEntry) Fatalf(format string, args ...interface{}) {
	le.entry.Fatalf(format, args...)
}

// Panic - Implements of Logger's Panic on LogrusLogEntry
func (le *logrusLogEntry) Panic(args ...interface{}) {
	le.entry.Fatal(args...)
}

// Panicf - Implements of Logger's Panicf on LogrusLogEntry
func (le *logrusLogEntry) Panicf(format string, args ...interface{}) {
	le.entry.Fatalf(format, args...)
}

// WithFields - Implements of Logger's WithFields on LogrusLogEntry
func (le *logrusLogEntry) WithFields(fields Fields) Logger {
	return &logrusLogEntry{
		entry: le.entry.WithFields(convertToLogrusFields(fields)),
	}
}

// WithError - Implements of Loggers WithError on LogrusLogEntry
func (le *logrusLogEntry) WithError(err error) Logger {
	return &logrusLogEntry{
		entry: le.entry.WithField(logrus.ErrorKey, err),
	}
}

// ===== [ Private Functions ] =====

// getFormatter - Returns an formatter for logrus
func getFormatter(isJSON bool) logrus.Formatter {
	if isJSON {
		return &logrus.JSONFormatter{}
	}
	return &logrus.TextFormatter{FullTimestamp: true, DisableLevelTruncation: true}
}

// convertToLogrusFields - Converts the fields to logrus fileds
func convertToLogrusFields(fields Fields) logrus.Fields {
	logrusFields := logrus.Fields{}
	for idx, val := range fields {
		logrusFields[idx] = val
	}
	return logrusFields
}

// newLogrusLogger - Returns an instance of logrus
func newLogrusLogger(conf Config) (Logger, error) {
	logLevel := conf.ConsoleLevel
	if logLevel == "" {
		logLevel = conf.FileLevel
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	stdOutHandler := os.Stdout
	fileHandler := &lumberjack.Logger{
		Filename: conf.FileLocation,
		MaxSize:  100,
		Compress: true,
		MaxAge:   28,
	}
	lLogger := &logrus.Logger{
		Out:       stdOutHandler,
		Formatter: getFormatter(conf.ConsoleJSONFormat),
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
	}

	if conf.EnableConsole && conf.EnableFile {
		lLogger.SetOutput(io.MultiWriter(stdOutHandler, fileHandler))
	} else {
		if conf.EnableFile {
			lLogger.SetOutput(fileHandler)
			lLogger.SetFormatter(getFormatter(conf.FileJSONFormat))
		}
	}

	return &logrusLogger{
		logger: lLogger,
	}, nil
}

// ===== [ Public Functions ] =====
