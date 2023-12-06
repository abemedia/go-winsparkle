//go:build windows
// +build windows

package dll

import _ "embed"

//go:embed x86/WinSparkle.dll
var dll []byte
