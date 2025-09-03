// Package common provides shared utilities and types for the ARK SDK.
//
// This package implements a custom logger with configurable log levels,
// color-coded output, and environment variable support for configuration.
package common

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Log level constants for ArkLogger.
//
// These constants define the severity levels used by the logging system.
// Higher numbers indicate more verbose output (Debug=4 is most verbose,
// Critical=0 is least verbose).
const (
	Debug    = 4
	Info     = 3
	Warning  = 2
	Error    = 1
	Critical = 0
	Unknown  = -1
)

// Environment variable names for logger configuration.
//
// These constants define the environment variables that can be used
// to configure the logger behavior at runtime.
const (
	LoggerStyle        = "LOGGER_STYLE"
	LogLevel           = "LOG_LEVEL"
	LoggerStyleDefault = "default"
)

// ArkLogger provides structured logging with configurable levels and output formatting.
//
// ArkLogger wraps the standard log.Logger and adds features like:
// - Configurable log levels (Debug, Info, Warning, Error, Critical)
// - Color-coded console output
// - Environment variable configuration
// - Verbose mode control
//
// The logger supports both static configuration and dynamic environment-based
// configuration for flexible deployment scenarios.
type ArkLogger struct {
	*log.Logger
	verbose                bool
	logLevel               int
	name                   string
	resolveLogLevelFromEnv bool
}

// NewArkLogger creates a new instance of ArkLogger with the specified configuration.
//
// This function initializes a new logger with the provided settings. The logger
// will output to stdout with timestamp formatting and can be configured for
// different verbosity levels and environment-based configuration.
//
// Parameters:
//   - name: The name/prefix for log messages
//   - level: The minimum log level (0=Critical, 1=Error, 2=Warning, 3=Info, 4=Debug)
//   - verbose: Whether to enable verbose output (false disables all logging)
//   - resolveLogLevelFromEnv: Whether to dynamically resolve log level from LOG_LEVEL env var
//
// Returns a configured ArkLogger instance ready for use.
//
// Example:
//
//	logger := NewArkLogger("myapp", Info, true, false)
//	logger.Info("Application started")
func NewArkLogger(name string, level int, verbose bool, resolveLogLevelFromEnv bool) *ArkLogger {
	return &ArkLogger{
		Logger:                 log.New(os.Stdout, name, log.LstdFlags),
		name:                   name,
		verbose:                verbose,
		logLevel:               level,
		resolveLogLevelFromEnv: resolveLogLevelFromEnv,
	}
}

// LogLevelFromEnv retrieves the log level from the LOG_LEVEL environment variable.
//
// This function reads the LOG_LEVEL environment variable and converts it to
// the corresponding integer log level. If the environment variable is not set
// or empty, it defaults to Critical level.
//
// Returns the integer log level corresponding to the environment variable value,
// or Critical (0) if the variable is not set or contains an invalid value.
//
// Example:
//
//	os.Setenv("LOG_LEVEL", "DEBUG")
//	level := LogLevelFromEnv() // Returns 4 (Debug)
func LogLevelFromEnv() int {
	logLevelStr := os.Getenv(LogLevel)
	if logLevelStr == "" {
		return Critical
	}
	return StrToLogLevel(logLevelStr)
}

// StrToLogLevel converts a string representation of a log level to its integer value.
//
// This function parses string log level names (case-insensitive) and returns
// the corresponding integer constant. Supported values are DEBUG, INFO, WARNING,
// ERROR, and CRITICAL. Any unrecognized value defaults to Critical.
//
// Parameters:
//   - logLevelStr: String representation of the log level (e.g., "DEBUG", "info", "Warning")
//
// Returns the integer constant corresponding to the log level, or Critical (0)
// for unrecognized input.
//
// Example:
//
//	level := StrToLogLevel("DEBUG")    // Returns 4
//	level := StrToLogLevel("invalid")  // Returns 0 (Critical)
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

// LogLevel returns the current effective log level of the logger.
//
// This method returns the log level that will be used for filtering log messages.
// If the logger is configured to resolve the level from environment variables,
// it will dynamically read the LOG_LEVEL environment variable. Otherwise,
// it returns the static log level set during logger creation.
//
// Returns the current log level as an integer (0=Critical, 1=Error, 2=Warning, 3=Info, 4=Debug).
func (l *ArkLogger) LogLevel() int {
	if l.resolveLogLevelFromEnv {
		return LogLevelFromEnv()
	}
	return l.logLevel
}

// SetVerbose sets the verbosity mode of the logger.
//
// When verbose is false, the logger will not output any messages regardless
// of the log level. This provides a master switch to disable all logging
// output from the logger instance.
//
// Parameters:
//   - value: true to enable verbose output, false to disable all output
func (l *ArkLogger) SetVerbose(value bool) {
	l.verbose = value
}

// Debug logs a debug message with green color formatting.
//
// Debug messages are only output when the logger is in verbose mode and
// the current log level is set to Debug (4) or higher. Debug messages
// are typically used for detailed diagnostic information.
//
// Parameters:
//   - msg: Format string for the log message
//   - v: Optional format arguments for the message string
//
// Example:
//
//	logger.Debug("Processing user %s with ID %d", username, userID)
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

// Info logs an informational message with green color formatting.
//
// Info messages are output when the logger is in verbose mode and
// the current log level is set to Info (3) or higher. Info messages
// are used for general application flow information.
//
// Parameters:
//   - msg: Format string for the log message
//   - v: Optional format arguments for the message string
//
// Example:
//
//	logger.Info("User %s logged in successfully", username)
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

// Warning logs a warning message with yellow color formatting.
//
// Warning messages are output when the logger is in verbose mode and
// the current log level is set to Warning (2) or higher. Warning messages
// indicate potentially problematic situations that don't prevent operation.
//
// Parameters:
//   - msg: Format string for the log message
//   - v: Optional format arguments for the message string
//
// Example:
//
//	logger.Warning("Rate limit approaching for user %s", username)
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

// Error logs an error message with red color formatting.
//
// Error messages are output when the logger is in verbose mode and
// the current log level is set to Error (1) or higher. Error messages
// indicate error conditions that should be investigated.
//
// Parameters:
//   - msg: Format string for the log message
//   - v: Optional format arguments for the message string
//
// Example:
//
//	logger.Error("Failed to connect to database: %v", err)
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

// Fatal logs a fatal error message with bright red color formatting and exits the program.
//
// Fatal messages are output when the logger is in verbose mode and
// the current log level is set to Critical (0) or higher. After logging
// the message, this method calls os.Exit(-1) to terminate the program.
// This should only be used for unrecoverable error conditions.
//
// Parameters:
//   - msg: Format string for the log message
//   - v: Optional format arguments for the message string
//
// Example:
//
//	logger.Fatal("Cannot start application: %v", err)
//	// Program terminates after this call
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

// GetLogger creates a new instance of ArkLogger with application-specific configuration.
//
// This is the primary factory function for creating logger instances. It handles
// environment variable resolution, logger style configuration, and sets up
// appropriate formatting. If logLevel is -1, the logger will dynamically
// resolve the log level from the LOG_LEVEL environment variable.
//
// Parameters:
//   - app: Application name used as the logger prefix
//   - logLevel: Static log level (0-4), or -1 to resolve from environment
//
// Returns a configured ArkLogger instance, or nil if an unsupported logger
// style is specified in the LOGGER_STYLE environment variable.
//
// Example:
//
//	logger := GetLogger("myapp", Info)        // Static Info level
//	envLogger := GetLogger("myapp", -1)       // Dynamic level from env
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
//
// This variable provides a package-level logger that can be used throughout
// the Ark SDK. It is configured to resolve its log level from the LOG_LEVEL
// environment variable and uses "ark-sdk" as its prefix.
//
// Example:
//
//	common.GlobalLogger.Info("SDK operation completed")
var GlobalLogger = GetLogger("ark-sdk", -1)
