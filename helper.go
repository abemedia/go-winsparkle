//go:build windows

package winsparkle

import (
	"syscall"
	"unsafe"
)

func char(s string) uintptr {
	b, err := syscall.BytePtrFromString(s)
	if err != nil {
		panic(err)
	}
	return uintptr(unsafe.Pointer(b))
}

func wchar(s string) uintptr {
	i, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		panic(err)
	}
	return uintptr(unsafe.Pointer(i))
}

func boolean(b bool) uintptr {
	if b {
		return 1
	}
	return 0
}

func utf8PtrToString(p *uint8) string {
	if p == nil || *p == 0 {
		return ""
	}

	// Find NUL terminator.
	n := 0
	for ptr := unsafe.Pointer(p); *(*uint8)(ptr) != 0; n++ {
		ptr = unsafe.Add(ptr, unsafe.Sizeof(*p))
	}
	return string(unsafe.Slice(p, n))
}

func utf16PtrToString(p *uint16) string {
	if p == nil || *p == 0 {
		return ""
	}

	// Find NUL terminator.
	n := 0
	for ptr := unsafe.Pointer(p); *(*uint16)(ptr) != 0; n++ {
		ptr = unsafe.Add(ptr, unsafe.Sizeof(*p))
	}
	return syscall.UTF16ToString(unsafe.Slice(p, n))
}

func configMethods(cs ConfigStore) unsafe.Pointer {
	if cs == nil {
		return nil
	}
	return unsafe.Pointer(&struct{ read, write, delete, _ uintptr }{
		read: syscall.NewCallbackCDecl(func(name *uint8, buf *uint16, size uint, _ uintptr) uintptr {
			s, ok := cs.Read(utf8PtrToString(name))
			if !ok {
				return 0
			}
			b := (*[1 << 29]uint8)(unsafe.Pointer(buf))[:size:size]
			copy(b, s)
			return 1
		}),
		write: syscall.NewCallbackCDecl(func(name *uint8, value *uint16, _ uintptr) uintptr {
			return boolean(cs.Write(utf8PtrToString(name), utf16PtrToString(value)))
		}),
		delete: syscall.NewCallbackCDecl(func(name *uint8, _ uintptr) uintptr {
			return boolean(cs.Delete(utf8PtrToString(name)))
		}),
	})
}
