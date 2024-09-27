package tusuploader

import (
	"os"

	"github.com/tus/tusd/v2/pkg/filestore"
)

func NewLocalStore(path string) *filestore.FileStore {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	return &filestore.FileStore{
		Path: path,
	}
}
