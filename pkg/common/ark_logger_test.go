package common

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestStrToLogLevel(t *testing.T) {
	tests := []struct {
		name          string
		logLevelStr   string
		expectedLevel int
	}{
		{
			name:          "success_debug_uppercase",
			logLevelStr:   "DEBUG",
			expectedLevel: Debug,
		},
		{
			name:          "success_debug_lowercase",
			logLevelStr:   "debug",
			expectedLevel: Debug,
		},
		{
			name:          "success_info_uppercase",
			logLevelStr:   "INFO",
			expectedLevel: Info,
		},
		{
			name:          "success_info_mixed_case",
			logLevelStr:   "InFo",
			expectedLevel: Info,
		},
		{
			name:          "success_warning_uppercase",
			logLevelStr:   "WARNING",
			expectedLevel: Warning,
		},
		{
			name:          "success_error_uppercase",
			logLevelStr:   "ERROR",
			expectedLevel: Error,
		},
		{
			name:          "success_critical_uppercase",
			logLevelStr:   "CRITICAL",
			expectedLevel: Critical,
		},
		{
			name:          "default_unknown_string",
			logLevelStr:   "UNKNOWN",
			expectedLevel: Critical,
		},
		{
			name:          "default_empty_string",
			logLevelStr:   "",
			expectedLevel: Critical,
		},
		{
			name:          "default_invalid_string",
			logLevelStr:   "invalid_level",
			expectedLevel: Critical,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := StrToLogLevel(tt.logLevelStr)

			if result != tt.expectedLevel {
				t.Errorf("Expected log level %d, got %d", tt.expectedLevel, result)
			}
		})
	}
}

func TestLogLevelFromEnv(t *testing.T) {
	tests := []struct {
		name          string
		envValue      string
		setEnv        bool
		expectedLevel int
	}{
		{
			name:          "success_debug_from_env",
			envValue:      "DEBUG",
			setEnv:        true,
			expectedLevel: Debug,
		},
		{
			name:          "success_info_from_env",
			envValue:      "INFO",
			setEnv:        true,
			expectedLevel: Info,
		},
		{
			name:          "success_warning_from_env",
			envValue:      "WARNING",
			setEnv:        true,
			expectedLevel: Warning,
		},
		{
			name:          "success_error_from_env",
			envValue:      "ERROR",
			setEnv:        true,
			expectedLevel: Error,
		},
		{
			name:          "success_critical_from_env",
			envValue:      "CRITICAL",
			setEnv:        true,
			expectedLevel: Critical,
		},
		{
			name:          "default_empty_env_var",
			envValue:      "",
			setEnv:        true,
			expectedLevel: Critical,
		},
		{
			name:          "default_unset_env_var",
			setEnv:        false,
			expectedLevel: Critical,
		},
		{
			name:          "default_invalid_env_value",
			envValue:      "INVALID",
			setEnv:        true,
			expectedLevel: Critical,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			originalValue := os.Getenv(LogLevel)
			defer func() {
				if originalValue != "" {
					os.Setenv(LogLevel, originalValue)
				} else {
					os.Unsetenv(LogLevel)
				}
			}()

			if tt.setEnv {
				os.Setenv(LogLevel, tt.envValue)
			} else {
				os.Unsetenv(LogLevel)
			}

			result := LogLevelFromEnv()

			if result != tt.expectedLevel {
				t.Errorf("Expected log level %d, got %d", tt.expectedLevel, result)
			}
		})
	}
}

func TestNewArkLogger(t *testing.T) {
	tests := []struct {
		name                   string
		loggerName             string
		level                  int
		verbose                bool
		resolveLogLevelFromEnv bool
		validateFunc           func(t *testing.T, logger *ArkLogger)
	}{
		{
			name:                   "success_basic_logger",
			loggerName:             "test-app",
			level:                  Info,
			verbose:                true,
			resolveLogLevelFromEnv: false,
			validateFunc: func(t *testing.T, logger *ArkLogger) {
				if logger.name != "test-app" {
					t.Errorf("Expected name 'test-app', got '%s'", logger.name)
				}
				if logger.logLevel != Info {
					t.Errorf("Expected log level %d, got %d", Info, logger.logLevel)
				}
				if !logger.verbose {
					t.Error("Expected verbose to be true")
				}
				if logger.resolveLogLevelFromEnv {
					t.Error("Expected resolveLogLevelFromEnv to be false")
				}
			},
		},
		{
			name:                   "success_env_resolver_logger",
			loggerName:             "env-app",
			level:                  Debug,
			verbose:                false,
			resolveLogLevelFromEnv: true,
			validateFunc: func(t *testing.T, logger *ArkLogger) {
				if logger.name != "env-app" {
					t.Errorf("Expected name 'env-app', got '%s'", logger.name)
				}
				if logger.logLevel != Debug {
					t.Errorf("Expected log level %d, got %d", Debug, logger.logLevel)
				}
				if logger.verbose {
					t.Error("Expected verbose to be false")
				}
				if !logger.resolveLogLevelFromEnv {
					t.Error("Expected resolveLogLevelFromEnv to be true")
				}
			},
		},
		{
			name:                   "success_critical_level",
			loggerName:             "critical-app",
			level:                  Critical,
			verbose:                true,
			resolveLogLevelFromEnv: false,
			validateFunc: func(t *testing.T, logger *ArkLogger) {
				if logger.logLevel != Critical {
					t.Errorf("Expected log level %d, got %d", Critical, logger.logLevel)
				}
			},
		},
		{
			name:                   "success_empty_name",
			loggerName:             "",
			level:                  Warning,
			verbose:                true,
			resolveLogLevelFromEnv: false,
			validateFunc: func(t *testing.T, logger *ArkLogger) {
				if logger.name != "" {
					t.Errorf("Expected empty name, got '%s'", logger.name)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger := NewArkLogger(tt.loggerName, tt.level, tt.verbose, tt.resolveLogLevelFromEnv)

			if logger == nil {
				t.Fatal("Expected logger to be created, got nil")
			}

			if logger.Logger == nil {
				t.Error("Expected embedded Logger to be initialized")
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, logger)
			}
		})
	}
}

func TestArkLogger_LogLevel(t *testing.T) {
	tests := []struct {
		name                   string
		staticLevel            int
		resolveLogLevelFromEnv bool
		envValue               string
		setEnv                 bool
		expectedLevel          int
	}{
		{
			name:                   "success_static_level",
			staticLevel:            Info,
			resolveLogLevelFromEnv: false,
			expectedLevel:          Info,
		},
		{
			name:                   "success_env_level_debug",
			staticLevel:            Warning,
			resolveLogLevelFromEnv: true,
			envValue:               "DEBUG",
			setEnv:                 true,
			expectedLevel:          Debug,
		},
		{
			name:                   "success_env_level_critical",
			staticLevel:            Info,
			resolveLogLevelFromEnv: true,
			envValue:               "CRITICAL",
			setEnv:                 true,
			expectedLevel:          Critical,
		},
		{
			name:                   "success_env_unset_defaults_critical",
			staticLevel:            Info,
			resolveLogLevelFromEnv: true,
			setEnv:                 false,
			expectedLevel:          Critical,
		},
		{
			name:                   "success_env_invalid_defaults_critical",
			staticLevel:            Info,
			resolveLogLevelFromEnv: true,
			envValue:               "INVALID",
			setEnv:                 true,
			expectedLevel:          Critical,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			originalValue := os.Getenv(LogLevel)
			defer func() {
				if originalValue != "" {
					os.Setenv(LogLevel, originalValue)
				} else {
					os.Unsetenv(LogLevel)
				}
			}()

			if tt.setEnv {
				os.Setenv(LogLevel, tt.envValue)
			} else {
				os.Unsetenv(LogLevel)
			}

			logger := NewArkLogger("test", tt.staticLevel, true, tt.resolveLogLevelFromEnv)
			result := logger.LogLevel()

			if result != tt.expectedLevel {
				t.Errorf("Expected log level %d, got %d", tt.expectedLevel, result)
			}
		})
	}
}

func TestArkLogger_SetVerbose(t *testing.T) {
	tests := []struct {
		name         string
		initialValue bool
		newValue     bool
	}{
		{
			name:         "success_set_true_from_false",
			initialValue: false,
			newValue:     true,
		},
		{
			name:         "success_set_false_from_true",
			initialValue: true,
			newValue:     false,
		},
		{
			name:         "success_set_true_from_true",
			initialValue: true,
			newValue:     true,
		},
		{
			name:         "success_set_false_from_false",
			initialValue: false,
			newValue:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger := NewArkLogger("test", Info, tt.initialValue, false)
			logger.SetVerbose(tt.newValue)

			if logger.verbose != tt.newValue {
				t.Errorf("Expected verbose to be %v, got %v", tt.newValue, logger.verbose)
			}
		})
	}
}

func TestArkLogger_Debug(t *testing.T) {
	tests := []struct {
		name             string
		verbose          bool
		logLevel         int
		msg              string
		args             []interface{}
		expectOutput     bool
		expectedContains []string
	}{
		{
			name:             "success_debug_output",
			verbose:          true,
			logLevel:         Debug,
			msg:              "Debug message %s",
			args:             []interface{}{"test"},
			expectOutput:     true,
			expectedContains: []string{"DEBUG", "Debug message test"},
		},
		{
			name:         "success_no_output_verbose_false",
			verbose:      false,
			logLevel:     Debug,
			msg:          "Debug message",
			expectOutput: false,
		},
		{
			name:         "success_no_output_level_too_low",
			verbose:      true,
			logLevel:     Info,
			msg:          "Debug message",
			expectOutput: false,
		},
		{
			name:         "success_no_output_critical_level",
			verbose:      true,
			logLevel:     Critical,
			msg:          "Debug message",
			expectOutput: false,
		},
		{
			name:             "success_debug_no_args",
			verbose:          true,
			logLevel:         Debug,
			msg:              "Simple debug message",
			expectOutput:     true,
			expectedContains: []string{"DEBUG", "Simple debug message"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			logger := NewArkLogger("test", tt.logLevel, tt.verbose, false)
			logger.Logger = log.New(&buf, "test", log.LstdFlags)

			logger.Debug(tt.msg, tt.args...)

			output := buf.String()
			if tt.expectOutput {
				if output == "" {
					t.Error("Expected output, got empty string")
				}
				for _, expected := range tt.expectedContains {
					if !strings.Contains(output, expected) {
						t.Errorf("Expected output to contain '%s', got '%s'", expected, output)
					}
				}
			} else {
				if output != "" {
					t.Errorf("Expected no output, got '%s'", output)
				}
			}
		})
	}
}

func TestArkLogger_Info(t *testing.T) {
	tests := []struct {
		name             string
		verbose          bool
		logLevel         int
		msg              string
		args             []interface{}
		expectOutput     bool
		expectedContains []string
	}{
		{
			name:             "success_info_output",
			verbose:          true,
			logLevel:         Info,
			msg:              "Info message %d",
			args:             []interface{}{42},
			expectOutput:     true,
			expectedContains: []string{"INFO", "Info message 42"},
		},
		{
			name:             "success_info_debug_level",
			verbose:          true,
			logLevel:         Debug,
			msg:              "Info message",
			expectOutput:     true,
			expectedContains: []string{"INFO"},
		},
		{
			name:         "success_no_output_verbose_false",
			verbose:      false,
			logLevel:     Info,
			msg:          "Info message",
			expectOutput: false,
		},
		{
			name:         "success_no_output_level_too_low",
			verbose:      true,
			logLevel:     Warning,
			msg:          "Info message",
			expectOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			logger := NewArkLogger("test", tt.logLevel, tt.verbose, false)
			logger.Logger = log.New(&buf, "test", log.LstdFlags)

			logger.Info(tt.msg, tt.args...)

			output := buf.String()
			if tt.expectOutput {
				if output == "" {
					t.Error("Expected output, got empty string")
				}
				for _, expected := range tt.expectedContains {
					if !strings.Contains(output, expected) {
						t.Errorf("Expected output to contain '%s', got '%s'", expected, output)
					}
				}
			} else {
				if output != "" {
					t.Errorf("Expected no output, got '%s'", output)
				}
			}
		})
	}
}

func TestArkLogger_Warning(t *testing.T) {
	tests := []struct {
		name             string
		verbose          bool
		logLevel         int
		msg              string
		args             []interface{}
		expectOutput     bool
		expectedContains []string
	}{
		{
			name:             "success_warning_output",
			verbose:          true,
			logLevel:         Warning,
			msg:              "Warning message %s",
			args:             []interface{}{"alert"},
			expectOutput:     true,
			expectedContains: []string{"WARNING", "Warning message alert"},
		},
		{
			name:             "success_warning_higher_level",
			verbose:          true,
			logLevel:         Debug,
			msg:              "Warning message",
			expectOutput:     true,
			expectedContains: []string{"WARNING"},
		},
		{
			name:         "success_no_output_verbose_false",
			verbose:      false,
			logLevel:     Warning,
			msg:          "Warning message",
			expectOutput: false,
		},
		{
			name:         "success_no_output_level_too_low",
			verbose:      true,
			logLevel:     Error,
			msg:          "Warning message",
			expectOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			logger := NewArkLogger("test", tt.logLevel, tt.verbose, false)
			logger.Logger = log.New(&buf, "test", log.LstdFlags)

			logger.Warning(tt.msg, tt.args...)

			output := buf.String()
			if tt.expectOutput {
				if output == "" {
					t.Error("Expected output, got empty string")
				}
				for _, expected := range tt.expectedContains {
					if !strings.Contains(output, expected) {
						t.Errorf("Expected output to contain '%s', got '%s'", expected, output)
					}
				}
			} else {
				if output != "" {
					t.Errorf("Expected no output, got '%s'", output)
				}
			}
		})
	}
}

func TestArkLogger_Error(t *testing.T) {
	tests := []struct {
		name             string
		verbose          bool
		logLevel         int
		msg              string
		args             []interface{}
		expectOutput     bool
		expectedContains []string
	}{
		{
			name:             "success_error_output",
			verbose:          true,
			logLevel:         Error,
			msg:              "Error message %v",
			args:             []interface{}{"failed"},
			expectOutput:     true,
			expectedContains: []string{"ERROR", "Error message failed"},
		},
		{
			name:             "success_error_higher_level",
			verbose:          true,
			logLevel:         Debug,
			msg:              "Error message",
			expectOutput:     true,
			expectedContains: []string{"ERROR"},
		},
		{
			name:         "success_no_output_verbose_false",
			verbose:      false,
			logLevel:     Error,
			msg:          "Error message",
			expectOutput: false,
		},
		{
			name:         "success_no_output_level_too_low",
			verbose:      true,
			logLevel:     Critical,
			msg:          "Error message",
			expectOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			logger := NewArkLogger("test", tt.logLevel, tt.verbose, false)
			logger.Logger = log.New(&buf, "test", log.LstdFlags)

			logger.Error(tt.msg, tt.args...)

			output := buf.String()
			if tt.expectOutput {
				if output == "" {
					t.Error("Expected output, got empty string")
				}
				for _, expected := range tt.expectedContains {
					if !strings.Contains(output, expected) {
						t.Errorf("Expected output to contain '%s', got '%s'", expected, output)
					}
				}
			} else {
				if output != "" {
					t.Errorf("Expected no output, got '%s'", output)
				}
			}
		})
	}
}

// Note: TestArkLogger_Fatal is not included as it calls os.Exit(-1) which would terminate the test process.
// In a real scenario, this would require special handling with a separate process or test isolation.

func TestGetLogger(t *testing.T) {
	tests := []struct {
		name         string
		app          string
		logLevel     int
		envStyle     string
		setStyleEnv  bool
		expectedNil  bool
		validateFunc func(t *testing.T, logger *ArkLogger)
	}{
		{
			name:        "success_default_style",
			app:         "test-app",
			logLevel:    Info,
			expectedNil: false,
			validateFunc: func(t *testing.T, logger *ArkLogger) {
				if logger.name != "test-app" {
					t.Errorf("Expected name 'test-app', got '%s'", logger.name)
				}
				if logger.logLevel != Info {
					t.Errorf("Expected log level %d, got %d", Info, logger.logLevel)
				}
				if !logger.verbose {
					t.Error("Expected verbose to be true")
				}
			},
		},
		{
			name:        "success_env_log_level",
			app:         "env-app",
			logLevel:    -1,
			expectedNil: false,
			validateFunc: func(t *testing.T, logger *ArkLogger) {
				if !logger.resolveLogLevelFromEnv {
					t.Error("Expected resolveLogLevelFromEnv to be true")
				}
			},
		},
		{
			name:        "success_explicit_default_style",
			app:         "styled-app",
			logLevel:    Warning,
			envStyle:    "default",
			setStyleEnv: true,
			expectedNil: false,
			validateFunc: func(t *testing.T, logger *ArkLogger) {
				if logger.logLevel != Warning {
					t.Errorf("Expected log level %d, got %d", Warning, logger.logLevel)
				}
			},
		},
		{
			name:        "success_unknown_style_returns_nil",
			app:         "unknown-app",
			logLevel:    Debug,
			envStyle:    "unknown",
			setStyleEnv: true,
			expectedNil: true,
		},
		{
			name:        "success_empty_app_name",
			app:         "",
			logLevel:    Critical,
			expectedNil: false,
			validateFunc: func(t *testing.T, logger *ArkLogger) {
				if logger.name != "" {
					t.Errorf("Expected empty name, got '%s'", logger.name)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			originalStyle := os.Getenv(LoggerStyle)
			defer func() {
				if originalStyle != "" {
					os.Setenv(LoggerStyle, originalStyle)
				} else {
					os.Unsetenv(LoggerStyle)
				}
			}()

			if tt.setStyleEnv {
				os.Setenv(LoggerStyle, tt.envStyle)
			} else {
				os.Unsetenv(LoggerStyle)
			}

			logger := GetLogger(tt.app, tt.logLevel)

			if tt.expectedNil {
				if logger != nil {
					t.Error("Expected nil logger, got non-nil")
				}
				return
			}

			if logger == nil {
				t.Fatal("Expected logger to be created, got nil")
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, logger)
			}
		})
	}
}

func TestGlobalLogger(t *testing.T) {
	// This test validates that the GlobalLogger is properly initialized
	if GlobalLogger == nil {
		t.Error("Expected GlobalLogger to be initialized, got nil")
	}

	if GlobalLogger.name != "ark-sdk" {
		t.Errorf("Expected GlobalLogger name to be 'ark-sdk', got '%s'", GlobalLogger.name)
	}

	if !GlobalLogger.resolveLogLevelFromEnv {
		t.Error("Expected GlobalLogger to resolve log level from environment")
	}
}
