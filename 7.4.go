package main

import (
	"fmt"
	"io"
	"bytes"
	"golang.org/x/net/html"
)

func main() {
	str := "<html><body><p>test</p></body></html>"
	doc, _ := html.Parse(NewReader(str))
	fmt.Printf("%v\n", htmlString(doc))
}

func ForEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ForEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

func htmlString(doc *html.Node) string {
	var buf bytes.Buffer
	deep := 0
	pre := func(n *html.Node) {
		if n.Type == html.ElementNode {
			tag := fmt.Sprintf("%*s<%s>\n", deep*4, "", n.Data)
			buf.WriteString(tag)
			deep++
		}
	}
	post := func(n *html.Node) {
		if n.Type == html.ElementNode {
			deep--
			closeTag := fmt.Sprintf("%*s</%s>\n", deep*4, "", n.Data)
			buf.WriteString(closeTag)
		}
	}
	ForEachNode(doc, pre, post);
	return buf.String()
}

type Reader struct {
	data string
	position int
}

func NewReader(str string) io.Reader {
	r := new(Reader)
	r.data = str
	return r
}

func (r *Reader) Read(p []byte) (n int, err error) {
	max := len(p)
	var i int
	p = p[:0]
	for i = 0; i < max; i++ {
		if r.position >= len(r.data) {
			break
		}
		p = append(p, r.data[r.position])
		r.position++
	}
	if i < max {
		err = io.EOF
	}
	return i, err
}