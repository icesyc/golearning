package main

import (
	"fmt"
	"io"
	"os"
	"encoding/xml"
	"strings"
)

type Node interface{} // CharData or *Element

type CharData string

type Element struct {
	Type xml.Name
	Attr []xml.Attr
	Children []Node
}

func (e *Element) String() string{
	html := fmt.Sprintf("<%s", e.Type.Local)
	for _, attr := range e.Attr {
		html += fmt.Sprintf(" %s=%q", attr.Name.Local, attr.Value)
	}
	html += ">\n"
	for _, node := range e.Children {
		switch node := node.(type) {
		case CharData:
			html += indent(string(node))
		case *Element:
			html += indent(node.String())
		}
		html += "\n"
	}
	html += fmt.Sprintf("</%s>", e.Type.Local)
	return html
}
func indent(s string) string {
	pad := fmt.Sprintf("%*s", 4, "")
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = pad + line
	}
	return strings.Join(lines, "\n")
}

func main(){
	dec := xml.NewDecoder(os.Stdin)
	var node *Element
	var stack []*Element
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break;
		}else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			current := &Element{Type: tok.Name, Attr: tok.Attr}
			stack = append(stack, current)
			if node == nil {
				node = current
				continue
			}
			node.Children = append(node.Children, current)
			node = current
		case xml.EndElement:
			stack = stack[:len(stack)-1]
			if len(stack) > 0 {
				node = stack[len(stack)-1]
			}
		case xml.CharData:
			if node == nil || strings.TrimSpace(string(tok)) == "" {
				continue
			}
			node.Children = append(node.Children, Node(CharData(tok)))
		}
	}
	fmt.Printf("%s\n", node)
}

