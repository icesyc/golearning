package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

var cancel = make(chan struct{})

func main() {
	fmt.Printf("%v\n", mirroredQuery())
}

func mirroredQuery() string{
	resp := make(chan string)
	go func() { resp <- request("http://www.baidu.com") }()
	go func() { resp <- request("http://www.qq.com") }()
	go func() { resp <- request("http://www.facebook.com") }()
	res := <-resp
	close(cancel)
	return res
}

func request(url string) string {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Cancel = cancel
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("%v\n", err)
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("error: %v\n", resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		return ""
	}
	return string(data)
}