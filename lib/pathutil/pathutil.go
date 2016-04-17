package pathutil

import (
	"errors"
	"os"
	pathp "path"
	"path/filepath"
	"strings"
)

type prefix int

const (
	PrefixNone prefix = iota
	PrefixDotSlash
	prefixParent
	// prefixBad
)

type suffix int

const (
	SuffixNone suffix = iota
	SuffixSlash
)

type separator string

const (
	FrwdSlash separator = "/"
	BackSlash separator = "\\"
)

func (s separator) Sep() string {
	return string(s)
}
func (s separator) String() string {
	return string(s)
}

type Formatter struct {
	sep  separator
	pre  prefix
	sufx suffix
}

type Conf struct {
	Separator separator
	Prefix    prefix
	Suffix    suffix
}

func New(conf Conf) *Formatter {

	f := &Formatter{
		sep:  conf.Separator,
		pre:  conf.Prefix,
		sufx: conf.Suffix,
	}

	if len(f.sep) == 0 {
		f.sep = FrwdSlash
	}

	return f
}

func (f *Formatter) sufStr(s suffix) string {
	switch s {
	case SuffixNone:
		return ""
	case SuffixSlash:
		return f.sep.Sep()
	}
	return ""
}

func (f *Formatter) getPreStr(p prefix) string {

	switch p {
	case PrefixNone:
		return ""
	case PrefixDotSlash:
		return "." + f.sep.Sep()
	default:
		return ""
	}

}

func (f *Formatter) Prefix() string {

	return f.getPreStr(f.pre)
}

func (f *Formatter) Suffix() string {
	return f.sufStr(f.sufx)
}

func (f *Formatter) Separator() string {
	return f.sep.Sep()
}

func (f *Formatter) Wd() string {
	oswd, _ := os.Getwd()
	oswd = f.replaceSep(oswd)

	return f.toSuffix(oswd)
}

func (f *Formatter) replaceSep(path string) (newpath string) {
	switch f.sep {
	case FrwdSlash:
		newpath = filepath.ToSlash(path)
	case BackSlash:
		newpath = filepath.FromSlash(path)
	default:
		newpath = path
	}
	return
}

func (f *Formatter) ReplaceSep(path string) (newpath string) {
	switch f.sep {
	case FrwdSlash:
		newpath = strings.Replace(path, "\\", f.sep.Sep(), -1)
	case BackSlash:
		newpath = strings.Replace(path, "/", f.sep.Sep(), -1)
	}
	return
}

func (f *Formatter) FormatRel(path string) string {
	path = filepath.Clean(path)
	path = f.replaceSep(path)

	if filepath.IsAbs(path) {
		return f.toSuffix(path)
	}

	return f.formatRel(path)
}

func (f *Formatter) formatRel(path string) string {

	if f.isBadPrefix(path) {
		path = strings.TrimLeft(path, f.sep.Sep())
	}

	return f.toPrefix(path)
}

//path must be absolute and has the right sep
func (f *Formatter) getRel(path string) (string, error) {
	wd := f.Wd()
	wdLen, pLen := len(wd), len(path)
	if wd == path {
		return ".", nil
	}
	if pLen < wdLen || path[0:wdLen] != wd {
		return "", errors.New("the path is not in the work directory")
	}

	p := path[wdLen:]

	if f.isBadPrefix(string(p)) {
		p = p[1:]
	}

	p = f.fix(p)
	if len(p) == 0 {
		p += "."
	}
	return p, nil

}

func (f *Formatter) RelToWd(path string) (string, error) {
	path = filepath.Clean(path)
	path = f.replaceSep(path)
	if filepath.IsAbs(path) {
		return f.getRel(path)
	}

	if f.isBadPrefix(path) {
		path = path[1:]
	}
	return f.fix(path), nil

}

func (f *Formatter) GetPath(path string) (abs, rel string, err error) {
	path = filepath.Clean(path)
	path = f.replaceSep(path)

	if filepath.IsAbs(path) {
		rel, err := f.getRel(path)
		return path, rel, err
	}
	return f.toSuffix(pathp.Join(f.Wd(), path)), f.fix(path), nil
}

func (f *Formatter) Format(path string) string {
	path = filepath.Clean(path)
	path = f.replaceSep(path)
	if filepath.IsAbs(path) {
		return f.toSuffix(path)
	}

	return f.fix(path)
}

//path must be a relative path
func (f *Formatter) toPrefix(path string) string {

	pathPrefix := f.prefixOf(path)
	if f.isParentPrefix(path) || pathPrefix == f.pre {
		return path
	}

	if pathPrefix == PrefixNone {
		return f.Prefix() + path
	}

	return strings.Replace(path, f.getPreStr(pathPrefix), f.Prefix(), 3)
}

func (f *Formatter) suffixOf(path string) suffix {
	if path == "" {
		return SuffixNone
	}

	switch path[len(path)-1] {
	case []byte(f.sep.Sep())[0]:
		return SuffixSlash
	default:
		return SuffixNone
	}
}

func (f *Formatter) toSuffix(path string) string {
	pathsuf := f.suffixOf(path)

	if pathsuf == f.sufx {
		return path
	}
	switch pathsuf {
	case SuffixNone:
		return path + f.sep.Sep()

	case SuffixSlash:
		return strings.TrimRight(path, f.sep.Sep())
	default:
		return path
	}

}

//fix the suffix and prefix of a relative path
func (f *Formatter) fix(path string) string {
	path = f.toSuffix(path)
	return f.toPrefix(path)
}

//path must be a relative path
func (f *Formatter) prefixOf(path string) prefix {
	if strings.HasPrefix(path, "."+f.sep.Sep()) {
		return PrefixDotSlash
	}
	if strings.HasPrefix(path, ".."+f.sep.Sep()) {
		return prefixParent
	}

	return PrefixNone

}

//path must be a relative path
func (f *Formatter) isParentPrefix(path string) bool {
	return strings.HasPrefix(path, ".."+f.Separator())
}

//path must be a relative path
func (f *Formatter) isDotSlashPrefix(path string) bool {
	return strings.HasPrefix(path, "."+f.Separator())
}

var defaultf = New(Conf{FrwdSlash, PrefixNone, SuffixNone})

func SetPrefix(p prefix) bool {
	defaultf.pre = p
	return defaultf.pre == p
}

func SetSeparator(sep separator) bool {
	defaultf.sep = sep
	return defaultf.sep == sep
}

func SetSuffix(s suffix) bool {
	defaultf.sufx = s
	return defaultf.sufx == s
}

func FormatRel(path string) string {
	return defaultf.FormatRel(path)
}

func RelToWd(path string) (string, error) {
	return defaultf.RelToWd(path)

}

func Format(path string) string {

	return defaultf.Format(path)
}

func GetPath(path string) (abs, rel string, err error) {
	return defaultf.GetPath(path)
}
