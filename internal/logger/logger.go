package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var (
	// Log is the default logger instance
	Log *logrus.Logger
)

// Config holds configuration for the logger
type Config struct {
	LogLevel      string
	LogFile       string
	MaxSize       int
	MaxBackups    int
	MaxAge        int
	Compress      bool
	ConsoleOutput bool
}

// Custom formatter that includes caller info in the main format
type CustomFormatter struct {
	*logrus.TextFormatter
}

// Format renders a single log entry
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Get caller info from the entry data if it exists
	caller := "unknown"
	if callerInfo, exists := entry.Data["caller_info"]; exists {
		caller = callerInfo.(string)
		// Remove it from the data so it doesn't appear twice
		delete(entry.Data, "caller_info")
	}

	// Create the custom format: LEVEL [timestamp] caller: message
	levelText := strings.ToUpper(entry.Level.String())
	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	// Build the log line
	logLine := fmt.Sprintf("%-7s[%s] %s: %s",
		levelText,
		timestamp,
		caller,
		entry.Message,
	)

	// Add fields if any exist
	if len(entry.Data) > 0 {
		for key, value := range entry.Data {
			logLine += fmt.Sprintf(" %s=%v", key, value)
		}
	}

	logLine += "\n"
	return []byte(logLine), nil
}

// Initialize sets up the logger with configuration
func Initialize(config Config) {
	if Log == nil {
		Log = logrus.New()
	}

	// Set log level
	level, err := logrus.ParseLevel(strings.ToLower(config.LogLevel))
	if err != nil {
		level = logrus.InfoLevel
	}
	Log.SetLevel(level)

	// Use our custom formatter
	Log.SetFormatter(&CustomFormatter{
		TextFormatter: &logrus.TextFormatter{
			ForceColors:            true,
			DisableColors:          false,
			FullTimestamp:          true,
			DisableLevelTruncation: false,
			TimestampFormat:        "2006-01-02 15:04:05",
		},
	})

	// 🔧 IMPORTANT: Disable built-in caller reporting since we handle it manually
	Log.SetReportCaller(false)

	outputs := []io.Writer{}

	// Set console output if enabled
	if config.ConsoleOutput {
		outputs = append(outputs, os.Stdout)
	}

	// Set file output if a log file is specified
	if config.LogFile != "" {
		// Ensure log directory exists
		logDir := filepath.Dir(config.LogFile)
		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			if err := os.MkdirAll(logDir, 0755); err != nil {
				fmt.Printf("Failed to create log directory: %v\n", err)
			}
		}

		// Configure log rotation
		logRotator := &lumberjack.Logger{
			Filename:   config.LogFile,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}
		outputs = append(outputs, logRotator)
	}

	// Set output to both console and file if needed
	if len(outputs) > 0 {
		multiWriter := io.MultiWriter(outputs...)
		Log.SetOutput(multiWriter)
	}
}

// Helper function to get caller info and create log entry
func getLogEntryWithCaller() *logrus.Entry {
	if Log == nil {
		return nil
	}

	// Get caller info, skipping this function and the wrapper function
	_, file, line, ok := runtime.Caller(2)
	if ok {
		fileName := filepath.Base(file)
		callerInfo := fmt.Sprintf("%s:%d", fileName, line)
		return Log.WithField("caller_info", callerInfo)
	}

	return Log.WithField("caller_info", "unknown")
}

// 🔧 UPDATED: Convenience methods with correct caller reporting
func Debug(args ...interface{}) {
	if entry := getLogEntryWithCaller(); entry != nil {
		entry.Debug(args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if entry := getLogEntryWithCaller(); entry != nil {
		entry.Debugf(format, args...)
	}
}

func Info(args ...interface{}) {
	if entry := getLogEntryWithCaller(); entry != nil {
		entry.Info(args...)
	}
}

func Infof(format string, args ...interface{}) {
	if entry := getLogEntryWithCaller(); entry != nil {
		entry.Infof(format, args...)
	}
}

func Warn(args ...interface{}) {
	if entry := getLogEntryWithCaller(); entry != nil {
		entry.Warn(args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if entry := getLogEntryWithCaller(); entry != nil {
		entry.Warnf(format, args...)
	}
}

func Error(args ...interface{}) {
	if entry := getLogEntryWithCaller(); entry != nil {
		entry.Error(args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if entry := getLogEntryWithCaller(); entry != nil {
		entry.Errorf(format, args...)
	}
}

func Fatal(args ...interface{}) {
	if entry := getLogEntryWithCaller(); entry != nil {
		entry.Fatal(args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	if entry := getLogEntryWithCaller(); entry != nil {
		entry.Fatalf(format, args...)
	}
}

// 🔧 BONUS: Structured logging methods that preserve the caller format
func InfoWithFields(message string, fields map[string]interface{}) {
	if entry := getLogEntryWithCaller(); entry != nil {
		entry.WithFields(fields).Info(message)
	}
}

func DebugWithFields(message string, fields map[string]interface{}) {
	if entry := getLogEntryWithCaller(); entry != nil {
		entry.WithFields(fields).Debug(message)
	}
}

func ErrorWithFields(message string, fields map[string]interface{}) {
	if entry := getLogEntryWithCaller(); entry != nil {
		entry.WithFields(fields).Error(message)
	}
}
