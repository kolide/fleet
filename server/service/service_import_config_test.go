package service

import (
	"testing"

	"github.com/kolide/kolide-ose/server/config"
	"github.com/kolide/kolide-ose/server/datastore/inmem"
	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createServiceMockForImport(t *testing.T) *service {
	ds, err := inmem.New(config.TestConfig())
	require.Nil(t, err)
	return &service{
		ds: ds,
	}
}

func TestHashQuery(t *testing.T) {
	q1 := `SELECT * FROM t1 INNER JOIN ON
		t1.id = t2.t1id
		WHERE t1.name = 'foo'  `
	q2 := "SELECT * from t1 INNER JOIN\tON t1.id = t2.t1id	WHERE t1.name = 'foo';"
	h1 := hashQuery("platform", q1)
	h2 := hashQuery("platform", q2)
	assert.Equal(t, h1, h2)
	q2 = "SELECT * from t1 INNER JOIN\tON t1.id = t2.t1id	WHERE t2.name = 'foo';"
	h2 = hashQuery("platform", q2)
	assert.NotEqual(t, h1, h2)

}

func TestOptionsImportConfig(t *testing.T) {
	opts := kolide.OptionNameToValueMap{
		"aws_access_key_id": "foo",
	}
	resp := kolide.NewImportConfigResponse()
	svc := createServiceMockForImport(t)
	err := svc.importOptions(opts, resp)
	require.Nil(t, err)
	status := resp.Status(kolide.OptionsSection)
	require.NotNil(t, status)
	assert.Equal(t, 1, status.ImportCount)
	opt, err := svc.ds.OptionByName("aws_access_key_id")
	require.Nil(t, err)
	assert.Equal(t, "foo", opt.GetValue())
	require.Len(t, status.Messages, 1)
	assert.Equal(t, "set aws_access_key_id value to foo", status.Messages[0])
}

func TestOptionsImportConfigWithSkips(t *testing.T) {
	opts := kolide.OptionNameToValueMap{
		"aws_access_key_id":     "foo",
		"aws_secret_access_key": "secret",
		// this should be skipped because it's already set
		"aws_firehose_period": 100,
		// these should be skipped because it's read only
		"disable_distributed": false,
		"pack_delimiter":      "x",
		// this should be skipped because it's not an option we know about
		"wombat": "not venomous",
	}
	resp := kolide.NewImportConfigResponse()
	svc := createServiceMockForImport(t)
	// set option val, it should be skipped
	opt, err := svc.ds.OptionByName("aws_firehose_period")
	require.Nil(t, err)
	opt.SetValue(23)
	err = svc.ds.SaveOptions([]kolide.Option{*opt})
	require.Nil(t, err)
	err = svc.importOptions(opts, resp)
	require.Nil(t, err)
	status := resp.Status(kolide.OptionsSection)
	require.NotNil(t, status)
	assert.Equal(t, 2, status.ImportCount)
	assert.Equal(t, 4, status.SkipCount)
	assert.Len(t, status.Warnings[kolide.OptionAlreadySet], 1)
	assert.Len(t, status.Warnings[kolide.OptionReadonly], 2)
	assert.Len(t, status.Warnings[kolide.OptionUnknown], 1)
	assert.Len(t, status.Messages, 2)
}

var boolptr = func(v bool) *bool {
	b := new(bool)
	*b = v
	return b
}
var stringptr = func(v string) *string {
	s := new(string)
	*s = v
	return s
}
var uintptr = func(n uint) *uint {
	i := new(uint)
	*i = n
	return i
}

func TestPacksImportConfig(t *testing.T) {
	svc := createServiceMockForImport(t)
	q1 := kolide.QueryDetails{
		Query:    "select * from foo",
		Interval: 100,
		Removed:  boolptr(false),
		Platform: stringptr("linux"),
		Version:  stringptr("1.0"),
	}
	q2 := kolide.QueryDetails{
		Query:    "select * from bar",
		Interval: 50,
		Removed:  boolptr(false),
		Platform: stringptr("linux"),
		Version:  stringptr("1.0"),
	}
	q3 := kolide.QueryDetails{
		Query:    "select * from baz",
		Interval: 500,
		Removed:  boolptr(false),
		Platform: stringptr("linux"),
		Version:  stringptr("1.0"),
	}

	importConfig := kolide.ImportConfig{
		Packs: kolide.PackNameMap{
			"ext1": "/home/usr/ext1.json",
			"pack1": kolide.PackDetails{
				Queries: kolide.QueryNameToQueryDetailsMap{
					"q1": q1,
					"q2": q2,
				},
				Discovery: []string{
					"select * from zz",
					"select id, xx from yy",
				},
			},
			"*": "/home/usr/packs/*",
		},
		ExternalPacks: kolide.PackNameToPackDetails{
			"ext1": kolide.PackDetails{
				Queries: kolide.QueryNameToQueryDetailsMap{
					"q1": q1,
				},
				Discovery: []string{
					"select * from zz",
					"select a, b, c from processes",
				},
			},
			"ext2": kolide.PackDetails{
				Queries: kolide.QueryNameToQueryDetailsMap{
					"q3": q3,
				},
			},
		},
		GlobPackNames: []string{"ext2"},
	}
	resp := kolide.NewImportConfigResponse()
	user := &kolide.User{
		Username: "bob",
		Password: []byte("secret"),
		Email:    "bob@something.com",
		Admin:    false,
		AdminForcedPasswordReset: false,
	}
	user, err := svc.ds.NewUser(user)
	require.Nil(t, err)

	packs, err := importConfig.CollectPacks()
	require.Nil(t, err)
	assert.Len(t, packs, 3)
	err = svc.importPacks(user.ID, &importConfig, resp)
	require.Nil(t, err)
	queries, err := svc.ds.ListQueries(kolide.ListOptions{})
	require.Nil(t, err)
	assert.Len(t, queries, 3)
	pack, ok, err := svc.ds.PackByName("pack1")
	require.Nil(t, err)
	require.True(t, ok)
	sqs, err := svc.ds.ListScheduledQueriesInPack(pack.ID, kolide.ListOptions{})
	require.Nil(t, err)
	assert.Len(t, sqs, 2)
	labels, err := svc.ds.ListLabels(kolide.ListOptions{})
	require.Nil(t, err)
	assert.Len(t, labels, 3)
}
