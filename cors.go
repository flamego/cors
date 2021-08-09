// Copyright 2021 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package cors is a middleware that generates CORS headers for Flamego.
package cors

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/flamego/flamego"
)

// Options contains options for the cors.CORS middleware.
type Options struct {
	// Scheme may be http or https as accepted schemes or the '*'
	// wildcard to accept any scheme. Default is "http".
	Scheme string
	// AllowDomain is a comma separated list of domains that are
	// allowed to run CORS requests. Special values are the  a single '*'
	// wildcard that will allow any domain to send requests without
	// credentials and the special '!*' wildcard which will reply with
	// requesting domain in the 'access-control-allow-origin' header and
	// hence allow requests from any domain *with* credentials.
	// Default is '*'.
	AllowDomain []string
	// AllowSubdomain allowed subdomains of domains to run CORS requests.
	// Default is false.
	AllowSubdomain bool
	// Methods may be a comma separated list of HTTP-methods to be accepted.
	// Default is ["GET", "POST", "OPTIONS"].
	Methods []string
	// MaxAgeSeconds may be the duration in secs for which the response is cached.
	// Default is 600 * time.Second.
	MaxAge time.Duration
	// AllowCredentials set to false rejects any request with credentials.
	// Default is false.
	AllowCredentials bool
}

func prepareOptions(options []Options) Options {
	var opt Options
	if len(options) > 0 {
		opt = options[0]
	}

	if opt.Scheme == "" {
		opt.Scheme = "http"
	}
	if len(opt.AllowDomain) == 0 {
		opt.AllowDomain = []string{"*"}
	}
	if len(opt.Methods) == 0 {
		opt.Methods = []string{
			http.MethodGet,
			http.MethodOptions,
			http.MethodPost,
		}
	}
	if opt.MaxAge.Seconds() <= 0 {
		opt.MaxAge = time.Duration(600) * time.Second
	}

	return opt
}

// CORS returns a middleware handler that responds to preflight requests with
// adequate "Access-Control-*" response headers.
func CORS(options ...Options) flamego.Handler {
	opt := prepareOptions(options)
	return func(ctx flamego.Context) {
		if ctx.Request().Method != http.MethodOptions {
			return
		}

		headers := map[string]string{
			"Access-Control-Allow-Methods": strings.Join(opt.Methods, ","),
			"Access-Control-Allow-Headers": ctx.Request().Header.Get("Access-Control-Request-Headers"),
			"Access-Control-Max-Age":       fmt.Sprintf("%.0f", opt.MaxAge.Seconds()),
		}
		if opt.AllowDomain[0] == "*" {
			headers["Access-Control-Allow-Origin"] = "*"
		} else {
			origin := ctx.Request().Header.Get("Origin")
			if origin == "" {
				http.Error(ctx.ResponseWriter(), "missing origin header in CORS request", http.StatusBadRequest)
				return
			}

			u, err := url.Parse(origin)
			if err != nil {
				http.Error(ctx.ResponseWriter(), fmt.Sprintf("Failed to parse CORS origin header. Reason: %v", err), http.StatusBadRequest)
				return
			}

			var ok bool
			for _, d := range opt.AllowDomain {
				if u.Hostname() == d || (opt.AllowSubdomain && strings.HasSuffix(u.Hostname(), "."+d)) || d == "!*" {
					ok = true
					break
				}
			}
			if ok {
				if opt.Scheme != "*" {
					u.Scheme = opt.Scheme
				}
				headers["Access-Control-Allow-Origin"] = u.String()
				headers["Access-Control-Allow-Credentials"] = strconv.FormatBool(opt.AllowCredentials)
				headers["Vary"] = "Origin"
			} else {
				http.Error(ctx.ResponseWriter(), fmt.Sprintf("CORS request from prohibited domain %v", origin), http.StatusBadRequest)
				return
			}
		}

		ctx.ResponseWriter().Before(func(w flamego.ResponseWriter) {
			for k, v := range headers {
				w.Header().Set(k, v)
			}
		})

		ctx.ResponseWriter().WriteHeader(http.StatusOK) // return response
		return
	}
}
