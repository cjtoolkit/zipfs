[![GoDoc](https://godoc.org/github.com/cjtoolkit/zipfs?status.svg)](https://godoc.org/github.com/cjtoolkit/zipfs)

# ZipFS for net/http

## Installation

```sh
$ go get github.com/cjtoolkit/zipfs
```

## Example

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cjtoolkit/zipfs"
)

func main() {
	// Init ZipFS, if the zip file does not exist it will fallback to embedded zip file within application, if that
	// also does not exist, it will panic.
	fs := zipfs.InitZipFs("asset.zip")

	fmt.Println("Running ZipFS via http server on port 8080.")
	log.Print(http.ListenAndServe(":8080", http.FileServer(fs)))
}
```

To embed zip to the application, first create a zip file with your preferred achiever with compression disabled.
Than append to the content of the zip to the compiled application.  You can use this simple shell command to append.

```sh
$ cat asset.zip >> application
```

## Credit

This project is based on the work of the following:

* Rémy Oudompheng
    * [ZipFS proof of concept](https://github.com/remyoudompheng/go-misc/blob/master/zipfs/zipfs.go)
* Mechiel Lukkien
    * [Embedding/Appending Zip File into Applications](https://godoc.org/bitbucket.org/mjl/httpasset)
* Derek Parker
    * [Trie Data Structure](https://github.com/derekparker/trie)