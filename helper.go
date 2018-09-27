package zipfs

import (
	"archive/zip"
	"log"
	"net/http"
	"os"
	"strings"
)

// Initialise ZipFS based on give zip file name.
// If the file does not exist, it will try to get the zip file that is embedded in the application itself.
// If the application also does not have zip embedded it will panic.
func InitZipFs(zipFileName string) http.FileSystem {
	f, err := os.Open(zipFileName)
	if err != nil {
		return initZipFsFromEmbed()
	}
	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return initZipFsFromEmbed()
	}

	z, err := zip.NewReader(f, fi.Size())
	if err == nil {
		return NewZipFSWithReaderAt(z, f)
	}

	return initZipFsFromEmbed()
}

func initZipFsFromEmbed() http.FileSystem {
	z, r, err := GetEmbeddedZip()
	if err != nil {
		log.Panic(err)
	}

	return NewZipFSWithReaderAt(z, r)
}

type fileSystemFunc func(name string) (http.File, error)

func (fn fileSystemFunc) Open(name string) (http.File, error) { return fn(name) }

// Add prefix to name
func Prefix(prefix string, fileSystem http.FileSystem) http.FileSystem {
	prefix = strings.TrimRight(prefix, "/")
	return fileSystemFunc(func(name string) (http.File, error) {
		return fileSystem.Open(strings.TrimRight(prefix+"/"+strings.TrimLeft(name, "/"), "/"))
	})
}
