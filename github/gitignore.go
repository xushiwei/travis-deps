package github

// GitIgnore of the GitHub API
// Includes the comments since I think they are much more useful in Gists, but rare for commits.
//
//	##  Gitignore Templates API
//		-  Listing available templates
//		-  Get a single template

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
)

type GitIgnore struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

// 
// GitHub Doc - GitIgnore: Listing available templates
// Url: https://api.github.com/gitignore/templates?access_token=...
// Request Type: GET /gitignore/templates
// Access Token: PUBLIC
// 
func (github *GitHubClient) ListTemplates() ([]string, error) {
	apiUrl := github.createUrl("/gitignore/templates")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		templates := &[]string{}
		templatesJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(templatesJson, templates); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*templates), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - GitIgnore: Get a single template
// Url: https://api.github.com/gitignore/templates/:template?access_token=...
// Request Type: GET /gitignore/templates/:template
// Access Token: PUBLIC
// 
func (github *GitHubClient) GetTemplate(template string) (*GitIgnore, error) {
	template = strings.TrimSpace(template)
	if len(template) == 0 {
		return nil, errors.New("The template value does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/gitignore/templates/" + template)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		template := &GitIgnore{}
		templateJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(templateJson, template); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return template, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}
