package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

var (
	ErrInvalidJson = errors.New("ErrInvalidJson")
)

type Config struct {
	App struct {
		UseHolidayJP bool `json:"use_holiday_jp"`
		UseDayOff    bool `json:"use_dayoff"`
	} `json:"app"`
	GitHub struct {
		AccessToken string `json:"access_token"`
	} `json:"github"`
	SlackWebhook struct {
		Channel    string `json:"channel"`
		IconEmoji  string `json:"icon_emoji"`
		Username   string `json:"username"`
		WebhookURL string `json:"webhook_url"`
	} `json:"slack_webhook"`

	DryRun bool
}

func (config *Config) validate() error {
	if len(config.GitHub.AccessToken) == 0 {
		fmt.Fprintln(os.Stderr, "Invalid config.json. You should set a github access_token.")
		return ErrInvalidJson
	}
	return nil
}

func NewConfig(path string, dryRun bool) (Config, error) {
	var config Config
	config.DryRun = dryRun

	usr, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Config: <error> get current user:", err)
		return config, err
	}

	if len(path) == 0 {
		path = filepath.Join(usr.HomeDir, "/.config/beacon/config.json")
	} else {
		p, absErr := filepath.Abs(path)
		if absErr != nil {
			fmt.Fprintln(os.Stderr, "Config: <error> get absolute representation of path:", absErr, path)
			return config, absErr
		}
		path = p
	}

	str, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Config: <error> read config.json:", err)
		return config, err
	}

	if err = json.Unmarshal(str, &config); err != nil {
		fmt.Fprintln(os.Stderr, "Config: <error> json unmarshal: config.json:", err)
		return config, err
	}

	if err = config.validate(); err != nil {
		return config, err
	}

	return config, nil
}
