package main

import (
	"fmt"

	sp "github.com/mnkd/slackposter"
)

type MessageBuilder struct {
	GitHubOwner     string
	GitHubRepo      string
	GitHubMilestone string
}

func filterIssues(issues []Issue, state string) []Issue {
	var filterd []Issue
	for _, issue := range issues {
		if issue.State == state {
			filterd = append(filterd, issue)
		}
	}
	return filterd
}

func (builder MessageBuilder) Summary() string {
	repo := fmt.Sprintf("%s/%s/milestone/%s", builder.GitHubOwner, builder.GitHubRepo, builder.GitHubMilestone)
	url := "https://github.com/" + repo
	link := fmt.Sprintf("<%s|%s>", url, repo)
	return fmt.Sprintf("%s\n", link)
}

func (builder MessageBuilder) buildMessage(issues []Issue) string {
	openIssues := filterIssues(issues, "open")
	closedIssues := filterIssues(issues, "closed")

	message := ""
	message = message + fmt.Sprintf("*OPEN (%v)*\n", len(openIssues))
	for _, issue := range openIssues {
		message = message + fmt.Sprintf("<%v|%v - %v> (%v)\n", issue.HTMLURL, issue.Number, issue.Title, issue.AssigneeName())
	}

	message = message + fmt.Sprintf("*CLOSED (%v)*\n", len(closedIssues))
	for _, issue := range closedIssues {
		message = message + fmt.Sprintf("<%v|%v - %v> (%v)\n", issue.HTMLURL, issue.Number, issue.Title, issue.AssigneeName())
	}

	return message
}

func (builder MessageBuilder) buildAttachment(issues []Issue, pretext string) sp.Attachment {
	message := builder.buildMessage(issues)

	return sp.Attachment{
		Pretext:  pretext,
		Fallback: message,
		Text:     message,
		Color:    "good",
		MrkdwnIn: []string{"pretext", "text", "fallback"},
	}
}

func (builder MessageBuilder) BuildAttachments(issues []Issue) []sp.Attachment {
	// Divide issues
	var issueItems []Issue
	var pullItems []Issue

	for _, issue := range issues {
		if len(issue.PullRequest.URL) > 0 {
			pullItems = append(pullItems, issue)
		} else {
			issueItems = append(issueItems, issue)
		}
	}

	return []sp.Attachment{
		builder.buildAttachment(issueItems, "*ISSUE*"),
		builder.buildAttachment(pullItems, "*PULL REQUEST*"),
	}
}

func NewMessageBuilder(gh GitHubAPI, milestone string) MessageBuilder {
	return MessageBuilder{
		GitHubOwner:     gh.Owner,
		GitHubRepo:      gh.Repo,
		GitHubMilestone: milestone,
	}
}
