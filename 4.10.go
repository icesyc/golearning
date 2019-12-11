package main

import (
	"time"
	"encoding/json"
	"net/http"
	"net/url"
	"fmt"
	"strings"
	"os"
)

const IssueURL = "https://api.github.com/search/issues"

type IssueSearchResult struct {
	TotalCount int `json:"total_count"`
	Items []*Issue
}

type Issue struct {
	Number int
	HTMLURL string `json:"html_url"`
	Title string
	State string
	User *User
	CreatedAt time.Time `json:"created_at"`
	Body string
}

type User struct {
	Login string
	HTMLURL string `json:"html_url"`
}

func main() {
	result, err := SearchIssue(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "searchIssue: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)

	timedIssues := make(map[string][]*Issue)
	now := time.Now()
	for _, item := range result.Items {
		duration := now.Sub(item.CreatedAt).Seconds()
		var category string
		if  duration < 30 * 86400 {
			category = "month"
		}else if duration < 364 * 86400 {
			category = "year"
		}else {
			category = "overyear"
		}
		timedIssues[category] = append(timedIssues[category], item)
	}
	for category, items := range timedIssues {
		fmt.Printf("\n%s\n", category)
		for _, item := range items {
			fmt.Printf("#%-5d\t%-20s\t%.55s\t%s\n", item.Number, item.User.Login, item.Title, item.CreatedAt.Format("2006-01-02"))
		}
	}
}


func SearchIssue(terms []string) (*IssueSearchResult, error) {
	query := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssueURL + "?q=" + query)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result IssueSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}