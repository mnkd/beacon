package main

import (
	"fmt"
	"os"

	sp "github.com/mnkd/slackposter"
)

type App struct {
	Config    Config
	Slack     sp.SlackPoster
	Milestone string
	GitHubAPI GitHubAPI
}

func (app App) Run() int {
	// Get milestone issues from GitHub
	var issues []Issue
	issues, err := app.GitHubAPI.GetMilestoneIssues(app.Milestone)
	if err != nil {
		return ExitCodeError
	}

	builder := NewMessageBuilder(app.GitHubAPI, app.Milestone)
	summary := builder.Summary()
	attachments := builder.BuildAttachments(issues)

	fmt.Println(summary)

	payload := app.Slack.NewPayload()
	payload.Channel = app.Slack.Channel
	payload.Username = app.Slack.Username
	payload.IconEmoji = app.Slack.IconEmoji
	payload.Mrkdwn = true
	payload.Text = summary
	payload.Attachments = attachments

	err = app.Slack.PostPayload(payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, "App: <error> send a payload to slack:", err)
		return ExitCodeError
	}

	return ExitCodeOK
}

func NewApp(config Config, printList bool, owner string, repo string, milestone string) (App, error) {
	var app = App{}
	var err error
	app.Config = config
	app.Milestone = milestone
	app.GitHubAPI = NewGitHubAPI(config, owner, repo)
	app.Slack = sp.NewSlackPoster(config.SlackWebhook)
	return app, err
}
