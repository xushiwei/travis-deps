package github

// 
// Gist - Section of the GitHub API v3
// Includes the comments since I think they are much more useful in Gists, but rare for commits.
//
//	## Gist API
//		-  List gists
//		-  Get a single gist
//		-  Create a gist
//		-  Edit a gist
//		-  Star a gist
//		-  Unstar a gist
//		-  Check if a gist is starred
//		-  Fork a gist
//		-  Delete a gist
//
//	## Gist Comments API
//		-  List comments on a gist
//		-  Get a single comment
//		-  Create a comment
//		-  Edit a comment
//		-  Delete a comment
//		-  Custom media types
// 

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type GistFork struct {
	User      GitUser `json:"user"`
	Url       string  `json:"url"`
	CreatedAt string  `json:"created_at"`
}
type GistForks []GistFork

type GistHistory struct {
	User         GitUser        `json:"user"`
	Version      string         `json:"version"`
	Url          string         `json:"url"`
	ChangeStatus map[string]int `json:"change_status"`
	CommittedAt  string         `json:"committed_at"`
}
type GistHistories []GistHistory

type GistFile struct {
	Filename string `json:"filename"`
	Type     string `json:"type"`
	Language string `json:"language"`
	RawUrl   string `json:"raw_url"`
	Size     int    `json:"size"`
	Content  string `json:"content"`
}

type GistFiles map[string]GistFile

type Gist struct {
	Url         string        `json:"url,omitempty"`
	ForksUrl    string        `json:"forks_url,omitempty"`
	CommitsUrl  string        `json:"commits_url,omitempty"`
	Files       GistFiles     `json:"files,omitempty"`
	CreatedAt   string        `json:"created_at,omitempty"`
	UpdatedAt   string        `json:"updated_at,omitempty"`
	User        GitUser       `json:"user,omitempty"`
	ID          string        `json:"id,omitempty"`
	Public      bool          `json:"public,omitempty"`
	Description string        `json:"description,omitempty"`
	Comments    int           `json:"comments,omitempty"`
	CommentsUrl string        `json:"comments_url,omitempty"`
	HtmlUrl     string        `json:"html_url,omitempty"`
	GitPullUrl  string        `json:"git_pull_url,omitempty"`
	GitPushUrl  string        `json:"git_push_url,omitempty"`
	Forks       GistForks     `json:"forks,omitempty"`
	History     GistHistories `json:"history,omitempty"`
}

type Gists []Gist

type PostGistFile struct {
	Content  string `json:"content,omitempty"`
	Filename string `json:"filename,omitempty"`
}

type PostGist struct {
	Description string                   `json:"description,omitempty"`
	Public      bool                     `json:"public,omitempty"`
	Files       map[string]*PostGistFile `json:"files,omitempty"`
}

type GistComment struct {
	ID        int     `json:"id"`
	Url       string  `json:"url"`
	Body      string  `json:"body"`
	User      GitUser `json:"user"`
	CreatedAt string  `json:"created_at"`
}
type GistComments []GistComment

// 
// GitHub Doc: Gists: List the authenticated user’s gists
// Url: https://api.github.com/gists?access_token=...
// Request Type: GET /gists
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetGists(getData map[string]string) (*Gists, error) {
	urlStr := github.UrlDataConvert(getData)
	apiUrl := github.createUrl("/gists?" + urlStr)

	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		gists := &Gists{}
		gistsJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(gistsJson, gists); err != nil {
			return nil, err
		}
		github.getLimits(res)
		return gists, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: Gists: List the authenticated user’s starred gists only
// Url: https://api.github.com/gists/starred?access_token=...
// Request Type: GET /gists/starred
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetStarredGists(since string) (*Gists, error) {
	apiUrl := ""
	if since == "" {
		apiUrl = github.createUrl("/gists/starred")
	} else {
		apiUrl = github.createUrl("/gists/starred?since=" + url.QueryEscape(since))
	}

	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		gists := &Gists{}
		gistsJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(gistsJson, gists); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return gists, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: Gists: List the authenticated user’s public gists only
// Url: https://api.github.com/gists/public?access_token=...
// Request Type: GET /gists/public
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetPublicGists(since string) (*Gists, error) {
	apiUrl := ""
	if since == "" {
		apiUrl = github.createUrl("/gists/public")
	} else {
		apiUrl = github.createUrl("/gists/public?since=" + url.QueryEscape(since))
	}

	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		gists := &Gists{}
		gistsJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(gistsJson, gists); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return gists, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: Gists: Get a single gist
// Url: https://api.github.com/gists/:id?access_token=...
// Request Type: GET /gists/:id
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetGistById(id string) (*Gist, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("The id must have a length greater then zero.")
	}

	apiUrl := github.createUrl("/gists/" + id)

	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		gist := &Gist{}
		gistJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(gistJson, gist); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return gist, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: Gists: Create a gist
// Url: https://api.github.com/gists?access_token=...
// Request Type: POST /gists
// Access Token: REQUIRED
// 
func (github *GitHubClient) CreateGist(postGist *PostGist) (*Gist, error) {
	fLen := len(postGist.Files)
	if fLen > 0 {
		return nil, errors.New("There are no files in your Gist. Please add a file to your Gist.")
	}

	apiUrl := github.createUrl("/gists")
	apiReader, err := github.CreateReader(postGist)
	if err != nil {
		return nil, err
	}

	res, err := github.Client.Post(apiUrl, "application/json", apiReader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		gist := &Gist{}
		gistJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(gistJson, gist); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return gist, nil
	}

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// 
// GitHub Doc: Gists: Edit a gist
// Url: https://api.github.com/gists:id?access_token=...
// Request Type: PATCH /gists/:id
// Access Token: REQUIRED
// 
func (github *GitHubClient) EditGist(id string, postGist *PostGist) (*Gist, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("The id must have a length greater then zero.")
	}

	apiUrl := github.createUrl("/gists/" + id)
	apiReader, err := github.CreateReader(postGist)
	if err != nil {
		return nil, err
	}

	apiRequest, err := http.NewRequest("PATCH", apiUrl, apiReader)
	if err != nil {
		return nil, err
	}
	apiRequest.ContentLength = int64(apiReader.Len())

	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		gist := &Gist{}
		gistJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(gistJson, gist); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return gist, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: Gists: Star a gist
// Url: https://api.github.com/gists/:id/star?access_token=...
// Request Type: PUT /gists/:id/star
// Access Token: REQUIRED
// 
func (github *GitHubClient) StarGist(id string) (bool, error) {
	if strings.TrimSpace(id) == "" {
		return false, errors.New("The id must have a length greater then zero.")
	}

	apiUrl := github.createUrl("/gists/" + id + "/star")
	apiRequest, err := http.NewRequest("PUT", apiUrl, nil)
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

// 
// GitHub Doc: Gists: Untar a gist
// Url: https://api.github.com/gists/:id/star?access_token=...
// Request Type: DELETE /gists/:id/star
// Access Token: REQUIRED
// 
func (github *GitHubClient) UnstarGist(id string) (bool, error) {
	if strings.TrimSpace(id) == "" {
		return false, errors.New("The id must have a length greater then zero.")
	}

	apiUrl := github.createUrl("/gists/" + id + "/star")
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
	if res.StatusCode == 404 {
		github.getLimits(res)
		return false, nil
	}

	return false, errors.New("Didn't receive 204 or 404 status from Github: " + res.Status)
}

// 
// GitHub Doc: Gists: Fork a gist
// Url: https://api.github.com/gists/:id/forks?access_token=...
// Request Type: POST /gists/:id/forks
// Access Token: REQUIRED
// 
func (github *GitHubClient) ForkGist(id string) (*Gist, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("The id must have a length greater then zero.")
	}

	apiUrl := github.createUrl("/gists/" + id + "/forks")
	res, err := github.Client.Post(apiUrl, "text/html", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		gist := &Gist{}
		gistJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(gistJson, gist); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return gist, nil
	}

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// 
// GitHub Doc: Gists: Delete a gist
// Url: https://api.github.com/gists/:id?access_token=...
// Request Type: DELETE /gists/:id
// Access Token: REQUIRED
// 
func (github *GitHubClient) DeleteGist(id string) (bool, error) {
	if strings.TrimSpace(id) == "" {
		return false, errors.New("The id must have a length greater then zero.")
	}

	apiUrl := github.createUrl("/gists/" + id)
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

// Gist - Comments Section
// 
// GitHub Doc: Gists - Comments: List comments on a gist
// Url: https://api.github.com/gists/:gist_id/comments?access_token=...
// Request Type: GET /gists/:gist_id/comments
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetGistComments(gistId string) (*GistComments, error) {
	if strings.TrimSpace(gistId) == "" {
		return nil, errors.New("The gistId must have a length greater then zero.")
	}

	apiUrl := github.createUrl("/gists/" + gistId + "/comments")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		comments := &GistComments{}
		commentJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(commentJson, comments); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return comments, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: Gists - Comments: Get a single commment of a gist
// Url: https://api.github.com/gists/:gist_id/comments/:id?access_token=...
// Request Type: GET /gists/:gist_id/comments/:id
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetGistCommentById(gistId, commentId string) (*GistComment, error) {
	if strings.TrimSpace(gistId) == "" || strings.TrimSpace(commentId) == "" {
		return nil, errors.New("gistId and commentId are both must have a length greater then zero")
	}

	apiUrl := github.createUrl("/gists/" + gistId + "/comments/" + commentId)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		comment := &GistComment{}
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
// GitHub Doc: Gists - Comment: Create a comment
// Url: https://api.github.com/gists/:gist_id/comments?access_token=...
// Request Type: POST /gists/:gist_id/comments
// Access Token: REQUIRED
// 
func (github *GitHubClient) CreateGistComment(gistId string, commentBody string) (*GistComment, error) {
	if strings.TrimSpace(gistId) == "" {
		return nil, errors.New("The gistId must have a length greater then zero.")
	}
	if strings.TrimSpace(commentBody) == "" {
		return nil, errors.New("The commentBody must have a length greater then zero.")
	}

	commentMap := make(map[string]string)
	commentMap["body"] = commentBody

	apiUrl := github.createUrl("/gists/" + gistId + "/comments")
	apiReader, err := github.CreateReader(commentMap)
	if err != nil {
		return nil, err
	}

	res, err := github.Client.Post(apiUrl, "application/json", apiReader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		comment := &GistComment{}
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

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// 
// GitHub Doc: Gists - Comments: Edit a comment
// Url: https://api.github.com/gists/:gist_id/comments/:id?access_token=...
// Request Type: PATCH /gists/:gist_id/comments/:id
// Access Token: REQUIRED
// 
func (github *GitHubClient) EditGistComment(gistId string, commentData map[string]string) (*GistComment, error) {
	if ok := github.AssertMapStrings([]string{"id", "body"}, commentData); !ok {
		return nil, errors.New("There is comment data missing. Both body and id are required and must have a length greater then zero.")
	}
	if strings.TrimSpace(gistId) == "" {
		return nil, errors.New("gistId must have a length greater then zero.")
	}

	apiUrl := github.createUrl("/gists/" + gistId + "/comments/" + commentData["id"])
	apiReader, err := github.CreateReader(commentData)
	if err != nil {
		return nil, err
	}

	apiRequest, err := http.NewRequest("PATCH", apiUrl, apiReader)
	if err != nil {
		return nil, err
	}

	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		comment := &GistComment{}
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
// GitHub Doc: Gists - Comments: Delete a comment
// Url: https://api.github.com/gists/:gist_id/comments/:id?access_token=...
// Request Type: DELETE /gists/:gist_id/comments/:id
// Access Token: REQUIRED
// 
func (github *GitHubClient) DeleteGistComment(gistId, commentId string) (bool, error) {
	if strings.TrimSpace(gistId) == "" || strings.TrimSpace(commentId) == "" {
		return false, errors.New("gistId and commentId are both must have a length greater then zero")
	}

	apiUrl := github.createUrl("/gists/" + gistId + "/comments/" + commentId)
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
