package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kolide/fleet/server/okforward"
)

// TODO remove, only used for debugging
func main() {

	// Logging.
	var logger log.Logger
	{
		logLevel := level.AllowInfo()
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = level.NewFilter(logger, logLevel)
	}

	w, err := okforward.New(logger, []string{"localhost"})
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	go func() {
		for {
			w.Write([]byte("hello\n"))
			time.Sleep(1 * time.Second)
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig

}
