package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime/debug"
	"sync"
	"time"
)

type LogLevel int

func (l LogLevel) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelDebug:
		return "DEBUG"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type LogEntry struct {
	Level      string            `json:"level"`
	Time       string            `json:"time"`
	Message    string            `json:"message"`
	Properties map[string]string `json:"properties,omitempty"`
	Trace      string            `json:"trace,omitempty"`
}

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelOff
)

type Logger struct {
	out      io.Writer
	minLevel LogLevel
	mu       sync.Mutex
}

func New(out io.Writer, minLevel LogLevel) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

func (l *Logger) Info(message string, props map[string]string) {
	l.log(LevelInfo, message, props)
}

func (l *Logger) Warn(message string, props map[string]string) {
	l.log(LevelWarn, message, props)
}

func (l *Logger) Debug(message string, props map[string]string) {
	l.log(LevelDebug, message, props)
}

func (l *Logger) Error(err error, props map[string]string) {
	l.log(LevelError, err.Error(), props)
}

func (l *Logger) Fatal(err error, props map[string]string) {
	l.log(LevelFatal, err.Error(), props)
}

func (l *Logger) log(level LogLevel, message string, props map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	entry := LogEntry{
		Level:      level.String(),
		Time:       time.Now().Format(time.RFC3339),
		Message:    message,
		Properties: props,
	}

	if level >= LevelFatal {
		entry.Trace = string(debug.Stack())
	}

	var line []byte

	line, err := json.Marshal(entry)

	if err != nil {
		line = []byte(fmt.Sprintf("%s: unable to marshal log message: %s", LevelError.String(), err.Error()))
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}
