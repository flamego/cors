// Copyright 2021 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cors

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/flamego/flamego"
)

func TestCORS(t *testing.T) {
	f := flamego.NewWithLogger(&bytes.Buffer{})
	f.Use(CORS())
	f.Get("/", func(c flamego.Context) string {
		return "ok"
	})

	tests := []struct {
		name    string
		method  string
		headers map[string]string
	}{
		{
			name:   "Error method",
			method: http.MethodGet,
			headers: map[string]string{
				"Access-Control-Allow-Origin": "",
				"Access-Control-Max-Age":      "",
			},
		},
		{
			name:   "Default cors response",
			method: http.MethodOptions,
			headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
				"Access-Control-Max-Age":      "600",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest(test.method, "/", nil)
			assert.Nil(t, err)

			f.ServeHTTP(resp, req)

			for headerKey, headerValue := range test.headers {
				assert.Equal(t, headerValue, resp.Header().Get(headerKey))
			}
		})
	}
}
