package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestValidationFailure(t *testing.T) {
	r := createEmptyTestServer(nil)

	type Foo struct {
		Foo string `json:"foo" binding:"required,len=10"`
		Bar string `json:"bar" binding:"required"`
	}

	r.POST("/foo", func(c *gin.Context) {
		var f Foo
		err := c.BindJSON(&f)

		t.Log(err)
	})

	buff := new(bytes.Buffer)
	buff.Write([]byte(`{"foo": "foo", "bar": "bar"}`))
	req, _ := http.NewRequest("POST", "/foo", buff)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code == http.StatusOK {
		t.Error("Binding should have failed")
	}

	t.Logf("JSON:\n%s", resp.Body.Bytes())

	var res map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &res); err != nil {
		t.Fatalf("Json parse error: %s", err.Error())
	}

	t.Error("fail")
}
