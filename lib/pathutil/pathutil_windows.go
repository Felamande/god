package pathutil

func (f *Formatter) isBadPrefix(path string) bool {

	return path[0] == f.sep[0]
}
