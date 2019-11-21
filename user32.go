package winapi

import (
	"syscall"
	"unsafe"
)

var (
	moduser32 = syscall.NewLazyDLL("user32.dll")
	procFindWindowW                   = moduser32.NewProc("FindWindowW")
)
func FindWindowW(className, windowName *uint16) syscall.Handle {
	ret, _, _ := procFindWindowW.Call(
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)))

	return syscall.Handle(ret)
}
