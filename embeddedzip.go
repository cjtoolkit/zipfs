package zipfs

import (
	"archive/zip"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"runtime"
)

func binself() (*os.File, error) {
	if runtime.GOOS != "windows" {
		return os.Open(os.Args[0])
	}

	bin, err := os.Open(os.Args[0])
	if err != nil && os.IsNotExist(err) {
		return os.Open(os.Args[0] + ".exe")
	}

	return bin, err
}

func GetEmbeddedZip() (*zip.Reader, error) {
	bin, err := binself()
	if err != nil {
		return nil, err
	}
	fi, err := bin.Stat()
	if err != nil {
		bin.Close()
		return nil, err
	}

	n := int64(65 * 1024)
	size := fi.Size()
	if size < n {
		n = size
	}
	buf := make([]byte, n)
	_, err = io.ReadAtLeast(io.NewSectionReader(bin, size-n, n), buf, len(buf))
	if err != nil {
		bin.Close()
		return nil, err
	}
	o := int64(findSignatureInBlock(buf))
	if o < 0 {
		bin.Close()
		return nil, errors.New("could not locate zip file, no end-of-central-directory signature found")
	}
	cdirsize := int64(binary.LittleEndian.Uint32(buf[o+12:]))
	cdiroff := int64(binary.LittleEndian.Uint32(buf[o+16:]))
	zipsize := cdiroff + cdirsize + (int64(len(buf)) - o)

	rr := io.NewSectionReader(bin, size-zipsize, zipsize)
	r, err := zip.NewReader(rr, zipsize)
	if err != nil {
		bin.Close()
		return nil, err
	}

	return r, err
}
