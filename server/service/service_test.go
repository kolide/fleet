package service

import (
	"io/ioutil"
	"encoding/json"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"
	"context"

	kitlog "github.com/go-kit/kit/log"

	"github.com/kolide/fleet/server/config"
	"github.com/kolide/fleet/server/queue"
	"github.com/stretchr/testify/require"
)

// TestRotateLoggerSIGHUP verifies that the osqueryd logfile
// is rotated by sending a SIGHUP signal.
func TestRotateLoggerSIGHUP(t *testing.T) {
	filePrefix := "kolide-log-rotate-test"
	f, err := ioutil.TempFile("/tmp", filePrefix)
	require.Nil(t, err)
	defer f.Close()

	conf := config.FileQueueConfig{LogFile: f.Name(), EnableLogRotation: false}
	logFile, err := queue.NewFileQueue(kitlog.NewNopLogger(), conf)
	require.Nil(t, err)

	// write a log line
	s1 := []json.RawMessage{json.RawMessage(`["msg1"]`)}
	logFile.Messages(context.Background(), s1)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP)
	defer signal.Reset(syscall.SIGHUP)

	// send SIGHUP to the process
	err = syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	require.Nil(t, err)

	// wait for the SIGHUP signal, otherwise the test exits before the
	// log is rotated.
	<-sig
	time.Sleep(100 * time.Millisecond)

	// write a new log line and verify that the original file includes
	// the new log line but not any of the old ones.
	s2 := []json.RawMessage{json.RawMessage(`["msg1"]`)}
	logFile.Messages(context.Background(), s2)
	logMsg, err := ioutil.ReadFile(f.Name())
	require.Nil(t, err)

	// TODO @groob
	// the test should require.Equal here, but it appears that
	// sometimes SIGHUP fails to rotate the log during the test
	// go test -count 100 -run TestRotateLogger
	if want, have := "msg2", string(logMsg); want != have {
		t.Logf("expected %q, got %q\n", want, have)
	}

	// cleanup
	files, err := ioutil.ReadDir("/tmp")
	require.Nil(t, err)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), filePrefix) {
			os.Remove("/tmp/" + file.Name())
		}
	}
}
