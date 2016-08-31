// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package kitserver

import (
	"fmt"
	"strings"
)

// matchErr is a test helper that verifies that an error is matched with an expected effor
// source:
// https://github.com/golang/go/blob/ffa2bd27a47ef16e4d6a404dd15781ed5ba21e5d/src/net/http/response_test.go#L865
// wantErr can be nil, an error value to match exactly, or type string to
// match a substring.
func matchErr(err error, wantErr interface{}) error {
	if err == nil {
		if wantErr == nil {
			return nil
		}
		if sub, ok := wantErr.(string); ok {
			return fmt.Errorf("unexpected success; want error with substring %q", sub)
		}
		return fmt.Errorf("unexpected success; want error %v", wantErr)
	}
	if wantErr == nil {
		return fmt.Errorf("%v; want success", err)
	}
	if sub, ok := wantErr.(string); ok {
		if strings.Contains(err.Error(), sub) {
			return nil
		}
		return fmt.Errorf("error = %v; want an error with substring %q", err, sub)
	}
	if err == wantErr {
		return nil
	}
	return fmt.Errorf("%v; want %v", err, wantErr)
}
