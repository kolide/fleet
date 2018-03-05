package logwriter

import (
	"bufio"
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kolide/fleet/server/config"

	kitlog "github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"gopkg.in/natefinch/lumberjack.v2"
)

type fileWriter struct {
	// No rotation.
	file *os.File
	buff *bufio.Writer

	// Rotation.
	jack *lumberjack.Logger

	logger kitlog.Logger
	mtx    sync.Mutex
}

func newFile(conf config.OsqueryLogFile, logger kitlog.Logger) (*fileWriter, error) {
	// TODO(thorduri): Deal with compat if required.
	if conf.Path == "" {
		return nil, nil
	}

	fw := new(fileWriter)

	if conf.EnableLogRotation {
		logger = kitlog.With(logger, "type", "file", "rotation", true)

		fw.jack = &lumberjack.Logger{
			Filename:   conf.Path,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
		}

		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGHUP)

		go func() {
			for {
				<-sig // Block on signal.
				if err := fw.jack.Rotate(); err != nil {
					logger.Log("err", err)
				}
			}
		}()
	} else {
		logger = kitlog.With(logger, "type", "file", "rotation", false)

		file, err := os.OpenFile(conf.Path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
		buff := bufio.NewWriter(file)

		fw.file = file
		fw.buff = buff
	}

	fw.logger = logger
	return fw, nil
}

// Write bytes to file
func (fw *fileWriter) Write(ctx context.Context, b []byte) error {
	fw.mtx.Lock()
	defer fw.mtx.Unlock()

	if fw.jack != nil {
		if _, err := fw.jack.Write(b); err != nil {
			return errors.Wrapf(err, "filewriter: write failed")
		}
		return nil
	}

	if fw.buff == nil || fw.file == nil {
		return errors.New("filewriter: can't write to closed file")
	}
	if _, statErr := os.Stat(fw.file.Name()); os.IsNotExist(statErr) {
		f, err := os.OpenFile(fw.file.Name(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return errors.Wrapf(err, "filewriter: can't create file %s", fw.file.Name())
		}
		fw.file = f
		fw.buff = bufio.NewWriter(f)
	}

	if _, err := fw.buff.Write(b); err != nil {
		return errors.Wrapf(err, "filewriter: write failed")
	}
	return nil
}

// Flush write all buffered bytes to log file
func (fw *fileWriter) Flush(ctx context.Context) error {
	fw.mtx.Lock()
	defer fw.mtx.Unlock()

	if fw.jack != nil {
		return nil
	}

	if fw.buff == nil {
		return errors.New("filewriter: can't write to a closed file")
	}

	if err := fw.buff.Flush(); err != nil {
		return errors.Wrapf(err, "filewriter: flush failed")
	}
	return nil
}

// Close log file
func (fw *fileWriter) Close(ctx context.Context) error {
	fw.mtx.Lock()
	defer fw.mtx.Unlock()

	if fw.jack != nil {
		if err := fw.jack.Close(); err != nil {
			return errors.Wrapf(err, "filewriter: close failed")
		}
		return nil
	}

	if fw.buff != nil {
		if err := fw.buff.Flush(); err != nil {
			return errors.Wrapf(err, "filewriter: flush in close failed")
		}
		fw.buff = nil
	}
	if fw.file != nil {
		if err := fw.file.Close(); err != nil {
			return errors.Wrapf(err, "filewriter: close failed")
		}
		fw.file = nil
	}

	return nil
}
