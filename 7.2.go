package main

import (
	"fmt"
	"os"
	"io"
)

func main() {
	w, n := CountingWriter(os.Stdout)
	fmt.Fprintf(w, "%s", "this is a test\n")
	fmt.Printf("total write %d bytes\n", *n)
	fmt.Fprintf(w, "%s", "new data\n")
	fmt.Printf("total write %d bytes\n", *n)
}
type CountWriter struct{
	writer io.Writer
	n int64
}

func (cw *CountWriter) Write(p []byte) (int, error) {
	n, err := cw.writer.Write(p)
	cw.n += int64(n)
	return n, err
}
func CountingWriter(w io.Writer) (io.Writer, *int64) {
	cw := CountWriter{writer: w}
	return &cw, &cw.n
}