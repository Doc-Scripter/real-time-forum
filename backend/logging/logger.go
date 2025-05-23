package logging

import (
	"log"
	"os"
	"path/filepath"
)

// InitializeLogger sets up the logger with file output
func InitializeLogger() {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	// Create log file
	logFile, err := os.OpenFile(
		filepath.Join(logsDir, "application.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Fatal("Failed to create log file:", err)
	}

	// Create a logger that writes to file
	fileLogger := log.New(logFile, "", log.LstdFlags|log.Lshortfile)

	// Create a logger that writes to terminal
	terminalLogger = log.New(os.Stdout, "", log.LstdFlags)

	// Store both loggers in package-level variables
	logger = fileLogger
	terminalLogger.Println("Logger initialized successfully")
}

// Logger that writes to file
var logger *log.Logger

// Logger that writes to terminal
var terminalLogger *log.Logger

// Log writes to file only
func Log(format string, v ...any) {
	logger.Printf(format, v...)
}

// TerminalLog writes to terminal only
func TerminalLog(format string, v ...interface{}) {
	terminalLogger.Printf(format, v...)
}
