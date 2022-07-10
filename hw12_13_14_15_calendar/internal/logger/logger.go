package logger

import (
	"log"
	"os"
	"time"

	"github.com/alexzvon/hw12_13_14_15_calendar/internal/myutils"
	"github.com/pkg/errors"
)

type Logger interface {
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	LogHTTPInfo(string)
	Close()
}

type sLogger struct {
	done chan struct{}
	mes  chan string
	path string
	file *os.File
}

func (l *sLogger) LogHTTPInfo(s string) {
	l.mes <- s
}

func (l *sLogger) Debug(mes string) {
	t := time.Now().Format(time.RFC822)
	l.mes <- myutils.ConCat(t, " DEBUG: ", mes, "\n")
}

func (l *sLogger) Info(mes string) {
	t := time.Now().Format(time.RFC822)
	l.mes <- myutils.ConCat(t, " INFO: ", mes, "\n")
}

func (l *sLogger) Warn(mes string) {
	t := time.Now().Format(time.RFC822)
	l.mes <- myutils.ConCat(t, " WARN: ", mes, "\n")
}

func (l *sLogger) Error(mes string) {
	t := time.Now().Format(time.RFC822)
	l.mes <- myutils.ConCat(t, " ERROR: ", mes, "\n")
}

func (l *sLogger) Close() {
	close(l.done)

	if err := l.file.Close(); err != nil {
		log.Fatalln(err)
	}
}

func New(path string) (Logger, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create logger")
	}

	sL := &sLogger{
		done: make(chan struct{}),
		mes:  make(chan string),
		path: path,
		file: file,
	}

	go writer(sL.done, sL.mes, sL.file)

	return sL, nil
}

func writer(done <-chan struct{}, mes chan string, file *os.File) {
	defer close(mes)

	for {
		select {
		case m, ok := <-mes:
			if !ok {
				return
			}

			select {
			case <-done:
				return
			default:
				_, err := file.WriteString(m)
				if err != nil {
					bL := []byte(myutils.ConCat("Cannot write logger file, message - ", m, "\n"))
					n, err := os.Stderr.Write(bL)
					if err != nil {
						log.Println(err)
						return
					}
					if n != len(bL) {
						log.Println("Cannot write to Stderr")
						return
					}
				}
			}
		case <-done:
			return
		}
	}
}
