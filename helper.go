//go:build windows
// +build windows

package winsparkle

import (
	"syscall"
	"unsafe"
)

func string2uintptr(s string) uintptr {
	b, err := syscall.BytePtrFromString(s)
	if err != nil {
		panic(err)
	}
	return uintptr(unsafe.Pointer(b))
}

func string2wchar(s string) uintptr {
	i, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		panic(err)
	}
	return uintptr(unsafe.Pointer(i))
}

func bool2uintptr(b bool) uintptr {
	if b {
		return 1
	}
	return 0
}
