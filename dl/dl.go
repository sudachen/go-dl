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
			if so.dlHandle != 0 {
				return so
			}
		case Cached:
			cached = expandCache(reflect.ValueOf(q).String())
			so.verbose(fmt.Sprintf("trying to load `%v`", cached), 2)
			so.dlHandle, _ = loadLibrary(cached)
			if so.dlHandle != 0 {
				return so
			}
		}
	}

	if cached != "" {
		if preload(cached,a...) {
			so.verbose(fmt.Sprintf("trying to load `%v`", cached), 2)
			so.dlHandle, _ = loadLibrary(cached)
			if so.dlHandle != 0 {
				return so
			}
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

func (so SO) Ok() bool {
	return so.dlHandle != 0
}

type SO struct {
	dlHandle uintptr
	verbose  Verbose
	onerror  OnError
}

type LzmaExternal string
type GzipExternal string
type External string
type Verbose func(string, int)
type OnError func(error)
type Cached string
type System string
type Custom string

func (c Custom) Preload(a ...interface{}) {
	preload(string(c),a...)
}

func (c Cached) Preload(a ...interface{}) {
	preload(expandCache(string(c)),a...)
}

func preload(sopath string, a ...interface{}) (ok bool) {
	verbose := fu.IfsOption(Verbose(func(text string, vl int) {
		if vl == 0 {
			fmt.Println(text)
		}
	}), a).(Verbose)
	onerror := fu.IfsOption(OnError(func(err error) {
		panic(err.Error())
	}), a).(OnError)

	external, zt := fu.StrMultiOption(a, External(""), GzipExternal(""), LzmaExternal(""))
	if external != "" {

		bf := bytes.Buffer{}

		verbose(fmt.Sprintf("downloading `%v`", external), 1)
		if resp, err := http.Get(external); err == nil {
			defer resp.Body.Close()
			if _, err = io.Copy(&bf, resp.Body); err == nil {
				if zt > 0 { // compressed
					bx := bf
					bf = bytes.Buffer{}
					switch zt {
					case 1: // GzipExternal
						verbose("unpacking gzip", 2)
						var z io.ReadCloser
						if z, err = gzip.NewReader(&bx); err == nil {
							defer z.Close()
							_, err = io.Copy(&bf, z)
						}
					case 2: // LzmaExternal
						verbose("unpacking lzma", 2)
						var z io.Reader
						if z, err = xz.NewReader(&bx); err == nil {
							_, err = io.Copy(&bf, z)
						}
						//default:
						//	panic("unknown compression")
					}
				}
				if err == nil {
					verbose(fmt.Sprintf("caching as `%v`", sopath), 2)
					_ = os.MkdirAll(filepath.Dir(sopath), 0755)
					err = ioutil.WriteFile(sopath, bf.Bytes(), 0644)
				}
			}
			if err != nil {
				onerror(err)
				return
			}
		}
		ok = true
	}
	return
}

func (c Cached) Remove() (err error) {
	s := expandCache(string(c))
	_, err = os.Stat(s);
	if err == nil {
		err = os.Remove(s)
	}
	return nil
}
