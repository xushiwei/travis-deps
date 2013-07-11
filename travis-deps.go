package main

import (
	"os"
	"fmt"
	"strings"
	"encoding/json"
	"github.com/qiniu/log"
	"github.com/qiniu/travis-deps/github"
)

type Config struct {
	Token string `json:"token"`
	Deps []string `json:"deps"`
}

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, `
Usage: traivs-deps <TravisDepsConf>

TravisDepsConf is json format file. Here is an example:

{
	"token": "<Personal API Access Token>",
	"deps": ["qiniu/errors", "qiniu/log", "qiniu/rpc"],
	"debug_level": 1
}
`)
		return
	}

	// load conf

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Warn("os.Open failed:", err)
		return
	}
	defer f.Close()

	var conf Config
	err = json.NewDecoder(f).Decode(&conf)
	if err != nil {
		log.Warn("load conf failed:", err)
		return
	}

	// token

	client := github.NewGitHubClient(conf.Token, "qiniu")

	// download dep-repos

	for _, repo := range conf.Deps {
		parts := strings.SplitN(repo, "/", 2)
		log.Info("repo:", repo, parts)
		if len(parts) != 2 {
			log.Warn("invalid repo:", repo)
			continue
		}
		keys, err := client.GetRepoKeys(map[string]string{
			"owner": parts[0],
			"repo": parts[1],
		})
		if err != nil {
			log.Warn("GetRepoKeys failed:", err)
			return
		}
		log.Info("GetRepoKeys:", keys)
	}
}

