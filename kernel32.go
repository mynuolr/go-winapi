package winapi

import (
	"runtime"
	"syscall"
	"unsafe"
)

var (
	modkernel32            = syscall.NewLazyDLL("kernel32.dll")
	procGetModuleHandle    = modkernel32.NewProc("GetModuleHandleW")
	procCreateProcessW     = modkernel32.NewProc("CreateProcessW")
	procVirtualAllocEx     = modkernel32.NewProc("VirtualAllocEx")
	procWriteProcessMemory = modkernel32.NewProc("WriteProcessMemory")
	procCreateRemoteThread = modkernel32.NewProc("CreateRemoteThread")
	procCreateMutex        = modkernel32.NewProc("CreateMutexW")
)

func Call(name string, parm ...uintptr) (r1, r2 uintptr, err error) {
	return modkernel32.NewProc(name).Call(parm...)
}
func GetModuleHandle(modulename string) syscall.Handle {
	var mn uintptr
	if modulename == "" {
		mn = 0
	} else {
		sPt, _ := syscall.UTF16PtrFromString(modulename)
		mn = uintptr(unsafe.Pointer(sPt))
	}
	ret, _, _ := procGetModuleHandle.Call(mn)
	return syscall.Handle(ret)
}

func CreateProcess(
	appName *uint16,
	commandLine *uint16,
	procSecurity *syscall.SecurityAttributes,
	threadSecurity *syscall.SecurityAttributes,
	inheritHandles bool,
	creationFlags uint32,
	env *uint16,
	currentDir *uint16,
	startupInfo *syscall.StartupInfo,
	outProcInfo *syscall.ProcessInformation) (err error) {
	var _p0 uint32
	if inheritHandles {
		_p0 = 1
	} else {
		_p0 = 0
	}
	r1, _, e1 := procCreateProcessW.Call(
		uintptr(unsafe.Pointer(appName)),
		uintptr(unsafe.Pointer(commandLine)),
		uintptr(unsafe.Pointer(procSecurity)),
		uintptr(unsafe.Pointer(threadSecurity)),
		uintptr(_p0), uintptr(creationFlags),
		uintptr(unsafe.Pointer(env)),
		uintptr(unsafe.Pointer(currentDir)),
		uintptr(unsafe.Pointer(startupInfo)),
		uintptr(unsafe.Pointer(outProcInfo)), 0, 0)
	if r1 == 0 {
		if e1.(syscall.Errno) != 0 {
			err = syscall.Errno(997)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
func VirtualAllocEx(h syscall.Handle, addr unsafe.Pointer, size, allocType, protect uint32) (uintptr, error) {
	r1, _, e1 := procVirtualAllocEx.Call(
		uintptr(h),
		uintptr(addr),
		uintptr(size),
		uintptr(allocType),
		uintptr(protect),
	)
	if int(r1) == 0 {
		return 0, e1
	}
	return r1, nil

}
func WriteProcessMemory(process syscall.Handle, addr uintptr, buf unsafe.Pointer, size uint32) (uint32, error) {
	var nLength uint32
	r1, _, e1 := procWriteProcessMemory.Call(
		uintptr(process),
		addr,
		uintptr(buf),
		uintptr(size),
		uintptr(unsafe.Pointer(&nLength)))

	if int(r1) == 0 {
		return nLength, e1
	}
	return nLength, nil
}
func CreateRemoteThread(process syscall.Handle, sa *syscall.SecurityAttributes, stackSize uint32, startAddress,
	parameter uintptr, creationFlags uint32) (syscall.Handle, uint32, error) {
	var threadId uint32
	r1, _, e1 := procCreateRemoteThread.Call(
		uintptr(process),
		uintptr(unsafe.Pointer(sa)),
		uintptr(stackSize),
		startAddress,
		parameter,
		uintptr(creationFlags),
		uintptr(unsafe.Pointer(&threadId)))
	runtime.KeepAlive(sa)
	if int(r1) == 0 {
		return syscall.InvalidHandle, 0, e1
	}
	return syscall.Handle(r1), threadId, nil
}
func CreateMutex(name string) (syscall.Handle, error) {
	n, e := syscall.UTF16PtrFromString(name)
	if e != nil {
		return 0, e
	}
	ret, _, err := procCreateMutex.Call(
		0,
		0,
		uintptr(unsafe.Pointer(n)),
	)
	switch int(err.(syscall.Errno)) {
	case 0:
		return syscall.Handle(ret), nil
	default:
		return 0, err
	}
}
