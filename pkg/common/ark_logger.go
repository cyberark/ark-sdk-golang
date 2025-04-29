package common

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Log Levels for ArkLogger
const (
	Debug    = 4
	Info     = 3
	Warning  = 2
	Error    = 1
	Critical = 0
	Unknown  = -1
)

// Logger Environment Variables
const (
	LoggerStyle        = "LOGGER_STYLE"
	LogLevel           = "LOG_LEVEL"
	LoggerStyleDefault = "default"
)

// ArkLogger is a custom logger for the Ark SDK.
type ArkLogger struct {
	*log.Logger
	verbose                bool
	logLevel               int
	name                   string
	resolveLogLevelFromEnv bool
}

// NewArkLogger creates a new instance of ArkLogger with the specified parameters.
func NewArkLogger(name string, level int, verbose bool, resolveLogLevelFromEnv bool) *ArkLogger {
	return &ArkLogger{
		Logger:                 log.New(os.Stdout, name, log.LstdFlags),
		name:                   name,
		verbose:                verbose,
		logLevel:               level,
		resolveLogLevelFromEnv: resolveLogLevelFromEnv,
	}
}

// LogLevelFromEnv retrieves the log level from the environment variable.
func LogLevelFromEnv() int {
	logLevelStr := os.Getenv(LogLevel)
	if logLevelStr == "" {
		return Critical
	}
	return StrToLogLevel(logLevelStr)
}

// StrToLogLevel converts a string representation of a log level to its integer value.
func StrToLogLevel(logLevelStr string) int {
	switch strings.ToUpper(logLevelStr) {
	case "DEBUG":
		return Debug
	case "INFO":
		return Info
	case "WARNING":
		return Warning
	case "ERROR":
		return Error
	case "CRITICAL":
		return Critical
	default:
		return Critical
	}
}

// LogLevel returns the current log level of the logger.
func (l *ArkLogger) LogLevel() int {
	if l.resolveLogLevelFromEnv {
		return LogLevelFromEnv()
	}
	return l.logLevel
}

// SetVerbose sets the verbosity of the logger.
func (l *ArkLogger) SetVerbose(value bool) {
	l.verbose = value
}

// Debug logs a debug message if the logger is verbose and the log level is set to Debug or higher.
func (l *ArkLogger) Debug(msg string, v ...interface{}) {
	if !l.verbose {
		return
	}
	if l.LogLevel() < Debug {
		return
	}
	colorMsg := fmt.Sprintf("| DEBUG | \033[1;32m%s\033[0m", fmt.Sprintf(msg, v...))
	l.Logger.Println(colorMsg)
}

// Info logs an info message if the logger is verbose and the log level is set to Info or higher.
func (l *ArkLogger) Info(msg string, v ...interface{}) {
	if !l.verbose {
		return
	}
	if l.LogLevel() < Info {
		return
	}
	colorMsg := fmt.Sprintf("| INFO | \033[32m%s\033[0m", fmt.Sprintf(msg, v...))
	l.Logger.Println(colorMsg)
}

// Warning logs a warning message if the logger is verbose and the log level is set to Warning or higher.
func (l *ArkLogger) Warning(msg string, v ...interface{}) {
	if !l.verbose {
		return
	}
	if l.LogLevel() < Warning {
		return
	}
	colorMsg := fmt.Sprintf("| WARNING | \033[33m%s\033[0m", fmt.Sprintf(msg, v...))
	l.Logger.Println(colorMsg)
}

// Error logs an error message if the logger is verbose and the log level is set to Error or higher.
func (l *ArkLogger) Error(msg string, v ...interface{}) {
	if !l.verbose {
		return
	}
	if l.LogLevel() < Error {
		return
	}
	colorMsg := fmt.Sprintf("| ERROR | \033[31m%s\033[0m", fmt.Sprintf(msg, v...))
	l.Logger.Println(colorMsg)
}

// Fatal logs a fatal message if the logger is verbose and the log level is set to Critical or higher, then exits the program.
func (l *ArkLogger) Fatal(msg string, v ...interface{}) {
	if !l.verbose {
		return
	}
	if l.LogLevel() < Critical {
		return
	}
	colorMsg := fmt.Sprintf("| FATAL | \033[1;31m%s\033[0m", fmt.Sprintf(msg, v...))
	l.Logger.Println(colorMsg)
	os.Exit(-1)
}

// GetLogger creates a new instance of ArkLogger with the specified application name and log level.
func GetLogger(app string, logLevel int) *ArkLogger {
	resolveLogLevelFromEnv := false
	if logLevel == -1 {
		resolveLogLevelFromEnv = true
		logLevel = LogLevelFromEnv()
	}
	envLoggerStyle := os.Getenv(LoggerStyle)
	if envLoggerStyle == "" {
		envLoggerStyle = LoggerStyleDefault
	}
	loggerStyle := strings.ToLower(envLoggerStyle)
	if loggerStyle == LoggerStyleDefault {
		logFormat := "%s | "
		logger := NewArkLogger(app, logLevel, true, resolveLogLevelFromEnv)
		logger.SetFlags(log.LstdFlags)
		logger.SetPrefix(fmt.Sprintf(logFormat, app))
		return logger
	}
	return nil
}

// GlobalLogger is the global logger instance for the Ark SDK.
var GlobalLogger = GetLogger("ark-sdk", -1)
