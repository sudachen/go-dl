package tests

import (
	"github.com/sudachen/go-dl/dl"
	"gotest.tools/assert"
	"testing"
)

const localLibSoName = "/tmp/go-dl-test/"+libSoName
const externalLibSoLzma = "https://github.com/sudachen/go-dl/releases/download/initial/"+libSoLzma
const externalLibSoGzip = "https://github.com/sudachen/go-dl/releases/download/initial/"+libSoGzip
const externalLibSo = "https://github.com/sudachen/go-dl/releases/download/initial/"+libSoName

func init() {
	dl.Custom(localLibSoName).Preload(
		dl.LzmaExternal(externalLibSoLzma))
}

func Test_LoadCustom(t *testing.T) {
	so := dl.Load(
			dl.OnError(func(err error){
				assert.NilError(t,err)
			}),
			dl.Custom(localLibSoName))
	assert.Assert(t,so.Ok())
	so.Bind("function", functionPtr())
	assert.Assert(t, function(1) == 2)
}

func Test_LoadLzmaExternal(t *testing.T) {
	err := dl.Cached("dl/.go-dl-test/loadlzmaexternal.so").Remove()
	assert.NilError(t,err)
	so := dl.Load(
		dl.OnError(func(err error){
			assert.NilError(t,err)
		}),
		dl.Cached("dl/.go-dl-test/loadlzmaexternal.so"),
		dl.LzmaExternal(externalLibSoLzma))
	assert.Assert(t,so.Ok())
	*(*uintptr)(functionPtr()) = 0
	so.Bind("function", functionPtr())
	assert.Assert(t, function(1) == 2)
}

func Test_LoadGzipExternal(t *testing.T) {
	err := dl.Cached("dl/.go-dl-test/loadgzipexternal"+SoExt).Remove()
	assert.NilError(t,err)
	so := dl.Load(
		dl.OnError(func(err error){
			assert.NilError(t,err)
		}),
		dl.Cached("dl/.go-dl-test/loadgzipexternal"+SoExt),
		dl.GzipExternal(externalLibSoGzip))
	assert.Assert(t,so.Ok())
	*(*uintptr)(functionPtr()) = 0
	so.Bind("function", functionPtr())
	assert.Assert(t, function(1) == 2)
}

func Test_LoadUncompressedExternal(t *testing.T) {
	err := dl.Cached("dl/.go-dl-test/loadexternal.so").Remove()
	assert.NilError(t,err)
	so := dl.Load(
		dl.OnError(func(err error){
			assert.NilError(t,err)
		}),
		dl.Cached("dl/.go-dl-test/loadexternal"+SoExt),
		dl.External(externalLibSo))
	assert.Assert(t,so.Ok())
	*(*uintptr)(functionPtr()) = 0
	so.Bind("function", functionPtr())
	assert.Assert(t, function(1) == 2)
}

func Test_LoadCached(t *testing.T) {
	err := dl.Cached("dl/go-dl/"+libSoName).Remove()
	assert.NilError(t,err)
	dl.Cached("dl/go-dl/"+libSoName).Preload(
		dl.LzmaExternal(externalLibSoLzma),
		dl.OnError(func(err error){
			assert.NilError(t,err)
		}))
	so := dl.Load(
		dl.OnError(func(err error){
			assert.NilError(t,err)
		}),
		dl.Cached("dl/go-dl/"+libSoName))
	assert.Assert(t,so.Ok())
	*(*uintptr)(functionPtr()) = 0
	so.Bind("function", functionPtr())
	assert.Assert(t, function(1) == 2)
}

