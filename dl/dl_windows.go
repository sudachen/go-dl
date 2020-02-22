// +build windows,amd64

package dl

import "C"
import (
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func loadLibrary(dlname string) (uintptr, error) {
	h, err := syscall.LoadLibrary(dlname)
	return uintptr(h), err
}

func bindFunction(h uintptr, funcname string, p unsafe.Pointer) (err error) {
	addr, err := syscall.GetProcAddress(syscall.Handle(h),funcname)
	if err == nil {
		q := (*uintptr)(p)
		(*q) = addr
	}
	return
}

func expandCache(s string) string {
	return os.Getenv("localappdata") + "\\.cache\\" + strings.ReplaceAll(s,"/","\\")
}
