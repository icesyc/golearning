package main

import (
	"fmt"
	"bufio"
	"strings"
)

func main() {
	var c Counter
	s := `
	this is a test
	haha I love you
	do you love me ?`
	fmt.Fprintf(&c, s)
	fmt.Printf("%s\n", c)
}

type Counter struct {
	lines int
	words int
}
func (c Counter) String() string{
	return fmt.Sprintf("{lines=%d, words=%d}", c.lines, c.words)
}

func (c *Counter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(p)))
	for scanner.Scan() {
		c.lines++
		wordScanner := bufio.NewScanner(strings.NewReader(scanner.Text()))
		wordScanner.Split(bufio.ScanWords)
		for wordScanner.Scan() {
			c.words++
		}
	}
	return c.lines, nil
}