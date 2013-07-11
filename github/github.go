// GitHub API v3 for the Go Programming Language
package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	GITSECRET = ""
	GITID     = ""
	BASEPATH  = ""
	APIURL    = "https://api.github.com"
)

type Nstring string

func (n *Nstring) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*string)(n))
}

type Nmap map[string]string

func (n *Nmap) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*Nmap)(n))
}

type GitHubClient struct {
	Type           string
	Token          string
	Login          string
	CallsLimit     int
	CallsRemaining int
	Client         *http.Client
}

type Markdown struct {
	Mode     string `json:"mode"`
	Context  string `json:"context"`
	Markdown string `json:"markdown"`
}

func NewGitHubClient(token, login string) *GitHubClient {
	httpClient := &http.Client{}

	gitClient := &GitHubClient{
		Type:           "oauth",
		Token:          token,
		Login:          login,
		CallsLimit:     5000,
		CallsRemaining: 5000,
		Client:         httpClient,
	}

	return gitClient
}

// GitHub v3 API - Utils to turn a single url into a full url making their management easier
//
// createUrl - path {string} - the path added to the base url https://api.github.com
// Also makes it easier to match with the docs
func (github *GitHubClient) createUrl(path string) string {
	apiUrl := ""

	if strings.Index(path, "?") == -1 {
		apiUrl = APIURL + path + "?access_token=" + url.QueryEscape(github.Token)
	} else {
		apiUrl = APIURL + path + "&access_token=" + url.QueryEscape(github.Token)
	}

	return apiUrl
}

func (github *GitHubClient) readResponse(res *http.Response, v interface{}) (interface{}, error) {
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return nil, err
	}

	github.getLimits(res)

	return v, nil
}

func (github *GitHubClient) AssertMapValue(key string, m map[string]interface{}) bool {
	if _, ok := m[key]; ok {
		return true
	}
	return false
}

func (github *GitHubClient) AssertMapValues(s []string, m map[string]interface{}) bool {
	for _, v := range s {
		if _, ok := m[v]; !ok {
			return false
		}
	}
	return true
}

func (github *GitHubClient) AssertMapString(key string, m map[string]string) bool {
	if v, ok := m[key]; ok && len(strings.TrimSpace(v)) != 0 {
		return true
	}
	return false
}

func (github *GitHubClient) AssertMapStrings(s []string, m map[string]string) bool {
	for _, key := range s {
		if val, ok := m[key]; !ok && strings.TrimSpace(val) != "" {
			return false
		}
	}
	return true
}

func (github *GitHubClient) UrlDataConvert(m map[string]string) string {
	s := ""
	for key, val := range m {
		if len(s) == 0 {
			s = s + url.QueryEscape(strings.TrimSpace(key)) + "=" + url.QueryEscape(strings.TrimSpace(val))
		} else {
			s = s + "&" + url.QueryEscape(strings.TrimSpace(key)) + "=" + url.QueryEscape(strings.TrimSpace(val))
		}
	}
	return s
}

func (github *GitHubClient) CreateReader(v interface{}) (*bytes.Reader, error) {
	jsonBuf, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(jsonBuf)
	return reader, nil
}

// Gets the limit headers from the response and saves them to the
// GitHubClient for determining rate limiting
func (github *GitHubClient) getLimits(res *http.Response) {
	remain, err := strconv.ParseInt(res.Header.Get("X-RateLimit-Remaining"), 10, 0)
	if err != nil {
		return
	}

	limit, err := strconv.ParseInt(res.Header.Get("X-RateLimit-Limit"), 10, 0)
	if err != nil {
		return
	}

	github.CallsRemaining = int(remain)
	github.CallsLimit = int(limit)
}

// *****************************
// * START: Markdown Section   *
// *****************************
//
// GitHub Docs: Render an arbitrary Markdown document
// Request Type: POST /markdown
// Access Token: NO Tokens needed
// Url: https://api.github.com/markdown?access_token=...
func (github *GitHubClient) RenderMarkdown(markdown *Markdown) (string, error) {
	if markdown.Markdown == "" {
		return "", errors.New("You must not send an empty string as the markdown contents.")
	}

	apiUrl := github.createUrl("/markdown")
	reader, err := github.CreateReader(markdown)

	res, err := github.Client.Post(apiUrl, "application/json", reader)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		htmlBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		html := string(htmlBytes)
		github.getLimits(res)
		return html, nil
	}

	return "", errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// ***************************
// *  END: Markdown Section  *
// ***************************
