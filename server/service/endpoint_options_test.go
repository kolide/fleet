package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testGetOptions(t *testing.T, r *testResource) {
	req, err := http.NewRequest("GET", r.server.URL+"/api/v1/kolide/options", nil)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.userToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var optsResp optionsResponse
	err = json.NewDecoder(resp.Body).Decode(&optsResp)
	require.Nil(t, err)
	require.NotNil(t, optsResp.Options)
	assert.Equal(t, "aws_access_key_id", optsResp.Options[0].Name)
}

func testModifyOptions(t *testing.T, r *testResource) {
	options, err := r.ds.Options()
	require.Nil(t, err)
	require.NotNil(t, options)
	pl := kolide.OptionRequest{options[0:2]}
	val := new(string)
	*val = "foo"
	pl.Options[0].Value = val
	val = new(string)
	*val = "10"
	pl.Options[1].Value = val
	var buff bytes.Buffer
	err = json.NewEncoder(&buff).Encode(pl)
	require.Nil(t, err)
	req, err := http.NewRequest("PATCH", r.server.URL+"/api/v1/kolide/options", &buff)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.userToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	var optsResp optionsResponse
	err = json.NewDecoder(resp.Body).Decode(&optsResp)
	require.Nil(t, err)
	require.NotNil(t, optsResp.Options)
	require.Len(t, optsResp.Options, 2)
	assert.Equal(t, "foo", *optsResp.Options[0].Value)
	assert.Equal(t, "10", *optsResp.Options[1].Value)
	options, err = r.ds.Options()
	assert.Equal(t, "foo", *options[0].Value)
	assert.Equal(t, "10", *options[1].Value)

}

func testModifyOptionsBadJSON(t *testing.T, r *testResource) {
	inJson := `{"options":[
  {"id":1,"name":"aws_access_key_id","type":"string","value":"foo","read_only":false},
  {"id":2,"name":"aws_firehose_period","type":"int","value":"xxs","read_only":false}]}`
	buff := bytes.NewBufferString(inJson)
	req, err := http.NewRequest("PATCH", r.server.URL+"/api/v1/kolide/options", buff)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.userToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode)

}
