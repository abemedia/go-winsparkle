//go:build windows
// +build windows

package dll

import _ "embed"

//go:embed arm64/WinSparkle.dll
var dll []byte
