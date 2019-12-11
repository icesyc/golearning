package main

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"html/template"
	"log"
)

type SearchIssueResult struct {
	TotalCount int `json:"total_count"`
	Items []*Issue
}

type Issue struct {
	Number int
	HTMLURL string `json:"html_url"`
	Title string
	CreatedAt time.Time `json:"created_at"`
	User *User
	State string
}
type User struct {
	Login string
	HTMLURL string `json:"html_url"`
}

func main() {
	http.HandleFunc("/", github)	
	http.ListenAndServe("localhost:8080", nil)
}

func github(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	result, err := requestIssue(query)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return 
	}
	tpl := `
		<h1>{{.TotalCount}}</h1>
		<table>
			<tr>
				<td>Number</td>
				<td>Title</td>
				<td>State</td>
				<td>User</td>
				<td>date</td>
			</tr>
			{{range .Items}}
			<tr>
				<td>{{.Number}}</td>
				<td><a href="{{.HTMLURL}}">{{.Title}}</a></td>
				<td>{{.State}}</td>
				<td><a href="{{.User.HTMLURL}}">{{.User.Login}}</td>
				<td>{{.CreatedAt | format "2006-01-02"}}</td>
			</tr>
			{{end}}
		</table>`
	html := template.Must(template.New("html").Funcs(template.FuncMap{"format": format}).Parse(tpl))
	if err = html.Execute(w, result); err != nil {
		log.Fatal(err)
	}
}

func format(layout string, t time.Time) string {
	return t.Format(layout)
}

func requestIssue(query string) (*SearchIssueResult, error) {
	const api = "https://api.github.com/search/issues"
	url := api + "?q=" + query
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("request error: %s" ,resp.Status)
	}
	var result SearchIssueResult
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}