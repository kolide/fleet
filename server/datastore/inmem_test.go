package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInmem(t *testing.T) {
	for _, f := range testFunctions {
		t.Run(functionName(f), func(t *testing.T) {
			ds, err := New("inmem", "")
			assert.Nil(t, err)
			f(t, ds)
		})
	}
}
