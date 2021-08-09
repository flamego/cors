// Copyright 2021 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package cors is a middleware that generates CORS headers for Flamego.
package cors

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flamego/flamego"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	f := flamego.NewWithLogger(&bytes.Buffer{})
	f.Use(CORS())
	f.Get("/", func(c flamego.Context) string {
		return "ok"
	})

	resp := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodOptions, "/", nil)
	assert.Nil(t, err)

	f.ServeHTTP(resp, req)
	assert.Equal(t, "*", resp.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "600", resp.Header().Get("Access-Control-Max-Age"))
}
