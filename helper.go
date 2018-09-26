package zipfs

import (
	"archive/zip"
	"log"
	"net/http"
	"os"
	"strings"
)

func InitZipFs(zipFileName string) http.FileSystem {
	{
		f, err := os.Open(zipFileName)
		if err != nil {
			log.Panic(err)
		}
		fi, err := f.Stat()
		if err != nil {
			f.Close()
			log.Panic(err)
		}

		z, err := zip.NewReader(f, fi.Size())
		if err == nil {
			return NewZipFSWithReadAt(z, f)
		}
	}

	z, r, err := GetEmbeddedZip()
	if err != nil {
		log.Panic(err)
	}

	return NewZipFSWithReadAt(z, r)
}

type fileSystemFunc func(name string) (http.File, error)

func (fn fileSystemFunc) Open(name string) (http.File, error) { return fn(name) }

func Prefix(prefix string, fileSystem http.FileSystem) http.FileSystem {
	return fileSystemFunc(func(name string) (http.File, error) {
		return fileSystem.Open(strings.TrimRight(
			strings.TrimRight(prefix, "/")+"/"+strings.TrimLeft(name, "/"), "/"))
	})
}
