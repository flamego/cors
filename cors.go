package cors

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/flamego/flamego"
)

const anyDomain = "!*"

// Options to configure the CORS middleware read from the [cors] section of the ini configuration file.
// ref: https://stackoverflow.com/questions/54300997/is-it-possible-to-cache-http-options-response?noredirect=1#comment95790277_54300997
// ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age
type Options struct {
	Section string
	// SCHEME may be http or https as accepted schemes or the '*' wildcard to accept any scheme.
	Scheme string
	// ALLOW_DOMAIN may be a comma separated list of domains that are allowed to run CORS requests
	// Special values are the  a single '*' wildcard that will allow any domain to send requests without
	AllowDomain []string
	// AllowSubdomain allowed
	AllowSubdomain bool
	// METHODS may be a comma separated list of HTTP-methods to be accepted.
	Methods []string
	// MAX_AGE_SECONDS may be the duration in secs for which the response is cached (default 600).
	MaxAgeSeconds int
	// ALLOW_CREDENTIALS set to false rejects any request with credentials.
	AllowCredentials bool
}

func prepareOptions(options []Options) Options {
	var opt Options
	if len(options) > 0 {
		opt = options[0]
	}

	if len(opt.Section) == 0 {
		opt.Section = "cors"
	}
	if len(opt.Scheme) == 0 {
		opt.Scheme = "http"
	}
	if len(opt.AllowDomain) == 0 {
		opt.AllowDomain = []string{"*"}
	}
	if len(opt.Methods) == 0 {
		opt.Methods = []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		}
	}
	if opt.MaxAgeSeconds <= 0 {
		opt.MaxAgeSeconds = 600
	}

	return opt
}

// CORS responds to preflight requests with adequate access-control-* respond headers
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
// https://fetch.spec.whatwg.org/#cors-protocol-and-credentials
func CORS(options ...Options) flamego.Handler {
	opt := prepareOptions(options)
	return func(ctx flamego.Context, log *log.Logger) {
		reqOptions := ctx.Request().Method == http.MethodOptions

		headers := map[string]string{
			"Access-Control-Allow-Methods": strings.Join(opt.Methods, ","),
			"Access-Control-Allow-Headers": ctx.Request().Header.Get("Access-Control-Request-Headers"),
			"Access-Control-Max-Age":       strconv.Itoa(opt.MaxAgeSeconds),
		}
		if opt.AllowDomain[0] == "*" {
			headers["Access-Control-Allow-Origin"] = "*"
		} else {
			origin := ctx.Request().Header.Get("Origin")
			if reqOptions && origin == "" {
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
				if u.Hostname() == d || (opt.AllowSubdomain && strings.HasSuffix(u.Hostname(), "."+d)) || d == anyDomain {
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
			}
			if reqOptions && !ok {
				http.Error(ctx.ResponseWriter(), fmt.Sprintf("CORS request from prohibited domain %v", origin), http.StatusBadRequest)
				return
			}
		}

		ctx.ResponseWriter().Before(func(w flamego.ResponseWriter) {
			for k, v := range headers {
				w.Header().Set(k, v)
			}
		})

		if reqOptions {
			ctx.ResponseWriter().WriteHeader(http.StatusOK) // return response
			return
		}
	}
}
