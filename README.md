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
