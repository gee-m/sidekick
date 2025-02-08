package logging

import (
	"log"
	"os"
	"time"
)

type Level string

const (
	InfoLevel  Level = "INFO"
	WarnLevel  Level = "WARN"
	ErrorLevel Level = "ERROR"
	DebugLevel Level = "DEBUG"
)

type Logger struct {
	logger *log.Logger
}

type LogEntry struct {
	Timestamp time.Time      `json:"timestamp"`
	Level     Level          `json:"level"`
	Message   string         `json:"message"`
	Data      map[string]any `json:"data,omitempty"`
}

func New() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", 0),
	}
}

func (l *Logger) log(level Level, msg string, data map[string]any) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Data:      data,
	}

	l.logger.Printf("[%s] %s | %s | %v\n",
		entry.Timestamp.Format(time.RFC3339),
		entry.Level,
		entry.Message,
		entry.Data,
	)
}

func (l *Logger) Info(msg string, data map[string]any) {
	l.log(InfoLevel, msg, data)
}

func (l *Logger) Warn(msg string, data map[string]any) {
	l.log(WarnLevel, msg, data)
}

func (l *Logger) Error(msg string, err error, data map[string]any) {
	if data == nil {
		data = make(map[string]any)
	}
	data["error"] = err.Error()
	l.log(ErrorLevel, msg, data)
}

func (l *Logger) Debug(msg string, data map[string]any) {
	l.log(DebugLevel, msg, data)
}
