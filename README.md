# WinSparkle Go Bindings

[![Go Reference](https://pkg.go.dev/badge/github.com/abemedia/go-winsparkle.svg)](https://pkg.go.dev/github.com/abemedia/go-winsparkle)

This package provides go bindings for [WinSparkle](https://github.com/vslavik/winsparkle) created by
Vaclav Slavik.

WinSparkle is a plug-and-forget software update library for Windows applications. It is heavily
inspired by the Sparkle framework for MacOS written by Andy Matuschak and others, to the point of
sharing the same updates format (appcasts) and having very similar user interface.

See <https://winsparkle.org> for more information about WinSparkle.

## Documentation

See the [WinSparkle wiki](https://github.com/vslavik/winsparkle/wiki) and the
[GoDoc](https://pkg.go.dev/github.com/abemedia/go-winsparkle?tab=doc).

## Important

WinSparkle.dll must be placed into the same directory as your app executable. Depending on your
architecture use the version from [dll/x64](./dll/x64/), [dll/x86](./dll/x86/) or
[dll/arm64](./dll/arm64/).

Alternatively you can embed the DLL into your application by importing
`github.com/abemedia/go-winsparkle/dll`.

## Example

```go
package main

import (
	"github.com/abemedia/go-winsparkle"
	_ "github.com/abemedia/go-winsparkle/dll" // Embed DLL.
)

//go:embed dsa-public-key.pem
var dsaPublicKey string

func main() {
	winsparkle.SetAppcastURL("https://dl.example.com/appcast.xml")
	winsparkle.SetAppDetails("example.com", "My Cool App", "1.0.0")
	winsparkle.SetAutomaticCheckForUpdates(true)

	if err := winsparkle.SetDSAPubPEM(dsaPublicKey); err != nil {
		panic(err)
	}

	winsparkle.Init()
	defer winsparkle.Cleanup()

	runMyApp()
}
```

## Versions

The version for `go-winsparkle` corresponds to the WinSparkle version. If you are not embedding the
DLL by importing `github.com/abemedia/go-winsparkle/dll` please make sure that the version of
`go-winsparkle` is the same as that of the DLL file or some functions might not work.

## Caveats

WinSparkle only runs on Windows. For MacOS see <https://github.com/abemedia/go-sparkle>.
