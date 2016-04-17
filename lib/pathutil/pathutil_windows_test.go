package pathutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreNone_SufxNone_FrwdSlash(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("", defaultf.Suffix())
	assert.Equal("", defaultf.Prefix())
	assert.Equal("/", defaultf.Separator())
	assert.Equal(PrefixNone, defaultf.pre)
	assert.Equal(SuffixNone, defaultf.sufx)
	assert.Equal(FrwdSlash, defaultf.sep)

	wd := "D:/Dev/gopath/src/github.com/Felamande/god/lib/pathutil"
	assert.Equal(wd, defaultf.Wd())

	path := "D:\\dir\\dir2\\file"
	assert.Equal("D:/dir/dir2/file", defaultf.replaceSep(path))

	path = "D:\\dir/dir2\\file"
	assert.Equal("D:/dir/dir2/file", defaultf.ReplaceSep(path))

	path = `D:/dir\file/`
	assert.Equal(`D:/dir/file`, defaultf.FormatRel(path))

	path = "./dir\\file/"
	assert.Equal("dir/file", defaultf.FormatRel(path))

	path = "D:/Dev/gopath/src/github.com/Felamande/god/lib/pathutil/dir/file"
	rel, err := defaultf.getRel(path)
	assert.Equal(err, nil)
	assert.Equal("dir/file", rel)

	path = "D:/Dev/gopath/src/github.com/Felamande/god/lib/pathutil/dir/file/"
	rel, err = defaultf.getRel(path)
	assert.Equal(err, nil)
	assert.Equal("dir/file", rel)

	path = "D:/Dev/gopath/src/github.com/Felamande/god/lib/dir/file/"
	rel, err = defaultf.getRel(path)
	assert.NotEqual(err, nil)
	assert.Equal("", rel)

	path = "D:/Dev/gopath/src/github.com/Felamande/god/lib/pathutil"
	rel, err = defaultf.getRel(path)
	assert.Equal(err, nil)
	assert.Equal(".", rel)

	path = "D:/Dev/gopath/src/github.com/Felamande/god/lib/pathutil/"
	rel, err = defaultf.getRel(path)
	assert.Equal(err, nil)
	assert.Equal(".", rel)

	path = "D:/Dev/gopath/src/github.com/Felamande/god/lib/pathutil/dir/file/"
	abs, rel, err := GetPath(path)
	assert.Equal("D:/Dev/gopath/src/github.com/Felamande/god/lib/pathutil/dir/file", abs)
	assert.Equal("dir/file", rel)
	assert.Equal(err, nil)

	path = "./dir/file/"
	abs, rel, err = GetPath(path)
	assert.Equal("D:/Dev/gopath/src/github.com/Felamande/god/lib/pathutil/dir/file", abs)
	assert.Equal("dir/file", rel)
	assert.Equal(err, nil)

	assert.Equal(false, defaultf.isParentPrefix(path))
}

func TestPreDS_SufxNone_FrwdSlash(t *testing.T) {
	SetPrefix(PrefixDotSlash)
	assert := assert.New(t)
	assert.Equal("", defaultf.Suffix())
	assert.Equal("./", defaultf.Prefix())
	assert.Equal("/", defaultf.Separator())
	assert.Equal(PrefixDotSlash, defaultf.pre)
	assert.Equal(SuffixNone, defaultf.sufx)
	assert.Equal(FrwdSlash, defaultf.sep)

	wd := "D:/Dev/gopath/src/github.com/Felamande/god/lib/pathutil"
	assert.Equal(wd, defaultf.Wd())

	path := "D:\\dir\\dir2\\file"
	assert.Equal("D:/dir/dir2/file", defaultf.replaceSep(path))

	path = "D:\\dir/dir2\\file"
	assert.Equal("D:/dir/dir2/file", defaultf.ReplaceSep(path))

	path = `D:/dir\file/`
	assert.Equal(`D:/dir/file`, defaultf.FormatRel(path))

	path = "dir\\file/"
	assert.Equal("./dir/file", defaultf.FormatRel(path))

	path = "D:/Dev/gopath/src/github.com/Felamande/god/lib/pathutil/dir/file"
	rel, err := defaultf.getRel(path)
	assert.Equal(err, nil)
	assert.Equal("./dir/file", rel)

	path = "D:/Dev/gopath/src/github.com/Felamande/god/lib/pathutil/dir/file/"
	rel, err = defaultf.getRel(path)
	assert.Equal(err, nil)
	assert.Equal("./dir/file", rel)

	path = "D:/Dev/gopath/src/github.com/Felamande/god/lib/dir/file/"
	rel, err = defaultf.getRel(path)
	assert.NotEqual(err, nil)
	assert.Equal("", rel)

	path = "D:/Dev/gopath/src/github.com/Felamande/god/lib/pathutil"
	rel, err = defaultf.getRel(path)
	assert.Equal(err, nil)
	assert.Equal(".", rel)

	path = "D:/Dev/gopath/src/github.com/Felamande/god/lib/pathutil/"
	rel, err = defaultf.getRel(path)
	assert.Equal(err, nil)
	assert.Equal("./", rel)

	path = "D:/Dev/gopath/src/github.com/Felamande/god/lib/pathutil/dir/file"
	abs, rel, err := GetPath(path)
	assert.Equal("D:/Dev/gopath/src/github.com/Felamande/god/lib/pathutil/dir/file", abs)
	assert.Equal("./dir/file", rel)
	assert.Equal(err, nil)

	path = "dir/file"
	abs, rel, err = GetPath(path)
	assert.Equal("D:/Dev/gopath/src/github.com/Felamande/god/lib/pathutil/dir/file", abs)
	assert.Equal("./dir/file", rel)
	assert.Equal(err, nil)

	assert.Equal(false, defaultf.isParentPrefix(path))
}
