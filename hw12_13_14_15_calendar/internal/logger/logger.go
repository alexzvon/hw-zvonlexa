package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
)

const (
	debugLevel = "debug"
	infoLevel  = "info"
	warnLevel  = "warn"
	errLevel   = "err"
)

type Logger interface {
	Debug(string) error
	Info(string) error
	Warn(string) error
	Error(string) error
	LogHTTPInfo(string) error
	Close() error
}

type sLogger struct {
	mu   sync.Mutex
	file *os.File
}

func (l *sLogger) LogHTTPInfo(mes string) error {
	return l.writer(mes, infoLevel)
}

func (l *sLogger) Debug(mes string) error {
	return l.writer(mes, debugLevel)
}

func (l *sLogger) Info(mes string) error {
	return l.writer(mes, infoLevel)
}

func (l *sLogger) Warn(mes string) error {
	return l.writer(mes, warnLevel)
}

func (l *sLogger) Error(mes string) error {
	return l.writer(mes, errLevel)
}

func (l *sLogger) Close() error {
	if err := l.file.Close(); err != nil {
		return errors.Wrap(err, "cannot close file")
	}

	return nil
}

func (l *sLogger) writer(mes, level string) error {
	t := time.Now().Format(time.RFC822)

	l.mu.Lock()
	_, err := l.file.WriteString(fmt.Sprintf("%s %s: %s\n", t, level, mes))
	l.mu.Unlock()

	if err != nil {
		return errors.Wrap(err, "cannot write")
	}

	return nil
}

func New(path string) (Logger, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create logger")
	}

	sL := &sLogger{
		file: file,
	}

	return sL, nil
}
