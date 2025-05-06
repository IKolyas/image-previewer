package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// LogLevel представляет уровень логирования
type LogLevel int

const (
	// LevelError означает уровень логирования ошибок
	LevelError LogLevel = iota
	// LevelWarn означает уровень логирования предупреждений
	LevelWarn
	// LevelInfo означает уровень логирования информации
	LevelInfo
	// LevelDebug означает уровень логирования отладки
	LevelDebug
)

// Logger представляет логгер
type Logger struct {
	mu     sync.Mutex
	level  LogLevel
	output io.Writer
}

// New создает новый логгер
func New(ctx context.Context, level string, outputFile string) (*Logger, error) {
	lvl, err := parseLogLevel(level)
	if err != nil {
		return nil, err
	}

	var output io.Writer = os.Stdout

	if outputFile != "" && outputFile != "stdout" {
		file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("failed to open log file: %v", err)
		}
		go func() {
			<-ctx.Done()
			file.Close()
		}()
		output = file
	}

	return &Logger{
		level:  lvl,
		output: output,
	}, nil
}

// parseLogLevel парсит уровень логирования из строки
func parseLogLevel(level string) (LogLevel, error) {
	switch strings.ToLower(level) {
	case "error":
		return LevelError, nil
	case "warn", "warning":
		return LevelWarn, nil
	case "info":
		return LevelInfo, nil
	case "debug":
		return LevelDebug, nil
	default:
		return LevelInfo, fmt.Errorf("unknown log level: %s", level)
	}
}

// SetLevel устанавливает уровень логирования
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// GetLevel возвращает текущий уровень логирования
func (l *Logger) GetLevel() LogLevel {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

// Log пишет сообщение в лог
func (l *Logger) Log(level LogLevel, msg string) {
	if level > l.level {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	_, err := fmt.Fprintf(l.output, "[%s] [%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), levelToString(level), msg)
	if err != nil {
		log.Fatalf("failed to write log: %v", err)
	}
}

// levelToString преобразует уровень логирования в строку
func levelToString(level LogLevel) string {
	switch level {
	case LevelError:
		return "ERROR"
	case LevelWarn:
		return "WARN"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

// Error пишет сообщение об ошибке в лог
func (l *Logger) Error(msg string) {
	l.Log(LevelError, msg)
}

// Warn пишет сообщение предупреждения в лог
func (l *Logger) Warn(msg string) {
	l.Log(LevelWarn, msg)
}

// Info пишет сообщение информации в лог
func (l *Logger) Info(msg string) {
	l.Log(LevelInfo, msg)
}

// Debug пишет сообщение отладки в лог
func (l *Logger) Debug(msg string) {
	l.Log(LevelDebug, msg)
}
