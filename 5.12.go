package main

import (
	"fmt"
	"os"
	"golang.org/x/net/html"
	"strings"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "html parse error: %s", err)
	}


	var depth int
	startElement := func(n *html.Node) {
		if n.Type == html.ElementNode {
			fmt.Printf("%*s<%s", depth*2, "", n.Data)
			for _, attr := range n.Attr {
				fmt.Printf(" %s=%q", attr.Key, attr.Val)
			}
			if n.FirstChild == nil {
				fmt.Println("/>")	
			} else{
				fmt.Println(">")
			}
		} else if n.Type == html.TextNode{
			text := strings.Trim(n.Data, "\n\t\r ")
			if text != "" {
				fmt.Printf("%*s%s\n", depth*2, "", text)
			}
		} else if n.Type == html.CommentNode {
			fmt.Printf("%*s<!--%s-->\n", depth*2, "", n.Data)
		} else if n.Type == html.DoctypeNode{
			fmt.Printf("<!DOCTYPE %s>\n", n.Data)
		}
		depth++
	}
	endElement := func(n *html.Node) {
		depth--
		if n.Type == html.ElementNode && n.FirstChild != nil{
			fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
		}
	}
	foreachNode(doc, startElement, endElement)
}

func foreachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		foreachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
