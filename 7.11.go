package main

import (
	"net/http"
	"fmt"
	"strconv"
	"html/template"
)

func main() {
	db := database{"shoes": 50, "socks": 5}
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/read", db.read)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)
	http.ListenAndServe("localhost:8080", nil)
}

type dollars float32
type database map[string]dollars

func (d dollars) String() string {
	return fmt.Sprintf("%.2f", d)
}

func (db database) list(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html;charset=utf-8")
	html := `
	<form method="get" action="/create">
		<input type=text name="item"/>
		<input type=text name="price"/>
		<input type="submit" value="添加"/>
	</form>
	`
	fmt.Fprintf(w, html)

	tpl := `
	{{range $item, $price := .}}
	<form method="get" action="/update">
		{{$item}}
		<input type=text name="price" value="{{$price}}"/>
		<input type="hidden" name="item" value="{{$item}}"/>
		<input type="submit" value="修改"/>
		<a href="/delete?item={{$item}}">删除</a>
	</form>
	{{end}}`
	temp := template.Must(template.New("html").Parse(tpl))
	temp.Execute(w, db)
}

func (db database) read(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	item := r.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %s\n", item)
		return 
	}
	fmt.Fprintf(w, "%s\n", price)
	fmt.Fprintf(w, "<p><a href=\"/list\">返回</a></p>")
}

func (db database) create(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price := r.URL.Query().Get("price")
	newPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		fmt.Fprintf(w, "error price: %s\n", err)
		return
	}
	if _, ok := db[item]; ok {
		fmt.Fprintf(w, "item exists: %s\n", item)
		return
	}
	db[item] = dollars(newPrice)
	w.Header().Set("Location", "/list")
	w.WriteHeader(http.StatusFound)
}

func (db database) update(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price := r.URL.Query().Get("price")
	newPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		fmt.Fprintf(w, "error price: %s\n", price)
		return
	}
	if _, ok := db[item]; !ok {
		fmt.Fprintf(w, "no such item: %s\n", item)
		return
	}
	db[item] = dollars(newPrice)
	w.Header().Set("Location", "/list")
	w.WriteHeader(http.StatusFound)
}
func (db database) delete(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	delete(db, item)
	w.Header().Set("Location", "/list")
	w.WriteHeader(http.StatusFound)
}
