package dl

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/sudachen/go-fp/fu"
	"github.com/ulikunitz/xz"
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
	verbose Verbose
	onerror OnError
}

func (so SO) Verbose(text string, verbosity int) {
	if so.verbose != nil {
		so.verbose(text,verbosity)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path);
	return err == nil
}

func Load(a ...interface{}) SO {
	so := SO{
		verbose: fu.IfsOption(Verbose(nil),a).(Verbose),
		onerror: fu.IfsOption(OnError(nil),a).(OnError),
	}

	cached := ""

	for _,x := range a {
		switch q := x.(type) {
		case System,Custom:
			s := reflect.ValueOf(q).String()
			so.Verbose(fmt.Sprintf("trying to load `%v`",s),2)
			so.dlHandle, _ = loadLibrary(s)
			if so.dlHandle != nil {
				return so
			}
		case Cached:
			cached = expandCache(reflect.ValueOf(q).String())
			so.Verbose(fmt.Sprintf("trying to load `%v`",cached),2)
			so.dlHandle, _ = loadLibrary(cached)
			if so.dlHandle != nil {
				return so
			}
		}
	}

	external, zt := fu.StrMultiOption(a,External(""),GzipExternal(""),LzmaExternal(""))
	if external != "" {

		bf := bytes.Buffer{}

		so.Verbose(fmt.Sprintf("downloading `%v`",external),1)
		if resp, err := http.Get(external); err != nil {
			panic(err.Error())
		} else {
			defer resp.Body.Close()
			if _,err = io.Copy(&bf,resp.Body); err != nil {
				panic(err.Error())
			}
		}

		if zt > 0 { // compressed
			so.Verbose("unpacking",2)
			bx := bf
			bf = bytes.Buffer{}
			switch zt {
			case 1: // GzipExternal
				z, err := gzip.NewReader(&bx);
				if err == nil {
					defer z.Close()
					_, err = io.Copy(&bf,z)
				}
				if err != nil {
					panic(err.Error())
				}
			case 2: // LzmaExternal
				z, err := xz.NewReader(&bx)
				if err == nil {
					_, err = io.Copy(&bf,z)
				}
				if err != nil {
					panic(err.Error())
				}
			//default:
			//	panic("unknown compression")
			}
		}

		so.Verbose(fmt.Sprintf("caching as `%v`",cached),2)
		_ = os.MkdirAll(filepath.Dir(cached),0755)
		if err := ioutil.WriteFile(cached, bf.Bytes(),0644); err != nil {
			panic(err.Error())
		}

		so.Verbose(fmt.Sprintf("trying to load `%v`",cached),2)
		so.dlHandle, _ = loadLibrary(cached)
		if so.dlHandle != nil {
			return so
		}
	}
	panic("not found or failed to load dynamic library")
}

func (so SO) Bind(funcname string, ptrptr unsafe.Pointer) {
	so.Verbose(fmt.Sprintf("binding SO.'%v' to *(%v)",funcname,ptrptr),2)
	err := bindFunction(so.dlHandle, funcname, ptrptr)
	if err != nil {
		if so.onerror == nil {
			panic(err.Error())
		}
		so.onerror(err)
	}
}

type LzmaExternal string
type GzipExternal string
type External string
type Verbose func(string,int)
type OnError func(error)
type Cached string
type System string
type Custom string

