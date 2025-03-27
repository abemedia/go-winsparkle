//go:build windows
// +build windows

package dll

import (
	"os"
	"path/filepath"
)

const version = "0.9.0"

func init() {
	dir := filepath.Join(os.TempDir(), "WinSparkle-"+version)
	file := filepath.Join(dir, "WinSparkle.dll")

	if _, err := os.Stat(file); err != nil {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			panic(err)
		}
		if err := os.WriteFile(file, dll, os.ModePerm); err != nil { //nolint:gosec
			panic(err)
		}
	}

	if err := os.Setenv("PATH", dir+";"+os.Getenv("PATH")); err != nil {
		panic(err)
	}
}
