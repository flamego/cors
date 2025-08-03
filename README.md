# cors

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/flamego/cors/Go?logo=github&style=for-the-badge)](https://github.com/flamego/cors/actions?query=workflow%3AGo)
[![GoDoc](https://img.shields.io/badge/GoDoc-Reference-blue?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/flamego/cors?tab=doc)

Package cors is a middleware that provides the Cross-Origin Resource Sharing for [Flamego](https://github.com/flamego/flamego).

## Installation

```zsh
go get github.com/flamego/cors
```

## Getting started

```go
package main

import (
	"github.com/flamego/cors"
	"github.com/flamego/flamego"
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

## Getting help

- Read [documentation and examples](https://flamego.dev/middleware/cors.html).
- Please [file an issue](https://github.com/flamego/flamego/issues) or [start a discussion](https://github.com/flamego/flamego/discussions) on the [flamego/flamego](https://github.com/flamego/flamego) repository.

## License

This project is under the MIT License. See the [LICENSE](LICENSE) file for the full license text.
