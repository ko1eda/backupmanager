package log

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// Logger writes logs to file
type Logger struct {
	mu     sync.Mutex
	logger *log.Logger
}

// New returns a logger instance that writes to stderr and to a file
// It also returns a function to close tthat file.
func New() (logger *Logger, closer func()) {
	// create a file if it does not exist, if it does exist open it and append to it
	f, err := os.OpenFile(makePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("LogFileOpenError: %v", err)
	}

	mw := io.MultiWriter(os.Stderr, f)

	l := log.New(mw, "", log.LstdFlags)

	return &Logger{logger: l}, func() { f.Close() }
}

// Log writes a log to stdout, stderr, and to the file specified in the path
func (l *Logger) Log(msg ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Println(msg...)
}

// Creates a path to log files in the current working directory
func makePath() string {
	wd, err := os.Getwd()

	if err != nil {
		log.Fatalf("LogGetWorkingDirError: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(wd, "storage", "logs"), 0770); err != nil {
		log.Fatalf("LogCreateStorageDirError: %v", err)
	}

	return filepath.Join(wd, "storage", "logs", "application.log")
}
