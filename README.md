# WinSparkle Go Bindings

[![Go Reference](https://pkg.go.dev/badge/github.com/abemedia/go-winsparkle.svg)](https://pkg.go.dev/github.com/abemedia/go-winsparkle)

This package provides go bindings for [WinSparkle](https://github.com/vslavik/winsparkle) created by
Vaclav Slavik.

WinSparkle is a plug-and-forget software update library for Windows applications. It is heavily
inspired by the Sparkle framework for OS X written by Andy Matuschak and others, to the point of
sharing the same updates format (appcasts) and having very similar user interface.

See https://winsparkle.org for more information about WinSparkle.

## Documentation

See the [WinSparkle wiki](https://github.com/vslavik/winsparkle/wiki) and the
[GoDoc](https://pkg.go.dev/github.com/abemedia/go-winsparkle?tab=doc).

## Important

WinSparkle.dll must be placed into the same directory as your app executable. Depending on your
architecture use the version from `dll/x64` or `dll/x86`.
