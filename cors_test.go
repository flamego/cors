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
			name:   "method get",
			method: http.MethodGet,
			wantHeaders: map[string]string{
				"Access-Control-Allow-Origin": "*",
				"Access-Control-Max-Age":      "600",
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
}

func TestCustomCORS(t *testing.T) {
	f := flamego.NewWithLogger(&bytes.Buffer{})
	f.Use(CORS(Options{
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
	f.Get("/", func(c flamego.Context) string {
		return "ok"
	})

	tests := []struct {
		name        string
		method      string
		reqHeaders  map[string]string
		wantHeaders map[string]string
		wantCode    int
	}{
		{
			name:   "method get",
			method: http.MethodGet,
			reqHeaders: map[string]string{
				"Origin": "https://example.com",
			},
			wantHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "https://example.com",
				"Access-Control-Max-Age":           "20",
				"Access-Control-Allow-Credentials": "true",
			},
			wantCode: http.StatusOK,
		},
		{
			name:   "default response",
			method: http.MethodOptions,
			reqHeaders: map[string]string{
				"Origin":                         "https://example.com",
				"Access-Control-Request-Headers": "Content-Type",
			},
			wantHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "https://example.com",
				"Access-Control-Max-Age":           "20",
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Allow-Headers":     "Content-Type",
			},
			wantCode: http.StatusOK,
		},
		{
			name:   "error subdomain",
			method: http.MethodOptions,
			reqHeaders: map[string]string{
				"Origin": "https://a.example.com",
			},
			wantHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "",
				"Access-Control-Max-Age":           "",
				"Access-Control-Allow-Credentials": "",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name:   "error scheme",
			method: http.MethodOptions,
			reqHeaders: map[string]string{
				"Origin": "http://example.com",
			},
			wantHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "https://example.com",
				"Access-Control-Max-Age":           "20",
				"Access-Control-Allow-Credentials": "true",
			},
			wantCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest(test.method, "/", nil)
			assert.Nil(t, err)
			for k, v := range test.reqHeaders {
				req.Header.Set(k, v)
			}

			f.ServeHTTP(resp, req)

			assert.Equal(t, test.wantCode, resp.Code)

			for headerKey, headerValue := range test.wantHeaders {
				assert.Equal(t, headerValue, resp.Header().Get(headerKey))
			}
		})
	}
}
