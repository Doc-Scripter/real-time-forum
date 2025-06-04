package logging

import (
	"log"
	"os"
	"path/filepath"
)

// InitializeLogger creates a logging system with both file and terminal output capabilities.
// It sets up a logs directory, creates a log file, and initializes two loggers - one for
// file logging and another for terminal output.
func InitializeLogger() {
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	logFile, err := os.OpenFile(
		filepath.Join(logsDir, "application.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Fatal("Failed to create log file:", err)
	}

	fileLogger := log.New(logFile, "", log.LstdFlags)
	terminalLogger = log.New(os.Stdout, "", log.LstdFlags)

	logger = fileLogger
	terminalLogger.Println("Logger initialized successfully")
}

// logger is the package-level variable for file logging
var logger *log.Logger

// terminalLogger is the package-level variable for terminal logging
var terminalLogger *log.Logger

// Log writes the formatted message to the log file
func Log(format string, v ...any) {
	logger.Printf(format, v...)
}

// TerminalLog writes the formatted message to the terminal
func TerminalLog(format string, v ...interface{}) {
	terminalLogger.Printf(format, v...)
}
