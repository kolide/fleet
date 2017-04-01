package logger

import (
	"bufio"
	"errors"
	"io"
	"os"
	"sync"
)

type logger struct {
	file *os.File
	buff *bufio.Writer
	mtx  sync.Mutex
}

func New(path string) (io.WriteCloser, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	buff := bufio.NewWriter(file)
	return &logger{file: file, buff: buff}, nil
}

func (l *logger) Write(b []byte) (int, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if l.buff == nil {
		return 0, errors.New("can't write to a closed file")
	}
	return l.buff.Write(b)
}

func (l *logger) Flush() error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if l.buff == nil {
		return errors.New("can't write to a closed file")
	}
	return l.buff.Flush()
}

func (l *logger) Close() error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if l.buff == nil || l.file == nil {
		return errors.New("file already closed")
	}
	if err := l.buff.Flush(); err != nil {
		return err
	}
	l.buff = nil
	return l.file.Close()
}
