package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/kolide/kolide/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupOptionTest(r *testResource) {
	opts := []kolide.Option{
		kolide.Option{
			ID:   6,
			Name: "aws_access_key_id",
			Value: kolide.OptionValue{
				Val: nil,
			},
			Type:     kolide.OptionTypeString,
			ReadOnly: kolide.NotReadOnly,
		},
		kolide.Option{
			ID:   7,
			Name: "aws_firehose_period",
			Value: kolide.OptionValue{
				Val: nil,
			},
			Type:     kolide.OptionTypeInt,
			ReadOnly: kolide.NotReadOnly,
		},
		kolide.Option{
			ID:   8,
			Name: "host_identifier",
			Value: kolide.OptionValue{
				Val: nil,
			},
			Type:     kolide.OptionTypeString,
			ReadOnly: kolide.NotReadOnly,
		},
		kolide.Option{
			ID:   9,
			Name: "schedule_splay_percent",
			Value: kolide.OptionValue{
				Val: nil,
			},
			Type:     kolide.OptionTypeInt,
			ReadOnly: kolide.NotReadOnly,
		},
	}
	r.SaveOptions(opts)
}

func testGetOptions(t *testing.T, r *testResource) {
	defer r.Close()
	setupOptionTest(r)
	req, err := http.NewRequest("GET", r.server.URL+"/api/v1/kolide/options", nil)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var optsResp optionsResponse
	err = json.NewDecoder(resp.Body).Decode(&optsResp)
	require.Nil(t, err)
	require.NotNil(t, optsResp.Options)
	assert.Len(t, optsResp.Options, 4)
}

func testModifyOptions(t *testing.T, r *testResource) {
	defer r.Close()
	setupOptionTest(r)
	inJson := `{"options":[
  {"id":6,"name":"aws_access_key_id","type":"string","value":"foo","read_only":false},
  {"id":7,"name":"aws_firehose_period","type":"int","value":23,"read_only":false}]}`
	buff := bytes.NewBufferString(inJson)
	req, err := http.NewRequest("PATCH", r.server.URL+"/api/v1/kolide/options", buff)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)

	var optsResp optionsResponse
	err = json.NewDecoder(resp.Body).Decode(&optsResp)
	require.Nil(t, err)
	require.NotNil(t, optsResp.Options)
	require.Len(t, optsResp.Options, 2)
	assert.Equal(t, "foo", optsResp.Options[0].GetValue())
	assert.Equal(t, float64(23), optsResp.Options[1].GetValue())
	assert.True(t, r.SaveOptionsFuncInvoked)
}

func testModifyOptionsValidationFail(t *testing.T, r *testResource) {
	defer r.Close()
	setupOptionTest(r)
	inJson := `{"options":[
  {"id":6,"name":"aws_access_key_id","type":"string","value":"foo","read_only":false},
  {"id":7,"name":"aws_firehose_period","type":"int","value":"xxs","read_only":false}]}`
	buff := bytes.NewBufferString(inJson)
	req, err := http.NewRequest("PATCH", r.server.URL+"/api/v1/kolide/options", buff)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	var errStruct mockValidationError
	err = json.NewDecoder(resp.Body).Decode(&errStruct)
	require.Nil(t, err)
	require.Len(t, errStruct.Errors, 1)
	assert.Equal(t, "aws_firehose_period", errStruct.Errors[0].Name)
	assert.Equal(t, "type mismatch", errStruct.Errors[0].Reason)
}
