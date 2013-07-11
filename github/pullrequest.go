package github

// GitHub API v3 Section - Search
// Allows you to search for emails, users, repos and issues
//
//	## Pull Request API
//		-  Link Relations
//		-  List pull requests
//		-  Get a single pull request
//		-  Create a pull request
//		-  Update a pull request
//		-  List commits on a pull request
//		-  List pull requests files
//		-  Get if a pull request has been merged
//		-  Merge a pull request (Merge Button™)
//		-  Custom media types

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type PullMerge struct {
	SHA     Nstring `json:"sha"`
	Merged  bool    `json:"merged"`
	Message Nstring `json:"message"`
}

type PullLinks struct {
	Self           map[string]string `json:"self"`
	Html           map[string]string `json:"html"`
	Comments       map[string]string `json:"comments"`
	ReviewComments map[string]string `json:"review_comments"`
}

type PullBase struct {
	Label Nstring `json:"label"`
	Ref   Nstring `json:"ref"`
	SHA   string  `json:"sha"`
	User  GitUser `json:"user"`
	Repo  Fork    `json:"repo"`
}

type PullHead struct {
	Label string  `json:"label"`
	Ref   string  `json:"ref"`
	SHA   string  `json:"sha"`
	User  GitUser `json:"user"`
	Repo  Fork    `json:"repo"`
}

type PullRequest struct {
	Url            string    `json:"url"`
	HtmlUrl        string    `json:"html_url"`
	DiffUrl        string    `json:"diff_url"`
	PatchUrl       Nstring   `json:"patch_url"`
	IssueUrl       Nstring   `json:"issue_url"`
	Number         int       `json:"number"`
	State          Nstring   `json:"state"`
	Title          string    `json:"title"`
	Body           Nstring   `json:"body"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      Nstring   `json:"updated_at"`
	ClosedAt       Nstring   `json:"closed_at"`
	MergedAt       Nstring   `json:"merged_at"`
	Head           PullHead  `json:"head"`
	Base           PullBase  `json:"base"`
	Links          PullLinks `json:"_links"`
	User           GitUser   `json:"user"`
	MergeCommitSHA Nstring   `json:"merge_commit_sha"`
	Merged         bool      `json:"merged,omitempty"`
	Mergable       bool      `json:"mergable,omitempty"`
	MergedBy       GitUser   `json:"merged_by,omitempty"`
	Comments       int       `json:"comments,omitempty"`
	Commits        int       `json:"commits,omitempty"`
	Additions      int       `json:"additions,omitempty"`
	Deletions      int       `json:"deletions,omitempty"`
	ChangedFiles   int       `json:"changed_files,omitempty"`
}

type CommentLinks struct {
	Self        map[string]string `json:"self"`
	Html        map[string]string `json:"html"`
	PullRequest map[string]string `json:"pull_request"`
}

type PullComment struct {
	Url       string       `json:"url"`
	ID        int          `json:"id"`
	Body      Nstring      `json:"body"`
	Path      Nstring      `json:"path"`
	Position  int          `json:"position"`
	CommitId  Nstring      `json:"commit_id"`
	User      GitUser      `json:"user"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
	Links     CommentLinks `json:"_links"`
}

type CreateComment struct {
	Body      string `json:"body,omitreply"`
	CommitId  string `json:"commit_id,omitreply"`
	Path      string `json:"path,omitreply"`
	Position  int    `json:"position,omitreply"`
	InReplyTo int    `json:"in_reply_to,omitreply"`
}

// GitHub Doc - GitData: Pull Requests - List pull requests
// Url: https://api.github.com/repos/:owner/:repo/pulls?state=open&access_token=...
// Request Type: GET /repos/:owner/:repo/pulls
// Access Token: REQUIRED

func (github *GitHubClient) GetPullRequests(urlData map[string]string, state string) ([]PullRequest, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("One or more fields are missing and/or do not have content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/pulls?state=" + url.QueryEscape(strings.TrimSpace(state)))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		pullreq := &[]PullRequest{}
		pullreqJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(pullreqJson, pullreq); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*pullreq), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - GitData: Pull Requests - Get a single pull request
// Url: https://api.github.com/repos/:owner/:repo/pulls?state=open&access_token=...
// Request Type: GET /repos/:owner/:repo/pulls
// Access Token: REQUIRED
// 

func (github *GitHubClient) GetAPullRequest(urlData map[string]string) (*PullRequest, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New("One or more fields are missing and/or do not have content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/pulls/" + urlData["numbner"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		pullreq := &PullRequest{}
		pullreqJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(pullreqJson, pullreq); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return pullreq, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: GitData: Pull Requests - Create a pull request
// Url: https://api.github.com/repos/:owner/:repo/pulls?access_token=...
// Request Type: POST /repos/:owner/:repo/pulls
// Access Token: REQUIRED
// 

func (github *GitHubClient) CreatePullRequest(urlData, pullData map[string]string) (*PullRequest, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("Your repo in your urlData is either missing or has a length of zero.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}
	if ok := github.AssertMapStrings([]string{"title", "base", "head"}, pullData); !ok {
		if ok2 := github.AssertMapStrings([]string{"issue", "base", "head"}, pullData); !ok2 {
			return nil, errors.New("pullData is either missing data or value(s) don't contain non-whitespace chracters.")
		}
	}

	pullReader, err := github.CreateReader(pullData)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/pulls")
	res, err := github.Client.Post(apiUrl, "application/json", pullReader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		pullreq := &PullRequest{}
		pullreqJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(pullreqJson, pullreq); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return pullreq, nil
	}

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// 
// GitHub Doc: GitData: Pull Requests - Create a pull request
// Url: https://api.github.com/repos/:owner/:repo/pulls?access_token=...
// Request Type: POST /repos/:owner/:repo/pulls
// Access Token: REQUIRED
// 

func (github *GitHubClient) EditPullRequest(urlData, pullData map[string]string) (*PullRequest, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("Your repo in your urlData is either missing or has a length of zero.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	pullReader, err := github.CreateReader(pullData)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/pulls/" + urlData["number"])
	apiRequest, err := http.NewRequest("PATCH", apiUrl, pullReader)
	if err != nil {
		return nil, err
	}

	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		pullreq := &PullRequest{}
		pullreqJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(pullreqJson, pullreq); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return pullreq, nil
	}

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// 
// GitHub Doc - GitData: Pull Requests - Get a single pull request
// Url: https://api.github.com/repos/:owner/:repo/pulls/:number/commits?access_token=...
// Request Type: GET /repos/:owner/:repo/pulls/:number/commits
// Access Token: REQUIRED
// 

func (github *GitHubClient) GetPullCommits(urlData map[string]string) (*Commits, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New("One or more fields are missing and/or do not have content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/pulls/" + urlData["number"] + "/commits")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		commits := &Commits{}
		commitsJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(commitsJson, commits); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return commits, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - GitData: Pull Requests - List pull requests files
// Url: https://api.github.com/repos/:owner/:repo/pulls/:number/files?state=open&access_token=...
// Request Type: GET /repos/:owner/:repo/pulls/:number/files
// Access Token: REQUIRED
// 

func (github *GitHubClient) GetPullFiles(urlData map[string]string) ([]CommitFile, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New("One or more fields are missing and/or do not have content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/pulls/" + urlData["number"] + "/files")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		files := &[]CommitFile{}
		filesJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(filesJson, files); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*files), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - GitData: Pull Requests - Get if a pull request has been merged
// Url: https://api.github.com/repos/:owner/:repo/pulls/:number/merge?state=open&access_token=...
// Request Type: GET /repos/:owner/:repo/pulls/:number/merge
// Access Token: REQUIRED
// 

func (github *GitHubClient) HasPullMerged(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return false, errors.New("One or more fields are missing and/or do not have content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/pulls/" + urlData["number"] + "/merge")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode == 204 {
		files := []CommitFile{}
		filesJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return false, err
		}

		if err = json.Unmarshal(filesJson, files); err != nil {
			return false, err
		}

		github.getLimits(res)
		return true, nil
	}

	if res.StatusCode == 404 {
		return false, nil
	}

	return false, errors.New("Didn't receive 204/404 status from Github: " + res.Status)
}

// 
// GitHub Doc: GitData: Pull Requests - Merge a pull request (Merge Button™)
// Url: https://api.github.com/repos/:owner/:repo/pulls/:number/merge?access_token=...
// Request Type: PUT /repos/:owner/:repo/pulls/:number/merge
// Access Token: REQUIRED
// 

func (github *GitHubClient) MergePullRequest(urlData map[string]string, message string) (*PullMerge, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New("Your repo in your urlData is either missing or has a length of zero.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/pulls/" + urlData["number"] + "/merge?commit_message=" + url.QueryEscape(strings.TrimSpace(message)))
	apiRequest, err := http.NewRequest("PATCH", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 || res.StatusCode == 405 {
		pullreq := &PullMerge{}
		pullreqJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(pullreqJson, pullreq); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return pullreq, nil
	}

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// Review Comments Section
// 
// GitHub Doc - GitData: Pull Requests - List comments on a pull request
// Url: https://api.github.com/repos/:owner/:repo/pulls/:number/comments?access_token=...
// Request Type: GET /repos/:owner/:repo/pulls/:number/comments
// Access Token: REQUIRED
// 

func (github *GitHubClient) GetPullComments(urlData map[string]string) ([]PullComment, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New("One or more fields are missing and/or do not have content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/pulls/" + urlData["number"] + "/comments")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		comments := &[]PullComment{}
		commentsJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(commentsJson, comments); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*comments), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - GitData: Pull Requests - List comments in a repository
// Url: https://api.github.com/repos/:owner/:repo/pulls/comments?access_token=...
// Request Type: GET /repos/:owner/:repo/pulls/comments
// Access Token: REQUIRED
// 

func (github *GitHubClient) GetRepoPullComments(urlData, getData map[string]string) ([]PullComment, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("One or more fields are missing and/or do not have content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	urlStr := github.UrlDataConvert(getData)
	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/pulls/comments?" + urlStr)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		comments := &[]PullComment{}
		commentsJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(commentsJson, comments); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*comments), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - GitData: Pull Requests - Get a single comment
// Url: https://api.github.com/repos/:owner/:repo/pulls/comments/:number?access_token=...
// Request Type: GET /repos/:owner/:repo/pulls/comments/:number
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetARepoPullComment(urlData map[string]string) (*PullComment, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("One or more fields are missing and/or do not have content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/pulls/comments/" + urlData["number"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		comment := &PullComment{}
		commentJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(commentJson, comment); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return comment, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: GitData: Pull Requests - Create a comment
// Url: https://api.github.com/repos/:owner/:repo/pulls/:number/comments?access_token=...
// Request Type: POST /repos/:owner/:repo/pulls/:number/comments
// Access Token: REQUIRED
// 
func (github *GitHubClient) CreateRepoPullComment(urlData, commentData map[string]string) (*PullComment, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New("Your repo in your urlData is either missing or has a length of zero.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	comReader, err := github.CreateReader(commentData)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/pulls/" + urlData["number"] + "/comments")

	res, err := github.Client.Post(apiUrl, "application/json", comReader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		pullcom := &PullComment{}
		pullcomJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(pullcomJson, pullcom); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return pullcom, nil
	}

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// 
// GitHub Doc: GitData: Pull Requests - Edit a comment
// Url: https://api.github.com/repos/:owner/:repo/pulls/:number/comments?access_token=...
// Request Type: POST /repos/:owner/:repo/pulls/:number/comments
// Access Token: REQUIRED
// 
func (github *GitHubClient) EditRepoPullComment(urlData map[string]string, commentData *CreateComment) (*PullComment, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New("Your repo in your urlData is either missing or has a length of zero.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	comReader, err := github.CreateReader(commentData)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/pulls/comments/" + urlData["number"])
	apiRequest, err := http.NewRequest("PATCH", apiUrl, comReader)
	if err != nil {
		return nil, err
	}

	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		pullcom := &PullComment{}
		pullcomJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(pullcomJson, pullcom); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return pullcom, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: GitData: Pull Requests - Edit a comment
// Url: https://api.github.com/repos/:owner/:repo/pulls/:number/comments?access_token=...
// Request Type: POST /repos/:owner/:repo/pulls/:number/comments
// Access Token: REQUIRED
// 
func (github *GitHubClient) DeleteRepoPullComment(urlData map[string]string, commentData *CreateComment) (bool, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return false, errors.New("Your repo in your urlData is either missing or has a length of zero.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/pulls/comments/" + urlData["number"])
	apiRequest, err := http.NewRequest("DELETE", apiUrl, nil)
	if err != nil {
		return false, err
	}

	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode == 204 {
		github.getLimits(res)
		return true, nil
	}

	return false, errors.New("Didn't receive 204 status from Github: " + res.Status)
}
