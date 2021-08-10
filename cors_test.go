// Copyright 2021 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cors

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/flamego/flamego"
)

func TestCORS(t *testing.T) {
	f := flamego.NewWithLogger(&bytes.Buffer{})
	f.Use(CORS())
	f.Get("/", func(c flamego.Context) string {
		return "ok"
	})

	defaultTests := []struct {
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

	for _, test := range defaultTests {
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

	f2 := flamego.NewWithLogger(&bytes.Buffer{})
	f2.Use(CORS(Options{
		Scheme: "https",
		AllowDomain: []string{
			"example.com",
		},
		AllowSubdomain: false,
		Methods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodOptions,
		},
		MaxAge:           time.Duration(20) * time.Second,
		AllowCredentials: false,
	}))
	f2.Get("/", func(c flamego.Context) string {
		return "ok"
	})

	customTests := []struct {
		name       string
		origin     string
		method     string
		headers    map[string]string
		statueCode int
	}{
		{
			name:   "Error method",
			origin: "https://example.com",
			method: http.MethodGet,
			headers: map[string]string{
				"Access-Control-Allow-Origin": "",
				"Access-Control-Max-Age":      "",
			},
			statueCode: 200,
		},
		{
			name:   "Default cors response",
			origin: "https://example.com",
			method: http.MethodOptions,
			headers: map[string]string{
				"Access-Control-Allow-Origin": "https://example.com",
				"Access-Control-Max-Age":      "20",
			},
			statueCode: 200,
		},
		{
			name:   "Error subdomain",
			origin: "https://a.example.com",
			method: http.MethodOptions,
			headers: map[string]string{
				"Access-Control-Allow-Origin": "",
				"Access-Control-Max-Age":      "",
			},
			statueCode: 400,
		},
		{
			name:   "Error scheme",
			origin: "http://example.com",
			method: http.MethodOptions,
			headers: map[string]string{
				"Access-Control-Allow-Origin": "https://example.com",
				"Access-Control-Max-Age":      "20",
			},
			statueCode: 200,
		},
	}

	for _, test := range customTests {
		t.Run(test.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest(test.method, "/", nil)
			assert.Nil(t, err)
			req.Header.Set("Origin", test.origin)

			f2.ServeHTTP(resp, req)

			assert.Equal(t, test.statueCode, resp.Code)

			for headerKey, headerValue := range test.headers {
				assert.Equal(t, headerValue, resp.Header().Get(headerKey))
			}
		})
	}

}
