package zipfs

import (
	"archive/zip"
	"log"
	"net/http"
	"strings"
)

func InitZipFs(zipFileName string) http.FileSystem {
	{
		z, err := zip.OpenReader(zipFileName)
		if err == nil {
			return NewZipFS(&z.Reader)
		}
	}

	z, err := GetEmbeddedZip()
	if err != nil {
		log.Panic(err)
	}

	return NewZipFS(z)
}

type fileSystemFunc func(name string) (http.File, error)

func (fn fileSystemFunc) Open(name string) (http.File, error) { return fn(name) }

func Prefix(prefix string, fileSystem http.FileSystem) http.FileSystem {
	return fileSystemFunc(func(name string) (http.File, error) {
		return fileSystem.Open(strings.TrimRight(prefix, "/") + "/" + strings.TrimLeft(name, "/"))
	})
}
