package queue 

import (
	"os"
	"fmt"
	"os/signal"
	"syscall"
	"context"
	"encoding/json"

	kitlog "github.com/go-kit/kit/log"
	"github.com/kolide/fleet/server/config"
	"gopkg.in/natefinch/lumberjack.v2"
)


type FileQueue struct {
	logger         kitlog.Logger
	enableRotation bool
	path           string
	l              *lumberjack.Logger
}

func NewFileQueue(appLogger kitlog.Logger, conf config.FileQueueConfig) (*FileQueue, error){
	var osquerydLogger *lumberjack.Logger

	if conf.EnableLogRotation {
		osquerydLogger = &lumberjack.Logger{
			Filename:   conf.LogFile,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
		}
		appLogger = kitlog.With(appLogger, "component", "osqueryd-logger")
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGHUP)
		go func() {
			for {
				<-sig //block on signal
				if err := osquerydLogger.Rotate(); err != nil {
					appLogger.Log("err", err)
				}
			}
		}()
	} else {
		osquerydLogger = &lumberjack.Logger{Filename: conf.LogFile}
	}
	q := &FileQueue{
		logger: appLogger,
		enableRotation: conf.EnableLogRotation,
		path: conf.LogFile,
		l: osquerydLogger,
	}
	return q, nil
}



func (fq *FileQueue) Messages(ctx context.Context, logs []json.RawMessage) error {
	for _, log := range logs {
		if _, err := fq.l.Write(append(log, '\n')); err != nil {
			return fmt.Errorf("error writing status log: " + err.Error())
		}
	}
	/*
	err := fq.l.Flush()
	if err != nil {
		error
		return fmt.Errorf("error flushing status log: " + err.Error())
	}
    */
	return nil
}
