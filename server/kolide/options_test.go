package kolide

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptionMarshaller(t *testing.T) {
	tests := []struct {
		value         string
		typ           OptionType
		expectSuccess bool
	}{
		{"23", OptionTypeInt, true},
		{"abc", OptionTypeInt, false},
		{"true", OptionTypeFlag, true},
		{"false", OptionTypeFlag, true},
		{"something", OptionTypeFlag, false},
		{"foobar", OptionTypeString, true},
	}

	for _, test := range tests {
		optIn := &Option{1, "foo", test.typ, &test.value, true}
		buff, err := json.Marshal(optIn)
		if !test.expectSuccess {
			assert.NotNil(t, err)
			continue
		}
		require.Nil(t, err)
		optOut := &Option{}
		err = json.Unmarshal(buff, optOut)
		require.Nil(t, err)
		assert.True(t, reflect.DeepEqual(optIn, optOut))
	}

	// test nil
	optIn := &Option{1, "bar", OptionTypeString, nil, true}
	buff, err := json.Marshal(optIn)
	require.Nil(t, err)
	optOut := &Option{}
	err = json.Unmarshal(buff, optOut)
	require.Nil(t, err)
	assert.True(t, reflect.DeepEqual(optIn, optOut))

}

func TestOptionUnmarshaller(t *testing.T) {
	errTypeMismatch := fmt.Errorf("option value type mismatch")

	tests := []struct {
		data string
		err  error
	}{
		{`{"id":1,"name":"foo","type":"string","value":"foobar","read_only":true}`, nil},
		{`{"id":1,"name":"foo","type":"float","value":"foobar","read_only":true}`, fmt.Errorf("option type 'float' invalid")},
		{`{"id":1,"name":"foo","type":"int","value":"foobar","read_only":true}`, errTypeMismatch},
		{`{"id":1,"name":"foo","type":"flag","value":"foobar","read_only":true}`, errTypeMismatch},
	}

	for _, test := range tests {
		buff := []byte(test.data)
		opt := &Option{}
		err := json.Unmarshal(buff, opt)
		if test.err != nil {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.Nil(t, err)
		}
	}

}
