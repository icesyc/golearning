package main

import (
	"fmt"
	"os"
	"bytes"
	"strings"
)
func main() {
	for _, arg := range os.Args[1:] {
		fmt.Printf("%s\n", comma(arg))
	}	
}

func comma(s string) string {
	var suffix, prefix string
	dot := strings.Index(s, ".")
	if dot > 0 {
		s, suffix = s[:dot], s[dot:]
	}
	if strings.HasPrefix(s, "+") || strings.HasPrefix(s, "-") {
		s, prefix = s[1:], s[0:1]
	}
	n := len(s)
	pre := n % 3
	//先写入前三个字节，如果是0会导致前面多一个,
	if pre == 0 {
		pre = 3
	}
	var buf bytes.Buffer
	buf.WriteString(prefix)
	buf.WriteString(s[:pre])
	for i := pre; i < n; i += 3{
		buf.WriteByte(',')
		buf.WriteString(s[i:i+3]);
	}
	buf.WriteString(suffix)
	return buf.String()
}