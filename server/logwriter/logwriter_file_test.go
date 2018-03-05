package logwriter

import (
	"context"
	"crypto/rand"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/kolide/fleet/server/config"

	kitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func TestLogger(t *testing.T) {
	tempPath, err := ioutil.TempDir("", "test")
	require.Nil(t, err)
	fileName := path.Join(tempPath, "logwriter")

	logger := kitlog.NewNopLogger()
	conf := config.OsqueryLogFile{
		Path: fileName,
	}
	lw, err := newFile(conf, logger)
	require.Nil(t, err)
	defer os.Remove(fileName)

	randInput := make([]byte, 512)
	rand.Read(randInput)

	ctx := context.Background()

	for i := 0; i < 100; i++ {
		err := lw.Write(ctx, randInput)
		require.Nil(t, err)
	}

	err = lw.Close(ctx)
	assert.Nil(t, err)

	// can't write to a closed logger
	err = lw.Write(ctx, randInput)
	assert.NotNil(t, err)

	// call close twice noop
	err = lw.Close(ctx)
	assert.Nil(t, err)

	info, err := os.Stat(fileName)
	require.Nil(t, err)
	assert.Equal(t, int64(51200), info.Size())
}

func BenchmarkLogger(b *testing.B) {
	tempPath, err := ioutil.TempDir("", "test")
	if err != nil {
		b.Fatal("temp dir failed", err)
	}
	fileName := path.Join(tempPath, "logwriter")

	logger := kitlog.NewNopLogger()
	conf := config.OsqueryLogFile{
		Path: fileName,
	}
	lw, err := newFile(conf, logger)
	if err != nil {
		b.Fatal("new failed ", err)
	}
	defer os.Remove(fileName)

	ctx := context.Background()

	randInput := make([]byte, 512)
	rand.Read(randInput)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := lw.Write(ctx, randInput)
		if err != nil {
			b.Fatal("write failed ", err)
		}
	}

	b.StopTimer()

	lw.Close(ctx)
}

func BenchmarkLumberjack(b *testing.B) {
	tempPath, err := ioutil.TempDir("", "test")
	if err != nil {
		b.Fatal("temp dir failed", err)
	}
	fileName := path.Join(tempPath, "lumberjack")
	lgr := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	}
	defer os.Remove(fileName)

	randInput := make([]byte, 512)
	rand.Read(randInput)
	// first lumberjack write opens file so we count that as part of initialization
	// just to make sure we're comparing apples to apples with our logger
	_, err = lgr.Write(randInput)
	if err != nil {
		b.Fatal("first write failed ", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := lgr.Write(randInput)
		if err != nil {
			b.Fatal("write failed ", err)
		}
	}

	b.StopTimer()

	lgr.Close()
}
