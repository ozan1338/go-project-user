package upload

import (
	"path/filepath"
	"runtime"
)

func GetFolderName() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable get the current foldername")
	}

	dirname := filepath.Dir(filename)
	return dirname
}