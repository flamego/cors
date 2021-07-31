# cors

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/flamego/cors/Go?logo=github&style=for-the-badge)](https://github.com/flamego/cors/actions?query=workflow%3AGo)
[![Codecov](https://img.shields.io/codecov/c/gh/flamego/cors?logo=codecov&style=for-the-badge)](https://app.codecov.io/gh/flamego/cors)
[![GoDoc](https://img.shields.io/badge/GoDoc-Reference-blue?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/flamego/cors?tab=doc)
[![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?style=for-the-badge&logo=sourcegraph)](https://sourcegraph.com/github.com/flamego/cors)

Package cors is a middleware that provides the Cross-Origin Resource Sharing.

## Getting started

```go
package main
import (
	"github.com/flamego/flamego"
	"github.com/flamego/cors"
)
func main() {
	f := flamego.Classic()
	f.Use(cors.CORS())
	f.Get("/", func(c flamego.Context) string {
        return "ok"
    })
	f.Run()
}
```

## Installation

The minimum requirement of Go is **1.16**.

	go get github.com/flamego/cors

## License

This project is under the MIT License. See the [LICENSE](LICENSE) file for the full license text.
