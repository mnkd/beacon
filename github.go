package main

type GitHubAPI struct {
	AccessToken string
	Owner       string
	Repo        string
}

func NewGitHubAPI(config Config, owner string, repo string) GitHubAPI {
	var gh = GitHubAPI{}
	gh.AccessToken = config.GitHub.AccessToken
	gh.Owner = owner
	gh.Repo = repo
	return gh
}

func (gh GitHubAPI) BaseURL() string {
	return "https://api.github.com/repos/" + gh.Owner + "/" + gh.Repo
}
