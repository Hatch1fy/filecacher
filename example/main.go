package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Hatch1fy/filecacher"
)

func main() {
	var (
		f   *filecacher.File
		err error
	)

	if f, err = filecacher.NewFile("./hello.txt"); err != nil {
		log.Fatal(err)
	}

	for {
		if err = f.Read(printReader); err != nil {
			break
		}

		time.Sleep(time.Second * 3)
	}
}

func printReader(r io.Reader) (err error) {
	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, r); err != nil {
		return
	}

	fmt.Println("Output", buf.String())
	return
}
