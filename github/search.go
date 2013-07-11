package github

//
// GitHub API v3 Section - Search
// Allows you to search for emails, users, repos and issues
//
//	## Search API
//		-  Search issues
//		-  Search repositories
//		-  Search users
//		-  Email search

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"strings"
)

type SearchIssue struct {
	Title      string   `json:"title"`
	User       string   `json:"user"`
	Body       Nstring  `json:"body"`
	Position   int      `json:"position"`
	Number     int      `json:"number"`
	Comments   int      `json:"comments"`
	Votes      int      `json:"votes"`
	Labels     []string `json:"labels"`
	State      Nstring  `json:"state"`
	GravatarId string   `json:"gravatar_id"`
	HtmlUrl    string   `json:"html_url"`
	UpdatedAt  Nstring  `json:"updated_at"`
	CreatedAt  string   `json:"created_at"`
}

type SearchUser struct {
	Name           string  `json:"name"`
	Id             string  `json:"id"`
	FullName       string  `json:"fullname"`
	Language       Nstring `json:"language"`
	Username       string  `json:"username"`
	Location       Nstring `json:"location"`
	Type           string  `json:"type"`
	Repos          int     `json:"repos"`
	Followers      int     `json:"followers"`
	PublicRepos    int     `json:"public_repo_count"`
	FollowersCount int     `json:"followers_count"`
	Score          float64 `json:"name"`
	GravatarId     string  `json:"gravatar_id"`
	CreatedAt      Nstring `json:"created_at"`
	Created        Nstring `json:"created"`
	Pushed         Nstring `json:"pushed"`
	Record         Nstring `json:"record"`
}

type SearchRepo struct {
	Name        string  `json:"name,omitempty"`
	Owner       string  `json:"owner,omitempty"`
	Type        string  `json:"type,omitempty"`
	Username    string  `json:"username,omitempty"`
	Url         string  `json:"url,omitempty"`
	Description Nstring `json:"description,omitempty"`
	Watchers    int     `json:"watchers"`
	Forks       int     `json:"forks"`
	Size        int     `json:"int"`
	Followers   int     `json:"followers"`
	OpenIssues  int     `json:"open_issues"`
	Language    Nstring `json:"language,omitempty"`
	Score       float64 `json:"score,omitempty"`
	UpdatedAt   Nstring `json:"updated_at,omitempty"`
	CreatedAt   string  `json:"created_at,omitempty"`
	Created     string  `json:"created,omitempty"`
	PushedAt    Nstring `json:"pushed_at,omitempty"`
	Homepage    Nstring `json:"homepage,omitempty"`
	Downloads   bool    `json:"has_downloads,omitempty"`
	Wiki        bool    `json:"has_wiki,omitempty"`
	Issues      bool    `json:"has_issues,omitempty"`
	Private     bool    `json:"private,omitempty"`
}

// 
// GitHub Doc - Search: Search Issues
// Url: https://api.github.com/legacy/issues/search/:owner/:repository/:state/:keyword?access_token=...
// Request Type: GET /legacy/issues/search/:owner/:repository/:state/:keyword
// Access Token: REQUIRED
// 

func (github *GitHubClient) SearchIssues(urlData map[string]string) ([]SearchIssue, error) {
	if ok := github.AssertMapStrings([]string{"repo", "state", "keyword"}, urlData); !ok {
		return nil, errors.New("One or more fields are missing and/or do not have content.")
	}
	if urlData["state"] != "open" && urlData["state"] != "closed" {
		return nil, errors.New("The state value in urlData is not a valid option - Only open and closed are acceptable.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/legacy/issues/search/" + urlData["owner"] + "/" + urlData["repo"] + "/" + urlData["state"] + "/" + urlData["keyword"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		issues := &map[string][]SearchIssue{}
		issuesJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(issuesJson, issues); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*issues)["issues"], nil
	}

	github.getLimits(res)
	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - Search: Search repositories
// Url: https://api.github.com/legacy/repos/search/:keyword?access_token=...
// Request Type: GET /legacy/repos/search/:keyword
// Access Token: REQUIRED
// 

func (github *GitHubClient) SearchRepos(keyword string, getData map[string]string) ([]SearchRepo, error) {
	if strings.TrimSpace(keyword) == "" {
		return nil, errors.New("The keyword does not contain any non-whitespace characters.")
	}

	apiUrl := github.createUrl("/legacy/repos/search/" + keyword + "?" + github.UrlDataConvert(getData))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		repos := &map[string][]SearchRepo{}
		reposJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(reposJson, repos); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*repos)["repositories"], nil
	}

	github.getLimits(res)
	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - Search: Search users - Find users by keyword.
// Url: https://api.github.com/legacy/repos/search/:keyword?access_token=...
// Request Type: GET /legacy/repos/search/:keyword
// Access Token: REQUIRED
// 

func (github *GitHubClient) SearchUsers(keyword, startPage string) ([]SearchUser, error) {
	if strings.TrimSpace(keyword) == "" {
		return nil, errors.New("The keyword does not contain any non-whitespace characters.")
	}

	apiUrl := ""
	if len(strings.TrimSpace(startPage)) > 0 {
		apiUrl = github.createUrl("/legacy/user/search/" + keyword + "?start_page=" + url.QueryEscape(startPage))
	} else {
		apiUrl = github.createUrl("/legacy/user/search/" + keyword)
	}

	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		users := &map[string][]SearchUser{}
		usersJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(usersJson, users); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*users)["users"], nil
	}

	github.getLimits(res)
	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - Search: Email search
// Url: https://api.github.com/legacy/user/email/:email?access_token=...
// Request Type: GET /legacy/user/email/:email
// Access Token: REQUIRED
// 

func (github *GitHubClient) SearchEmail(email string) (*SearchUser, error) {
	if strings.TrimSpace(email) == "" {
		return nil, errors.New("The email does not contain any non-whitespace characters.")
	}

	apiUrl := github.createUrl("/legacy/user/email/" + email)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		user := &map[string]*SearchUser{}
		userJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(userJson, user); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*user)["user"], nil
	}

	github.getLimits(res)
	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}
