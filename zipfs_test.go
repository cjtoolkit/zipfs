package zipfs

import (
	"testing"
)

func TestZipFS_Open(t *testing.T) {
	t.Run("Without Compression", func(t *testing.T) {
		fs := InitZipFs("testdata/uncompressed.zip")

		{
			file, err := fs.Open("/")
			if err != nil {
				t.Fail()
			}
			if _, ok := file.(*zipRoot); !ok {
				t.Fail()
			}
			if dirs, _ := file.Readdir(-1); len(dirs) != 2 {
				t.Fail()
			}
		}

		{
			file, err := fs.Open("/text1.txt")
			if err != nil {
				t.Fail()
			}
			if _, ok := file.(*uncompressedFile); !ok {
				t.Fail()
			}
		}

		{
			file, err := fs.Open("/dirA")
			if err != nil {
				t.Fail()
			}
			if _, ok := file.(*zipDir); !ok {
				t.Fail()
			}
			if dirs, _ := file.Readdir(-1); len(dirs) != 3 {
				t.Fail()
			}
		}

		{
			file, err := fs.Open("/dirA/dirB")
			if err != nil {
				t.Fail()
			}
			if _, ok := file.(*zipDir); !ok {
				t.Fail()
			}
			if dirs, _ := file.Readdir(-1); len(dirs) != 2 {
				t.Fail()
			}
		}

		{
			file, err := fs.Open("/dirA/dirB/text3.txt")
			if err != nil {
				t.Fail()
			}
			if _, ok := file.(*uncompressedFile); !ok {
				t.Fail()
			}
		}

		{
			file, err := fs.Open("/dirA/dirB/text4.txt")
			if err != nil {
				t.Fail()
			}
			if _, ok := file.(*uncompressedFile); !ok {
				t.Fail()
			}
		}

		{
			file, err := fs.Open("/dirA/dirC")
			if err != nil {
				t.Fail()
			}
			if _, ok := file.(*zipDir); !ok {
				t.Fail()
			}
			if dirs, _ := file.Readdir(-1); len(dirs) != 2 {
				t.Fail()
			}
		}

		{
			file, err := fs.Open("/dirA/dirC/text5.txt")
			if err != nil {
				t.Fail()
			}
			if _, ok := file.(*uncompressedFile); !ok {
				t.Fail()
			}
		}

		{
			file, err := fs.Open("/dirA/dirC/text6.txt")
			if err != nil {
				t.Fail()
			}
			if _, ok := file.(*uncompressedFile); !ok {
				t.Fail()
			}
		}

		{
			_, err := fs.Open("/dirA/dirC/text7.txt")
			if err == nil {
				t.Fail()
			}
		}

		{
			_, err := fs.Open("")
			if err == nil {
				t.Fail()
			}
		}
	})

	t.Run("With Compression", func(t *testing.T) {
		fs := InitZipFs("testdata/compressed.zip")

		{
			file, err := fs.Open("/")
			if err != nil {
				t.Fail()
			}
			if _, ok := file.(*zipRoot); !ok {
				t.Fail()
			}
			if dirs, _ := file.Readdir(-1); len(dirs) != 2 {
				t.Fail()
			}
		}

		{
			file, err := fs.Open("/text1.txt")
			if err != nil {
				t.Fail()
			}
			if _, ok := file.(*compressedFile); !ok {
				t.Fail()
			}
		}

		{
			file, err := fs.Open("/dirA")
			if err != nil {
				t.Fail()
			}
			if _, ok := file.(*zipDir); !ok {
				t.Fail()
			}
			if dirs, _ := file.Readdir(-1); len(dirs) != 3 {
				t.Fail()
			}
		}

		{
			file, err := fs.Open("/dirA/dirB")
			if err != nil {
				t.Fail()
			}
			if _, ok := file.(*zipDir); !ok {
				t.Fail()
			}
			if dirs, _ := file.Readdir(-1); len(dirs) != 2 {
				t.Fail()
			}
		}

		{
			file, err := fs.Open("/dirA/dirB/text3.txt")
			if err != nil {
				t.Fail()
			}
			if _, ok := file.(*compressedFile); !ok {
				t.Fail()
			}
		}

		{
			file, err := fs.Open("/dirA/dirB/text4.txt")
			if err != nil {
				t.Fail()
			}
			if _, ok := file.(*compressedFile); !ok {
				t.Fail()
			}
		}

		{
			file, err := fs.Open("/dirA/dirC")
			if err != nil {
				t.Fail()
			}
			if _, ok := file.(*zipDir); !ok {
				t.Fail()
			}
			if dirs, _ := file.Readdir(-1); len(dirs) != 2 {
				t.Fail()
			}
		}

		{
			file, err := fs.Open("/dirA/dirC/text5.txt")
			if err != nil {
				t.Fail()
			}
			if _, ok := file.(*compressedFile); !ok {
				t.Fail()
			}
		}

		{
			file, err := fs.Open("/dirA/dirC/text6.txt")
			if err != nil {
				t.Fail()
			}
			if _, ok := file.(*compressedFile); !ok {
				t.Fail()
			}
		}

		{
			_, err := fs.Open("/dirA/dirC/text7.txt")
			if err == nil {
				t.Fail()
			}
		}

		{
			_, err := fs.Open("")
			if err == nil {
				t.Fail()
			}
		}
	})
}

func BenchmarkZipFS_Open(b *testing.B) {
	b.Run("Without Compression", func(b *testing.B) {
		fs := InitZipFs("testdata/uncompressed.zip")
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			fs.Open("/dirA/dirC/text6.txt")
		}
	})

	b.Run("With Compression", func(b *testing.B) {
		fs := InitZipFs("testdata/compressed.zip")
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			fs.Open("/dirA/dirC/text6.txt")
		}
	})
}
