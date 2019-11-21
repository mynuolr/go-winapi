package winapi

import (
	"fmt"
	"syscall"
)

func Call(dll interface{},name string, parm ...uintptr) (r1, r2 uintptr, err error) {
	switch dll {
	case USER32:
		r1,r2,err=moduser32.NewProc(name).Call(parm...)
		break
	case KERNEL32:
		r1,r2,err=modkernel32.NewProc(name).Call(parm...)
		break
	default:
		if s,ok:=dll.(string);ok {
			r1,r2,err=syscall.NewLazyDLL(s).NewProc(name).Call(parm...)
		} else {

		}
		err = fmt.Errorf(" invalid parameter ")
	}
	return
}
