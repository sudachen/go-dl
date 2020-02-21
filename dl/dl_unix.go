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

func loadLibrary(dlname string) (unsafe.Pointer, error) {
	s := C.CString(dlname)
	defer C.free(unsafe.Pointer(s))
	h := C.dlopen(s, C.RTLD_LAZY)
	if h == nil {
		return h, xerrors.Errorf("failed to load dynamic library")
	}
	return h, nil
}

func bindFunction(h unsafe.Pointer, funcname string, p unsafe.Pointer) error {
	n := C.CString(funcname)
	defer C.free(unsafe.Pointer(n))
	fp := (*C.void)(C.dlsym(h, n))
	if fp == nil {
		return xerrors.Errorf("dynamic library does not have symbol %v", funcname)
	}
	q := (**C.void)(p)
	(*q) = fp
	return nil
}

func expandCache(s string) string {
	if usr, err := user.Current(); err != nil {
		panic(err.Error())
	} else {
		return usr.HomeDir + "/.cache/" + s
	}
}
