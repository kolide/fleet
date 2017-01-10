package datastore

import (
	"testing"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testFileIntegrityMonitoring(t *testing.T, ds kolide.Datastore) {
	fp := &kolide.FilePath{
		SectionName: "fp1",
		Paths: []string{
			"path1",
			"path2",
			"path3",
		},
	}
	fp, err := ds.NewFilePath(fp)
	require.Nil(t, err)
	assert.True(t, fp.ID > 0)
	fp = &kolide.FilePath{
		SectionName: "fp2",
		Paths: []string{
			"path4",
			"path5",
		},
	}
	_, err = ds.NewFilePath(fp)
	require.Nil(t, err)

	actual, err := ds.FilePaths()
	require.Nil(t, err)
	assert.Len(t, actual, 2)
	assert.Len(t, actual["fp1"], 3)
	assert.Len(t, actual["fp2"], 2)
}
