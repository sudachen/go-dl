package dl

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/sudachen/go-fp/fu"
	"github.com/ulikunitz/xz"
	"golang.org/x/xerrors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"unsafe"
)

type SO struct {
	dlHandle unsafe.Pointer
	verbose  Verbose
	onerror  OnError
}

func (so SO) Ok() bool {
	return so.dlHandle != nil
}

func Load(a ...interface{}) SO {
	so := SO{
		verbose: fu.IfsOption(Verbose(func(text string, vl int) {
			if vl == 0 {
				fmt.Println(text)
			}
		}), a).(Verbose),
		onerror: fu.IfsOption(OnError(func(err error) {
			panic(err.Error())
		}), a).(OnError),
	}

	cached := ""

	for _, x := range a {
		switch q := x.(type) {
		case System, Custom:
			s := reflect.ValueOf(q).String()
			so.verbose(fmt.Sprintf("trying to load `%v`", s), 2)
			so.dlHandle, _ = loadLibrary(s)
			if so.dlHandle != nil {
				return so
			}
		case Cached:
			cached = expandCache(reflect.ValueOf(q).String())
			so.verbose(fmt.Sprintf("trying to load `%v`", cached), 2)
			so.dlHandle, _ = loadLibrary(cached)
			if so.dlHandle != nil {
				return so
			}
		}
	}

	external, zt := fu.StrMultiOption(a, External(""), GzipExternal(""), LzmaExternal(""))
	if external != "" {

		bf := bytes.Buffer{}

		so.verbose(fmt.Sprintf("downloading `%v`", external), 1)
		if resp, err := http.Get(external); err != nil {
			so.onerror(err)
			return SO{}
		} else {
			defer resp.Body.Close()
			if _, err = io.Copy(&bf, resp.Body); err != nil {
				so.onerror(err)
				return SO{}
			}
		}

		if zt > 0 { // compressed
			so.verbose("unpacking", 2)
			bx := bf
			bf = bytes.Buffer{}
			err := error(nil)
			switch zt {
			case 1: // GzipExternal
				var z io.ReadCloser
				if z, err = gzip.NewReader(&bx); err == nil {
					defer z.Close()
					_, err = io.Copy(&bf, z)
				}
			case 2: // LzmaExternal
				var z io.Reader
				if z, err = xz.NewReader(&bx); err == nil {
					_, err = io.Copy(&bf, z)
				}
				//default:
				//	panic("unknown compression")
			}
			if err != nil {
				so.onerror(err)
				return SO{}
			}
		}

		so.verbose(fmt.Sprintf("caching as `%v`", cached), 2)
		_ = os.MkdirAll(filepath.Dir(cached), 0755)
		if err := ioutil.WriteFile(cached, bf.Bytes(), 0644); err != nil {
			so.onerror(err)
			return SO{}
		}

		so.verbose(fmt.Sprintf("trying to load `%v`", cached), 2)
		so.dlHandle, _ = loadLibrary(cached)
		if so.dlHandle != nil {
			return so
		}
	}
	so.onerror(xerrors.Errorf("not found or failed to load dynamic library"))
	return SO{}
}

/*

	//
	//int function(int);
	//
	//#define DEFINE_JUMPER(x) \
	//        void *_godl_##x = (void*)0; \
	//        __asm__(".global "#x"\n\t"#x":\n\tmovq _godl_"#x"(%rip),%rax\n\tjmp *%rax\n")
	//
	//DEFINE_JUMPER(function)
	//
	import "C"

	func init() {
		if runtime.GOOS == "linux" && runtime.GOARCH == "amd64"{
			so := dl.Load(
				dl.Cache("dl/go-dl/libfunction.so"),
				dl.LzmaExternal("https://github.com/sudachen/go-dl/releases/download/initial/libfunction_lin64.lzma"))
		} else if runtime.GOOS == "windows" && runtime.GOARCH == "amd64" {
			so := dl.Load(
				dl.Cache("dl/go-dl/function.dll"),
				dl.LzmaExternal("https://github.com/sudachen/go-dl/releases/download/initial/libfunction_win64.lzma"))
		}
		so.Bind("function",unsafe.Pointer(&C._godl_function))
	}

	func main() {
		C.function(0)
	}
*/
func (so SO) Bind(funcname string, ptrptr unsafe.Pointer) {
	if !so.Ok() {
		so.onerror(xerrors.Errorf("dynamic library object is not initialized"))
		return
	}
	so.verbose(fmt.Sprintf("binding SO.'%v' to *(%v)", funcname, ptrptr), 2)
	if err := bindFunction(so.dlHandle, funcname, ptrptr); err != nil {
		so.onerror(err)
	}
}

type LzmaExternal string
type GzipExternal string
type External string
type Verbose func(string, int)
type OnError func(error)
type Cached string
type System string
type Custom string
