package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gdbu/filecacher"
)

func main() {
	var (
		err error
	)

	f := filecacher.New("./")

	for {
		if err = f.Read("./hello.txt", printReader); err != nil {
			log.Println("Error reading", err)
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
