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

func setupDecoratorTest(r *testResource) {
	decs := []kolide.Decorator{
		kolide.Decorator{
			Type:  kolide.DecoratorLoad,
			Query: "select something from foo;",
		},
		kolide.Decorator{
			Type:  kolide.DecoratorLoad,
			Query: "select bar from foo;",
		},
		kolide.Decorator{
			Type:  kolide.DecoratorAlways,
			Query: "select x from y;",
		},
		kolide.Decorator{
			Type:     kolide.DecoratorInterval,
			Query:    "select name from system_info;",
			Interval: 3600,
		},
	}
	for _, d := range decs {
		r.NewDecorator(&d)
	}
}

func testListDecorator(t *testing.T, r *testResource) {
	defer r.Close()
	setupDecoratorTest(r)
	req, err := http.NewRequest("GET", r.server.URL+"/api/v1/kolide/decorators", nil)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var decs listDecoratorResponse
	err = json.NewDecoder(resp.Body).Decode(&decs)
	require.Nil(t, err)

	assert.Len(t, decs.Decorators, 4)
}

func testNewDecorator(t *testing.T, r *testResource) {
	defer r.Close()
	newDec := newDecoratorRequest{
		Payload: kolide.DecoratorPayload{
			DecoratorType: "load",
			Query:         "select x from y;",
		},
	}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(newDec)
	require.Nil(t, err)
	req, err := http.NewRequest("POST", r.server.URL+"/api/v1/kolide/decorators", &buffer)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var dec decoratorResponse
	err = json.NewDecoder(resp.Body).Decode(&dec)
	require.Nil(t, err)
	require.NotNil(t, dec.Decorator)
	assert.Equal(t, kolide.DecoratorLoad, dec.Decorator.Type)
	assert.Equal(t, "select x from y;", dec.Decorator.Query)
}

// invalid json
func testNewDecoratorFailType(t *testing.T, r *testResource) {
	defer r.Close()
	newDec := newDecoratorRequest{
		Payload: kolide.DecoratorPayload{
			DecoratorType: "zip",
			Query:         "select x from y;",
		},
	}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(newDec)
	require.Nil(t, err)
	req, err := http.NewRequest("POST", r.server.URL+"/api/v1/kolide/decorators", &buffer)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	var errStruct mockValidationError
	err = json.NewDecoder(resp.Body).Decode(&errStruct)
	require.Nil(t, err)
	require.Len(t, errStruct.Errors, 1)
	assert.Equal(t, "undefined decorator type", errStruct.Errors[0].Reason)

}

func testNewDecoratorFailValidation(t *testing.T, r *testResource) {
	defer r.Close()
	newDec := newDecoratorRequest{
		Payload: kolide.DecoratorPayload{
			DecoratorType: "interval",
			Query:         "select x from y;",
			Interval:      uintPtr(3601),
		},
	}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(newDec)
	require.Nil(t, err)
	req, err := http.NewRequest("POST", r.server.URL+"/api/v1/kolide/decorators", &buffer)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	var errStruct mockValidationError
	err = json.NewDecoder(resp.Body).Decode(&errStruct)
	require.Nil(t, err)
	require.Len(t, errStruct.Errors, 1)
	assert.Equal(t, "interval value must be divisible by 60", errStruct.Errors[0].Reason)
}

func testDeleteDecorator(t *testing.T, r *testResource) {
	defer r.Close()
	setupDecoratorTest(r)
	req, err := http.NewRequest("DELETE", r.server.URL+"/api/v1/kolide/decorators/1", nil)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	decs, _ := r.ListDecorators()
	assert.Len(t, decs, 3)
}
