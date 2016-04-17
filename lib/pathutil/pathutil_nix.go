// +build !windows

package pathutil

func (f *Formatter) isBadPrefix(path string) bool {
	return false
}
