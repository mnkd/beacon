package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Issue struct {
	Assignee struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
		URL   string `json:"url"`
	} `json:"assignee"`
	Body    string `json:"body"`
	ID      int64  `json:"id"`
	Number  int64  `json:"number"`
	State   string `json:"state"`
	Title   string `json:"title"`
	HTMLURL string `json:"html_url"`
	User    struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
		URL   string `json:"url"`
	} `json:"user"`
	PullRequest struct {
		URL string `json:"url"`
	} `json:"pull_request"`
}

func (issue Issue) AssigneeName() string {
	if len(issue.Assignee.Login) == 0 {
		return "?"
	}
	return issue.Assignee.Login
}

func (gh GitHubAPI) GetMilestoneIssues(milestone string) ([]Issue, error) {
	var issues []Issue

	// Prepare HTTP Request
	url := gh.BaseURL() + "/issues" + "?access_token=" + gh.AccessToken + "&milestone=" + milestone + "&state=all&sort=created-asc"

	req, err := http.NewRequest("GET", url, nil)

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		fmt.Fprintln(os.Stderr, "[GitHubAPI] Issue Milestone: <error> parse http request form:", parseFormErr)
		return issues, parseFormErr
	}

	// Fetch Request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[GitHubAPI] Issue Milestone: <error> fetch issues:", err)
		return issues, err
	}

	// Read Response Body
	resBody, _ := ioutil.ReadAll(res.Body)

	// Decode JSON
	if err := json.Unmarshal(resBody, &issues); err != nil {
		fmt.Fprintln(os.Stderr, "[GitHubAPI] Issue Milestone: <error> json unmarshal:", err)
		return issues, err
	}

	return issues, nil
}
