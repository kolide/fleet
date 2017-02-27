package service

import (
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/require"
)

// TestRotateLoggerSIGHUP verifies that the osqueryd logfile
// is rotated by sending a SIGHUP signal.
func TestRotateLoggerSIGHUP(t *testing.T) {
	filePrefix := "kolide-log-rotate-test"
	f, err := ioutil.TempFile("/tmp", filePrefix)
	require.Nil(t, err)
	defer f.Close()

	logFile := osqueryLogFile(f.Name(), log.NewNopLogger())

	// write a log line
	logFile.Write([]byte("msg1"))

	// send SIGHUP to the process
	err = syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	require.Nil(t, err)

	// wait for the SIGHUP signal, otherwise the test exits before the
	// log is rotated.
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGHUP)
	<-sig

	// write a new log line and verify that the original file includes
	// the new log line but not any of the old ones.
	logFile.Write([]byte("msg2"))
	logMsg, err := ioutil.ReadFile(f.Name())
	require.Nil(t, err)
	require.Equal(t, string(logMsg), "msg2")

	// cleanup
	files, err := ioutil.ReadDir("/tmp")
	require.Nil(t, err)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), filePrefix) {
			os.Remove("/tmp/" + file.Name())
		}
	}
}
