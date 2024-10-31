package logger

import (
	"io"
	"log"
	"os"
	"sync"
	"time"
)

// Logger is a custom logger with support for verbosity and timing.
type Logger struct {
	mu       sync.Mutex
	verbose  bool
	logger   *log.Logger
	startTimes map[string]time.Time
}

// singleton instance
var (
	instance *Logger
	once     sync.Once
)

// GetLogger returns the singleton Logger instance.
func GetLogger(verbose bool) *Logger {
	once.Do(func() {
		instance = &Logger{
			verbose:  verbose,
			logger:   log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
			startTimes: make(map[string]time.Time),
		}
		if !verbose {
			instance.logger.SetOutput(io.Discard)
		}
	})
	return instance
}

// Info logs informational messages.
func (l *Logger) Info(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.verbose {
		l.logger.Printf("INFO: "+format, v...)
	}
}

// Error logs error messages.
func (l *Logger) Error(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.verbose {
		l.logger.Printf("ERROR: "+format, v...)
	}
}

// StartTimer records the start time of an operation.
func (l *Logger) StartTimer(operation string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.verbose {
		l.startTimes[operation] = time.Now()
		l.logger.Printf("Started operation: %s", operation)
	}
}

// StopTimer logs the duration of an operation.
func (l *Logger) StopTimer(operation string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.verbose {
		start, exists := l.startTimes[operation]
		if exists {
			duration := time.Since(start)
			l.logger.Printf("Completed operation: %s in %v", operation, duration)
			delete(l.startTimes, operation)
		} else {
			l.logger.Printf("StopTimer called for unknown operation: %s", operation)
		}
	}
}
