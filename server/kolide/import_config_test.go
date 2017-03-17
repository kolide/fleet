package kolide

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntervalUnmarshal(t *testing.T) {
	scenarios := []struct {
		testVal        interface{}
		errExpected    bool
		expectedResult uint
	}{
		{"100", false, 100},
		{float64(123), false, 123},
		{nil, false, 0},
		{"hi there", true, 0},
	}
	for _, scenario := range scenarios {
		v, e := unmarshalInterval(scenario.testVal)
		if scenario.errExpected {
			assert.NotNil(t, e)
			continue
		}
		assert.Nil(t, e)
		assert.Equal(t, scenario.expectedResult, v)
	}
}

func TestPackNameMapUnmarshal(t *testing.T) {
	pnm := PackNameMap{
		"path": "/this/is/a/path",
		"details": PackDetails{
			Queries: QueryNameToQueryDetailsMap{
				"q1": QueryDetails{
					Query:    "select from foo",
					Interval: 100,
					Removed:  new(bool),
					Platform: strptr("linux"),
					Shard:    new(uint),
					Snapshot: new(bool),
				},
			},
			Discovery: []string{
				"select from something",
			},
		},
	}
	b, _ := json.Marshal(pnm)
	actual := make(PackNameMap)
	err := json.Unmarshal(b, &actual)
	require.Nil(t, err)
	assert.Len(t, actual, 2)

	pnm = PackNameMap{
		"path": "/this/is/a/path",
		"details": PackDetails{
			Queries: QueryNameToQueryDetailsMap{
				"q1": QueryDetails{
					Query:    "select from foo",
					Interval: 100,
					Removed:  new(bool),
					Platform: strptr("linux"),
					Shard:    new(uint),
					Snapshot: new(bool),
				},
			},
			Shard:    uintptr(float64(10)),
			Version:  strptr("1.0"),
			Platform: "linux",
			Discovery: []string{
				"select from something",
			},
		},
		"details2": PackDetails{
			Queries: QueryNameToQueryDetailsMap{
				"q1": QueryDetails{
					Query:    "select from bar",
					Interval: 100,
					Removed:  new(bool),
					Platform: strptr("linux"),
					Shard:    new(uint),
					Snapshot: new(bool),
				},
			},
			Shard:    uintptr(float64(10)),
			Version:  strptr("1.0"),
			Platform: "linux",
		},
	}

	b, _ = json.Marshal(pnm)
	actual = make(PackNameMap)
	err = json.Unmarshal(b, &actual)
	require.Nil(t, err)
	assert.Len(t, actual, 3)
}
