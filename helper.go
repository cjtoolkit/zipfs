package zipfs

import (
	"archive/zip"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Initialise ZipFS based on given zip file name.
// If the file does not exist, it will try to get the zip file that is embedded in the application itself.
// If the application also does not have zip embedded it will panic.
func InitZipFs(zipFileName string) http.FileSystem {
	f, err := os.Open(zipFileName)
	if err != nil {
		return initEmbeddedZipFs()
	}
	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return initEmbeddedZipFs()
	}

	z, err := zip.NewReader(f, fi.Size())
	if err == nil {
		return NewZipFSWithReaderAt(z, f)
	}

	return initEmbeddedZipFs()
}

func initEmbeddedZipFs() http.FileSystem {
	z, r, err := GetEmbeddedZip()
	if err != nil {
		log.Panic(err)
	}

	return NewZipFSWithReaderAt(z, r)
}

// Init Zip FS from HTTP File, must be uncompressed. Does not support compressed files!
func InitZipFsFromHttpFile(f http.File) http.FileSystem {
	r, ok := f.(io.ReaderAt)
	if !ok {
		log.Panic("Does not implemented io.ReaderAt, must use uncompressed file. Does not support compressed files!")
	}
	fi, err := f.Stat()
	if err != nil {
		log.Panic(err)
	}

	z, err := zip.NewReader(r, fi.Size())
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

// Must not have any error, in HTTP File.
func Must(f http.File, err error) http.File {
	if err != nil {
		log.Panic(err)
	}
	return f
}
