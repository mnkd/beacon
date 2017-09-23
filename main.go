package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	ExitCodeOK int = iota
	ExitCodeError
	ExitCodeFileError
)

var (
	version  string
	revision string
)

func usage() {
	str := `
Usage:

 beacon -o owner -r repo -m milestone [-d]

Examples:

 Print owner/awesome-app/milesstone/15
 $ beacon -o owner -r awesome-app -m 15
`
	fmt.Fprintln(os.Stderr, str)
}

var app App

func init() {
	var configPath, milestone, owner, repo string
	var list, ver, dryRun bool

	flag.StringVar(&configPath, "c", "", "/path/to/config.json. (default: $HOME/.config/beacon/config.json)")
	flag.StringVar(&owner, "o", "", "owner (e.g. github)")
	flag.StringVar(&repo, "r", "", "repo (e.g. hub)")
	flag.StringVar(&milestone, "m", "", "milestone number")
	flag.BoolVar(&dryRun, "d", false, "A dry run will not send any message to Slack. (defualt: false)")
	flag.BoolVar(&ver, "v", false, "Print version.")
	flag.Parse()

	if ver {
		fmt.Fprintln(os.Stdout, "Version:", version)
		fmt.Fprintln(os.Stdout, "Revision:", revision)
		os.Exit(ExitCodeOK)
	}

	if len(owner) == 0 || len(repo) == 0 {
		flag.Usage()
		os.Exit(ExitCodeOK)
	}

	// Prepare config
	config, err := NewConfig(configPath, dryRun)
	if err != nil {
		os.Exit(ExitCodeError)
	}

	// Prepare app
	app, err = NewApp(config, list, owner, repo, milestone)
	if err != nil {
		os.Exit(ExitCodeError)
	}
}

func main() {
	os.Exit(app.Run())
}
