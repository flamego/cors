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

	tests := []struct {
		name        string
		method      string
		wantHeaders map[string]string
	}{
		{
			name:   "error method",
			method: http.MethodGet,
			wantHeaders: map[string]string{
				"Access-Control-Allow-Origin": "",
				"Access-Control-Max-Age":      "",
			},
		},
		{
			name:   "default response",
			method: http.MethodOptions,
			wantHeaders: map[string]string{
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

			for headerKey, headerValue := range test.wantHeaders {
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
		AllowCredentials: true,
	}))
	f2.Get("/", func(c flamego.Context) string {
		return "ok"
	})

	customTests := []struct {
		name        string
		method      string
		reqHeaders  map[string]string
		respHeaders map[string]string
		statueCode  int
	}{
		{
			name:   "Error method",
			method: http.MethodGet,
			reqHeaders: map[string]string{
				"Origin": "https://example.com",
			},
			respHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "",
				"Access-Control-Max-Age":           "",
				"Access-Control-Allow-Credentials": "",
			},
			statueCode: 200,
		},
		{
			name:   "Default cors response",
			method: http.MethodOptions,
			reqHeaders: map[string]string{
				"Origin":                         "https://example.com",
				"Access-Control-Request-Headers": "Content-Type",
			},
			respHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "https://example.com",
				"Access-Control-Max-Age":           "20",
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Allow-Headers":     "Content-Type",
			},
			statueCode: 200,
		},
		{
			name:   "Error subdomain",
			method: http.MethodOptions,
			reqHeaders: map[string]string{
				"Origin": "https://a.example.com",
			},
			respHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "",
				"Access-Control-Max-Age":           "",
				"Access-Control-Allow-Credentials": "",
			},
			statueCode: 400,
		},
		{
			name:   "Error scheme",
			method: http.MethodOptions,
			reqHeaders: map[string]string{
				"Origin": "http://example.com",
			},
			respHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "https://example.com",
				"Access-Control-Max-Age":           "20",
				"Access-Control-Allow-Credentials": "true",
			},
			statueCode: 200,
		},
	}

	for _, test := range customTests {
		t.Run(test.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest(test.method, "/", nil)
			assert.Nil(t, err)
			for k, v := range test.reqHeaders {
				req.Header.Set(k, v)
			}

			f2.ServeHTTP(resp, req)

			assert.Equal(t, test.statueCode, resp.Code)

			for headerKey, headerValue := range test.respHeaders {
				assert.Equal(t, headerValue, resp.Header().Get(headerKey))
			}
		})
	}

}
