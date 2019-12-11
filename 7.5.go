package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	str := "this is a test, I love you!"	
	var p [100]byte
	reader := LimitReader(strings.NewReader(str), 10)
	reader.Read(p[:])
	fmt.Printf("%s\n", p)
}

type StrLimitReader struct {
	reader io.Reader
	max int
}

func (lr StrLimitReader) Read(p []byte) (n int, err error) {
	if lr.max >= len(p) {
		n = len(p)
	} else {
		n = lr.max
	}
	n, err = lr.reader.Read(p[:n])
	lr.max -= n
	return n, err
}

func LimitReader(r io.Reader, n int64) io.Reader {
	return &StrLimitReader{reader: r, max: int(n)};
}