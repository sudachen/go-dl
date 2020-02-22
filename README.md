
[![CircleCI](https://circleci.com/gh/sudachen/go-dl.svg?style=svg)](https://circleci.com/gh/sudachen/go-dl)
[![Maintainability](https://api.codeclimate.com/v1/badges/3b8d5bd3fe992a6ce7f2/maintainability)](https://codeclimate.com/github/sudachen/go-dl/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/3b8d5bd3fe992a6ce7f2/test_coverage)](https://codeclimate.com/github/sudachen/go-dl/test_coverage)
[![Go Report Card](https://goreportcard.com/badge/github.com/sudachen/go-dl)](https://goreportcard.com/report/github.com/sudachen/go-dl)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

```golang
/*

int function(int);

#define DEFINE_JUMPER(x) \
        void *_godl_##x = (void*)0; \
        __asm__(".global "#x"\n\t"#x":\n\tmovq _godl_"#x"(%rip),%rax\n\tjmp *%rax\n")
  
DEFINE_JUMPER(function);

*/
import "C"

import (
	"github.com/sudachen/go-dl/dl"
	"runtime"
	"unsafe"
)

func init() {
    urlbase := "https://github.com/sudachen/go-dl/releases/download/initial/"
    if runtime.GOOS == "linux" && runtime.GOARCH == "amd64"{
        so := dl.Load(
            dl.Cache("dl/go-dl/libfunction.so"),
            dl.LzmaExternal(urlbase+"libfunction_lin64.lzma"))
    } else if runtime.GOOS == "windows" && runtime.GOARCH == "amd64" {
        so := dl.Load(
            dl.Cache("dl/go-dl/function.dll"),
            dl.LzmaExternal(urlbase+"libfunction_win64.lzma"))
    }
    so.Bind("function",unsafe.Pointer(&C._godl_function))
}

func main() {
    C.function(0)
}
```
