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
	message := ""
	for _, issue := range issues {
		message = message + fmt.Sprintf("<%v|%v - %v> (%v)\n", issue.HTMLURL, issue.Number, issue.Title, issue.AssigneeName())
	}
	return message
}

func (builder MessageBuilder) buildAttachment(issues []Issue, pretext string, color string) sp.Attachment {
	message := builder.buildMessage(issues)

	return sp.Attachment{
		Pretext:  pretext,
		Fallback: message,
		Text:     message,
		Color:    color,
		MrkdwnIn: []string{"pretext", "text", "fallback"},
	}
}

func (builder MessageBuilder) buildAttachments(issues []Issue, pretext string) []sp.Attachment {
	openIssues := filterIssues(issues, "open")
	closedIssues := filterIssues(issues, "closed")

	openTitle := fmt.Sprintf("*%v | OPEN (%d)*", pretext, len(openIssues))
	a1 := builder.buildAttachment(openIssues, openTitle, "danger")

	closedTitle := fmt.Sprintf("*%v | CLOSED (%d)*", pretext, len(closedIssues))
	a2 := builder.buildAttachment(closedIssues, closedTitle, "good")

	return []sp.Attachment{a1, a2}
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

	var attachments []sp.Attachment
	a1 := builder.buildAttachments(issueItems, "ISSUE")
	a2 := builder.buildAttachments(pullItems, "PULL REQUEST")
	attachments = append(attachments, a1...)
	attachments = append(attachments, a2...)
	return attachments
}

func NewMessageBuilder(gh GitHubAPI, milestone string) MessageBuilder {
	return MessageBuilder{
		GitHubOwner:     gh.Owner,
		GitHubRepo:      gh.Repo,
		GitHubMilestone: milestone,
	}
}
