package main

import (
	"net/http"
	"fmt"
	"ioutil"
)

func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

type Func func(key string) (interface{}, error)
type result struct {
	value interface{}
	err error
}
type entry struct {
	res result
	ready chan struct{}
}
type request struct {
	key string
	response chan<- result
}
type Memo struct {
	requests chan request
}
func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (memo *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) server(f Func) {
	cache := make([string]*entry)
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key)
		}
		go e.deliver()
	}
}

func (e *entry) call(f Func, key string){
	e.res.value, e.res.err = f(key)
	<-e.ready
	response <- e.res
}