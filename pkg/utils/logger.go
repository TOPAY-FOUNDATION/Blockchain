package utils

import (
	"log"
	"os"
)

// Logger is a custom logger
type Logger struct {
	*log.Logger
}

// NewLogger creates a new logger instance
func NewLogger(prefix string) *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, prefix, log.LstdFlags),
	}
}
