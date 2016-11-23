package datastore

import (
	"testing"

	"github.com/kolide/kolide-ose/server/config"
	"github.com/kolide/kolide-ose/server/datastore/inmem"
	"github.com/stretchr/testify/require"
)

func TestInmem(t *testing.T) {
	config := config.KolideConfig{
		Auth: config.AuthConfig{
			SaltKeySize: 20,
			BcryptCost:  10,
		},
	}

	for _, f := range testFunctions {
		t.Run(functionName(f), func(t *testing.T) {
			ds, err := inmem.New(inmem.WithConfig(&config))
			defer func() { require.Nil(t, ds.Drop()) }()
			require.Nil(t, err)
			f(t, ds)
		})
	}
}
