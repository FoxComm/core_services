package test

import (
	"path/filepath"
	"runtime"
	"strings"
)

func RouterRoot() string {
	_, file, _, _ := runtime.Caller(0)
	filelist := strings.Split(filepath.Clean(file), string(filepath.Separator))
	// return / at the beginning of path
	filelist[0] = "/" + filelist[0]

	l := len(filelist) - 3
	return filepath.Join(filelist[:l]...)
}
