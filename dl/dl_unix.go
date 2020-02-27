// +build !windows

package dl

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdlib.h>
*/
import "C"

import (
	"golang.org/x/xerrors"
	"os/user"
	"unsafe"
)

func loadLibrary(dlname string) (uintptr, error) {
	s := C.CString(dlname)
	defer C.free(unsafe.Pointer(s))
	h := C.dlopen(s, C.RTLD_LAZY)
	if h == nil {
		return 0, xerrors.Errorf("failed to load dynamic library")
	}
	return uintptr(h), nil
}

func bindFunction(h uintptr, funcname string, p unsafe.Pointer) error {
	n := C.CString(funcname)
	defer C.free(unsafe.Pointer(n))
	fp := uintptr(C.dlsym(unsafe.Pointer(h), n))
	if fp == 0 {
		return xerrors.Errorf("dynamic library does not have symbol %v", funcname)
	}
	q := (*uintptr)(p)
	(*q) = fp
	return nil
}

func expandCache(s string) string {
	usr, err := user.Current()
	if err != nil {
		panic(err.Error())
	}
	return usr.HomeDir + "/.cache/" + s
}
