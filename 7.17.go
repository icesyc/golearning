package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"encoding/xml"
	"text/scanner"

)

type Selector struct {
	Type int
	Value string
}

const (
	Id = iota + 1
	Class
	Tag
)

func main(){
	if len(os.Args) < 2 {
		fmt.Printf("please input the selector\n")
		os.Exit(1)
	}
	dec := xml.NewDecoder(os.Stdin)
	var stack []xml.StartElement
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
			stack = append(stack, tok)
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		case xml.CharData:
			if matchCss(stack, os.Args[1:]) {
				fmt.Printf("%s: %s\n", stackString(stack), tok)
			}
		}
	}
}

func stackString(stack []xml.StartElement) string {
	var res []string
	for _, el := range stack {
		html := "[" + el.Name.Local
		for _, attr := range el.Attr {
			html += fmt.Sprintf(" %s=%q", attr.Name.Local, attr.Value)
		}
		html += "]"
		res = append(res, html)
	}
	return strings.Join(res, " > ")
}

func parseSelector(selector string) []Selector {
	var sc scanner.Scanner
	sc.Init(strings.NewReader(selector))
	sc.Mode = scanner.ScanIdents
	token := sc.Scan()
	var selList []Selector
	for {
		if token == scanner.EOF {
			break
		}
		switch token {
		case scanner.Ident:
			selList = append(selList, Selector{Tag, sc.TokenText()})
		case '.', '#':
			selectorType := Id
			if token == '.' {
				selectorType = Class
			}
			token := sc.Scan()
			if token != scanner.Ident {
				panic(fmt.Sprintf("selector parse error, unexpected %s", sc.TokenText()))
			}
			value := sc.TokenText()	
			selList = append(selList, Selector{selectorType, value})
		}
		token = sc.Scan()
	}
	return selList
}
func matchCss(stack []xml.StartElement, selectors []string) bool {
	for len(selectors) <= len(stack) {
		if len(selectors) == 0 {
			return true
		}
		if matchSelector(stack[0], selectors[0]) {
			selectors = selectors[1:]
		}
		stack = stack[1:]
	}
	return false
}
func matchSelector(el xml.StartElement, selector string) bool {
	selectorList := parseSelector(selector)
	for _, sel := range selectorList {
		switch sel.Type {
		case Tag:
			if el.Name.Local != sel.Value {
				return false
			}
		case Id, Class:
			selectorType := "id"
			if sel.Type == Class {
				selectorType = "class"
			}
			if !hasAttr(el, selectorType, sel.Value) {
				return false
			}
		}
	}
	return true
}

func hasAttr(e xml.StartElement, name, value string) bool{
	for _, attr := range e.Attr {
		if attr.Name.Local != name {
			continue
		}
		if name == "id" && attr.Value == value {
			return true
		}
		if name == "class" {
			for _, cls := range strings.Split(attr.Value, " ") {
				if cls == value {
					return true
				}
			}
		}

	}
	return false
}