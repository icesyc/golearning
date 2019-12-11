package main

import (
	"fmt"
	"os"
	"net/http"
	"net/url"
	"strings"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"time"
)

var host string
//用于取消http请求的channel
var cancel = make(chan struct{})

func main() {
	var workList = make(chan []string)
	var linkQueue = make(chan string)
	var seen = make(map[string]bool)
	var n = 0

	uu, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	host = uu.Hostname()
	os.Mkdir(host, 0755)
	err = os.Chdir(host)
	if err != nil {
		log.Fatal(err)
	}

	go func(){
		workList <- os.Args[1:]
	}()
	n++

	//开启10个worker
	for i := 0; i < 10; i++ {
		go func() {
			for link := range linkQueue {
				newLinks := crawl(link)
				go func(){
					workList <- newLinks
				}()
			}
		}()
	}

	go func(){
		time.Sleep(5 * time.Second)
		close(cancel)
	}()
	for ; n > 0; n-- {
		links := <- workList
		for _, link := range links {
			if !seen[link] {
				u, _ := url.Parse(link)
				//只抓原始host的内容
				if u.Hostname() != host || strings.Contains(link, "/bbs/") {
					continue
				}
				seen[link] = true
				linkQueue <- link
				n++
			}
		}
	}
}

func ForeachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ForeachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

func Extract(url string) (string, []string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Cancel = cancel
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		err := fmt.Errorf("getting %s error: %s\n", url, resp.Status)
		return "", nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	doc, err := html.Parse(strings.NewReader(string(content)))
	resp.Body.Close()
	if err != nil {
		return "", nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link, err := resp.Request.URL.Parse(attr.Val)
					if err != nil {
						continue
					}
					links = append(links, link.String())
				}
			}
		}
	}
	ForeachNode(doc, visitNode, nil)
	return string(content), links, nil
}


func crawl(u string) []string {
	uu, err := url.Parse(u)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s error: %v\n", u, err)
		return nil
	}
	
	content, links, err := Extract(u)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return nil
	}
	var values []string
	for _, v := range uu.Query() {
		values = append(values, v...)
	}
	fname := strings.Join(values, "_") + ".html"
	fname = strings.Replace(fname, "/", "_", -1)
	saveFile(fname, content)
	return links
}

func saveFile(name string, data string) error {
	f, err := os.Create(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "-> can not create file: %s, %v ", name, err)
		return err
	}
	_, err = f.WriteString(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "-> can not write file: %s, %v ", name, err)
		return err
	}
	err = f.Close()
	return err
}