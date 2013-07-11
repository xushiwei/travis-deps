package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Message struct {
	Message string
}

type SubCommit struct {
	Url       string             `json:"url"`
	SHA       string             `json:"sha"`
	Message   string             `json:"message"`
	Author    map[string]Nstring `json:"author"`
	Committer map[string]Nstring `json:"committer"`
	Tree      map[string]Nstring `json:"tree"`
}

type CommitFile struct {
	Filename  string  `json:"filename"`
	Additions int     `json:"additions"`
	Deletions int     `json:"deletions"`
	Changes   int     `json:"changes"`
	Status    Nstring `json:"status"`
	RawUrl    string  `json:"raw_url"`
	BlobUrl   string  `json:"blob_url"`
	Patch     Nstring `json:"patch"`
}

type GitUser struct {
	Login             string  `json:"login,omitempty"`
	ID                int     `json:"id,omitempty"`
	Avatar            Nstring `json:"avatar_url,omitempty"`
	Url               string  `json:"url,omitempty"`
	Gravatar          Nstring `json:"gravatar_id,omitempty"`
	HtmlUrl           string  `json:"html_url,omitempty"`
	FollowersUrl      string  `json:"followers_url,omitempty"`
	FollowingUrl      string  `json:"following_url,omitempty"`
	GistsUrl          string  `json:"gists_url,omitempty"`
	StarredUrl        string  `json:"starred_url,omitempty"`
	SubscriptionsUrl  string  `json:"subscriptions_url,omitempty"`
	OrganizationUrl   string  `json:"organizations_url,omitempty"`
	ReposUrl          string  `json:"repos_url,omitempty"`
	EventsUrl         string  `json:"events_url,omitempty"`
	ReceivedEventsUrl string  `json:"received_events_url,omitempty"`
	Type              string  `json:"type,omitempty"`
}

type Status struct {
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   Nstring `json:"updated_at"`
	State       Nstring `json:"state"`
	ID          int     `json:"id"`
	Url         string  `json:"url"`
	TargetUrl   Nstring `json:"target_url"`
	Description Nstring `json:"description"`
	Creator     GitUser `json:"creator"`
}
type Statuses []Status

type Commit struct {
	Url       string               `json:"url"`
	SHA       string               `json:"sha"`
	Commit    SubCommit            `json:"commit"`
	Author    GitUser              `json:"author"`
	Committer GitUser              `json:"committer"`
	Parents   []map[string]Nstring `json:"parents"`
	Stats     map[string]int       `json:"stats,omitempty"`
	Files     []CommitFile         `json:"files,omitempty"`
}
type Commits []Commit

type Repo struct {
	ID               int             `json:"id"`
	Name             string          `json:"name"`
	FullName         string          `json:"full_name"`
	Owner            User            `json:"owner"`
	Private          bool            `json:"private"`
	HtmlURL          string          `json:"html_url"`
	Description      Nstring         `json:"description"`
	Fork             bool            `json:"fork"`
	Watchers         int             `json:"watchers"`
	URL              string          `json:"url"`
	Homepage         Nstring         `json:"homepage"`
	ForksUrl         string          `json:"forks_url"`
	KeysUrl          string          `json:"keys_url"`
	CollaboratorsUrl string          `json:"collaborators_url"`
	TeamsUrl         string          `json:"teams_url"`
	HooksUrl         string          `json:"hooks_url"`
	IssueEventsUrl   string          `json:"issue_events_url"`
	EventsUrl        string          `json:"events_url"`
	AssigneeUrl      string          `json:"assignee_url"`
	BranchesUrl      string          `json:"branches_url"`
	TagsUrl          string          `json:"tags_url"`
	StatusesUrl      string          `json:"statuses_url"`
	LanguagesUrl     string          `json:"languages_url"`
	StargazersUrl    string          `json:"stargazers_url"`
	ContributorsUrl  string          `json:"contributors_url"`
	BlobsUrl         string          `json:"blobs_url"`
	TreesUrl         string          `json:"trees_url"`
	GitRefsUrl       string          `json:"git_refs_url"`
	GitTagsUrl       string          `json:"git_tags_url"`
	SubscribersUrl   string          `json:"subscribers_url"`
	SubscriptionUrl  string          `json:"subscription_url"`
	CommitsUrl       string          `json:"commits_url"`
	GitCommitsUrl    string          `json:"git_commits_url"`
	CommentsUrl      string          `json:"comments_url"`
	IssueCommentUrl  string          `json:"issue_comment_url"`
	CompareUrl       string          `json:"compare_url"`
	MergesUrl        string          `json:"merges_url"`
	ArchiveUrl       string          `json:"archive_url"`
	DownloadsUrl     Nstring         `json:"downloads_url"`
	IssuesUrl        string          `json:"issues_url"`
	PullsUrl         string          `json:"pulls_url"`
	MilestonesUrl    string          `json:"milestones_url"`
	NotificationsUrl string          `json:"notifications_url"`
	LabelsUrl        string          `json:"labels_url"`
	DefaultBranch    string          `json:"defaul_url"`
	SSHURL           Nstring         `json:"ssh_url"`
	CloneURL         string          `json:"clone_url"`
	MasterBranch     string          `json:"master_branch"`
	CreatedAt        Nstring         `json:"created_at"`
	UpdatedAt        Nstring         `json:"updated_at"`
	PushedAt         Nstring         `json:"pushed_at"`
	Language         Nstring         `json:"language"`
	Size             int             `json:"size"`
	ContentsUrl      string          `json:"contents_url"`
	BlobUrl          string          `json:"blobs_url"`
	Permissions      map[string]bool `json:"permissions"`
	MirrorUrl        Nstring         `json:"mirror_url"`
	OpenIssues       int             `json:"opne_issues"`
	Forks            int             `json:"forks"`
	ForksCount       int             `json:"forks_count"`
	WatchersCount    int             `json:"watchers_count"`
	HasIssues        bool            `json:"has_issues"`
	HasDownloads     bool            `json:"has_downloads"`
	HasWiki          bool            `json:"has_wiki"`
}

type Hook struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	Url       string            `json:"url"`
	Events    []string          `json:"events"`
	Active    bool              `json:"active"`
	Config    map[string]string `json:"config"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}
type Hooks []Hook

type Repos []Repo

type NewRepo struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Homepage    string `json:"homepage,omitempty"`
	Private     bool   `json:"private,omitempty"`
	Issues      bool   `json:"has_issues,omitempty"`
	Wiki        bool   `json:"has_wiki,omitempty"`
	Downloads   bool   `json:"has_downloads,omitempty"`
	AutoInit    bool   `json:"auto_init,omitempty"`
	TeamId      int    `json:"team_id,omitempty"`
	IgnoreTemp  string `json:"gitignore_template,omitempty"`
}

type Branch struct {
	Name     string                 `json:"name"`
	Commit   map[string]interface{} `json:"commit,omitempty"`
	Author   map[string]interface{} `json:"author,omitempty"`
	Parents  []map[string]string    `json:"parents,omitempty"`
	Url      string                 `json:"url,omitempty"`
	Commiter map[string]interface{} `json:"commiter,omitempty"`
	Links    map[string]string      `json:"_links,omitempty"`
}

type Branches []Branch

type Contributor struct {
	Login         string `json:"login"`
	ID            int    `json:"id"`
	Avatar        string `json:"avatar_url"`
	Url           string `json:"url"`
	Contributions int    `json:"contributions"`
}
type Contributors []Contributor

type Collaborator GitUser
type Collaborators []Collaborator

type Tag struct {
	Name   string            `json:"name"`
	Commit map[string]string `json:"commit"`
	ZipUrl string            `json:"zipball_url"`
	TarUrl string            `json:"tarball_url"`
}
type Tags []Tag

type Content struct {
	Url      string             `json:"url"`
	Type     string             `json:"type"`
	Encoding Nstring            `json:"encoding"`
	GitUrl   string             `json:"git_url"`
	Path     string             `json:"path"`
	HtmlUrl  string             `json:"html_url"`
	Size     int                `json:"size,omitempty"`
	Links    map[string]Nstring `json:"_links"`
	Name     Nstring            `json:"name"`
	Content  Nstring            `json:"content"`
	SHA      string             `json:"sha"`
}
type Contents []Content

type Fork struct {
	ID            int     `json:"id"`
	Owner         GitUser `json:"owner"`
	Name          string  `json:"name"`
	FullName      string  `json:"full_name"`
	Description   Nstring `json:"description"`
	Private       bool    `json:"private"`
	Fork          bool    `json:"private"`
	Url           string  `json:"url"`
	HtmlUrl       string  `json:"html_url"`
	SSHUrl        string  `json:"ssh_url"`
	GitUrl        string  `json:"git_url"`
	SvnUrl        string  `json:"svn_url"`
	MirrorUrl     Nstring `json:"mirror_url"`
	Homepage      Nstring `json:"homepage"`
	Language      Nstring `json:"language,omitempty"`
	Forks         int     `json:"forks"`
	ForksCount    int     `json:"forks_count"`
	Watchers      int     `json:"watchers"`
	WatchersCount int     `json:"watchers_count"`
	Size          int     `json:"size"`
	Master        string  `json:"master_branch"`
	OpenIssues    int     `json:"open_issues"`
	PushedAt      Nstring `json:"pushed_at"`
	CreatedAt     Nstring `json:"created_at"`
	UpdatedAt     Nstring `json:"updated_at"`
}
type Forks []Fork

//Start Repo Keys
type Key struct {
	ID    int    `json:"id,omitempty"`
	Key   string `json:"key,omitempty"`
	Url   string `json:"url,omitempty"`
	Title string `json:"title,omitempty"`
}
type Keys []Key

// 
// GitHub Doc: List your repositories - List repositories for the authenticated user
// Url: https://api.github.com/user/repos?access_token=...
// Request Type: GET 
// Access Token: REQUIRED
// Options: 
//		type - string: all, owner, public, private, member. Default: all
//		sort - string: created, updated, pushed, full_name, default: full_name
//		direction - string: asc or desc, default: when using full_name: asc, otherwise desc
// 
func (github *GitHubClient) GetUserRepos(getData map[string]string) (*Repos, error) {
	optionString := github.UrlDataConvert(getData)
	apiUrl := github.createUrl("/user/repos?" + optionString)

	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		repoJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		repos := &Repos{}
		if err = json.Unmarshal(repoJson, repos); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return repos, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: Repo - GET - GET /repos/:owner/:repo
// Url: https://api.github.com/user/:repo/:owner?access_token=...
// Request Type: GET
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetRepo(urlData map[string]string) (*Repo, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("Data to create the url is missing. Both user and id are required fields for this map.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		repo := &Repo{}
		repoJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(repoJson, repo); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return repo, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: List organization repositories. - List repositories for the specified org.
// Url: https://api.github.com/orgs/:org/repos?access_token=...
// Request Type: GET /orgs/:org/repos
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetOrgRepos(org, repoType string) (*Repos, error) {
	repoType = url.QueryEscape(strings.TrimSpace(repoType))
	apiUrl := ""
	if repoType == "" {
		apiUrl = github.createUrl("/orgs/" + org + "/repos")
	} else {
		apiUrl = github.createUrl("/orgs/" + org + "/repos?type=" + repoType)
	}

	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		repos := &Repos{}
		repoJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(repoJson, repos); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return repos, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: Create - Create a new repository for the authenticated user. OAuth users must supply repo scope.
// Url: https://api.github.com/user/repos?access_token=...
// Request Type: POST /user/repos
// Access Token: REQUIRED
// 
func (github *GitHubClient) CreateRepo(newRepo *NewRepo) (*Repo, error) {
	apiUrl := github.createUrl("/user/repos")
	if newRepo.Name != "" { // If there is a name it is good to go
		repoReader, err := github.CreateReader(newRepo)

		res, err := github.Client.Post(apiUrl, "application/json", repoReader)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode == 201 {
			repo := &Repo{}
			repoJson, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}

			if err = json.Unmarshal(repoJson, repo); err != nil {
				return nil, err
			}

			github.getLimits(res)
			return repo, nil
		}

		return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
	}
	return nil, errors.New("There was no name given to the repo you wanted to create")
}

//ORGANIZATION VERSION
func (github *GitHubClient) CreateOrgRepo(newRepo *NewRepo, company string) (*Repo, error) {
	if company == "" {
		return nil, errors.New("There is no company name try using CreateRepo instead")
	}

	apiUrl := github.createUrl("/user/repos")
	if newRepo.Name != "" { // If there is a name it is good to go
		repoBuffer, err := json.Marshal(newRepo)
		if err != nil {
			return nil, err
		}
		repoReader := bytes.NewReader(repoBuffer)

		res, err := github.Client.Post(apiUrl, "application/json", repoReader)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode == 201 {
			repo := &Repo{}
			repoJson, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}

			if err = json.Unmarshal(repoJson, repo); err != nil {
				return nil, err
			}

			github.getLimits(res)
			return repo, nil
		}

		return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
	}
	return nil, errors.New("There was no name given to the repo you wanted to create")
}

// 
// GitHub  Docs: Edit
// Url: https://api.github.com/repos/:owner/:repo?access_token=...
// Request Type: PATCH /repos/:owner/:repo
// Access Token: REQUIRED
// 
func (github *GitHubClient) EditRepo(urlData map[string]string, editRepo *NewRepo) (*Repo, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("Data to create the url is missing. Both user and id are required fields for this map.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	//Setup Request Data
	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"])
	repoBuffer, err := json.Marshal(editRepo)
	if err != nil {
		return nil, err
	}
	repoReader := bytes.NewReader(repoBuffer)                       //Reader
	apiRequest, err := http.NewRequest("PATCH", apiUrl, repoReader) // PATCH Request 
	if err != nil {
		return nil, err
	}

	// Execute Request
	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	//If Request is Successful then return data
	if res.StatusCode == 200 {
		repo := &Repo{}
		repoJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(repoJson, repo); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return repo, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: Repo - List contributors
// Url: https://api.github.com/repos/:owner/:repo/contributors?access_token=...
// Request Type: GET /repos/:owner/:repo/contributors
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetRepoContributors(urlData map[string]string, anon string) (*Contributors, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("Data to create the url is missing. Both user and id are required fields for this map.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	anonStr := ""
	if anon == "1" || anon == "true" {
		anonStr = "?anon=" + anon
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/contributors" + anonStr)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		contribs := &Contributors{}
		ContribJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(ContribJson, contribs); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return contribs, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: Repo - List languages
// Url: https://api.github.com/repos/:owner/:repo/contributors?access_token=...
// Request Type: GET /repos/:owner/:repo/contributors
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetRepoLanguages(urlData map[string]string) (*map[string]int, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("Data to create the url is missing. Both user and id are required fields for this map.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/languages")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		langMap := &map[string]int{}
		langJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(langJson, langMap); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return langMap, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: Repo - List Teams
// Url: https://api.github.com/repos/:owner/:repo/teams?access_token=...
// Request Type: GET /repos/:owner/:repo/teams
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetRepoTeams(urlData map[string]string) (*Teams, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("Data to create the url is missing")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/teams")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		teamJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		teams := &Teams{}

		err = json.Unmarshal(teamJson, teams)
		if err != nil {
			return nil, err
		}

		github.getLimits(res)
		return teams, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: Repo - List Tags
// Url: https://api.github.com/repos/:owner/:repo/tags?access_token=...
// Request Type: GET /repos/:owner/:repo/tags
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetRepoTags(urlData map[string]string) (*Tags, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("Data to create the url is missing")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/tags")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		tags := &Tags{}
		tagsJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(tagsJson, tags); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return tags, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: Repo - List Branches
// Url: https://api.github.com/repos/:owner/:repo/branches?access_token=...
// Request Type: GET /repos/:owner/:repo/branches
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetRepoBranches(urlData map[string]string) (*Branches, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("Data to create the url is missing")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/branches")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		branches := &Branches{}
		branchJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(branchJson, branches); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return branches, nil
	}

	github.getLimits(res)
	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: Repo - Get Branch
// Url: https://api.github.com/repos/:owner/:repo/branches/:branch?access_token=...
// Request Type: GET /repos/:owner/:repo/branches/:branch
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetRepoBranch(urlData map[string]string) (*Branch, error) {
	if ok := github.AssertMapStrings([]string{"repo", "branch"}, urlData); !ok {
		return nil, errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "-" + urlData["repo"] + "/branches/" + urlData["branch"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		branch := &Branch{}
		branchJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(branchJson, branch); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return branch, nil
	}

	github.getLimits(res)
	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub  Docs: Delete Repo
// Url: https://api.github.com/repos/:owner/:repo?access_token=...
// Request Type: DELETE /repos/:owner/:repo
// Access Token: REQUIRED
// 
func (github *GitHubClient) DeleteRepo(urlData map[string]string) error {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	//Setup Request Data
	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"])
	apiRequest, err := http.NewRequest("DELETE", apiUrl, nil) // DELETE Request 
	if err != nil {
		return err
	}

	// Execute Request
	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	//If Request is Successful then return data
	if res.StatusCode == 204 {
		github.getLimits(res)
		return nil
	}

	github.getLimits(res)
	return errors.New("Didn't receive 204 status from Github: " + res.Status)
}

// 
// GitHub  Docs: Get Repo + Path Contents
// Url: https://api.github.com/repos/:owner/:repo/contents/:path?access_token=...
// Request Type: GET /repos/:owner/:repo/contents/:path
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetPathContents(urlData map[string]string) (*Contents, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	//Setup Request Data
	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/contents/" + urlData["path"])
	// Execute Request
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	//If Request is Successful then return data
	if res.StatusCode == 200 {
		content := &Contents{}
		contentJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(contentJson, content); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return content, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub  Docs: Get the README - This method returns the preferred README for a repository.
// Url: https://api.github.com/repos/:owner/:repo/contents/:path?access_token=...
// Request Type: GET /repos/:owner/:repo/readme
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetReadme(urlData map[string]string) (*Content, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	//Setup Request Data
	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "-" + urlData["repo"] + "/readme")
	// Execute Request
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	//If Request is Successful then return data
	if res.StatusCode == 200 {
		content := &Content{}
		contentJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(contentJson, content); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return content, nil
	}

	github.getLimits(res)
	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub  Docs: Get archive link - For private repositories, these links are temporary and expire quickly.
// Url: https://api.github.com/repos/:owner/:repo/:archive_format/:ref?access_token=...
// Request Type: GET /repos/:owner/:repo/:archive_format/:ref
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetZip(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return false, errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	if len(urlData["format"]) == 0 {
		urlData["format"] = "zipball"
	}
	ext := ".zip"
	if urlData["format"] == "tarball" {
		ext = ".tar.gz"
	}

	zipOut, err := os.Create(BASEPATH + "github/zip/" + urlData["owner"] + "-" + urlData["repo"] + "-" + urlData["branch"] + ext)
	if err != nil {
		return false, err
	}
	defer zipOut.Close()

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/" + urlData["format"] + "/" + urlData["branch"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	_, err = io.Copy(zipOut, res.Body)
	if err != nil {
		github.getLimits(res)
		return false, err
	}

	github.getLimits(res)
	return true, nil
}

// Start of Collaborators
// 
// GitHub  Docs: List Collabs
// Url: https://api.github.com/repos/:owner/:repo/collaborators?access_token=...
// Request Type: GET /repos/:owner/:repo/collaborators
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetCollabs(urlData map[string]string) (*Collaborators, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/collaborators")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		collabs := &Collaborators{}
		collabJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(collabJson, collabs); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return collabs, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub  Docs: Get - Is User a Collab
// Url: https://api.github.com/repos/:owner/:repo/collaborators/:user?access_token=...
// Request Type: GET /repos/:owner/:repo/collaborators/:user
// Access Token: REQUIRED
// 
func (github *GitHubClient) IsCollab(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapStrings([]string{"repo", "user"}, urlData); !ok {
		return false, errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/collaborators/" + urlData["user"])
	res, err := github.Client.Get(apiUrl)
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

	github.getLimits(res)
	return false, errors.New("Didn't receive 204 or 404 status from Github: " + res.Status)
}

// 
// GitHub  Docs: Add Collab
// Url: https://api.github.com/repos/:owner/:repo/collaborators/:user?access_token=...
// Request Type: PUT /repos/:owner/:repo/collaborators/:user
// Access Token: REQUIRED
// 
func (github *GitHubClient) AddCollab(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapStrings([]string{"repo", "user"}, urlData); !ok {
		return false, errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/collaborators/" + urlData["user"])
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

	if res.StatusCode == 404 {
		github.getLimits(res)
		return false, nil
	}

	github.getLimits(res)
	return false, errors.New("Didn't receive 204 or 404 status from Github: " + res.Status)
}

// 
// GitHub  Docs: Delete Collab
// Url: https://api.github.com/repos/:owner/:repo/collaborators/:user?access_token=...
// Request Type: PUT /repos/:owner/:repo/collaborators/:user
// Access Token: REQUIRED
// 
func (github *GitHubClient) DeleteCollab(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapStrings([]string{"repo", "user"}, urlData); !ok {
		return false, errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/collaborators/" + urlData["user"])
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

	github.getLimits(res)
	return false, errors.New("Didn't receive 204 or 404 status from Github: " + res.Status)
}

// Start of Forks - List, Create Forks

// 
// GitHub  Docs: List Forks
// Url: https://api.github.com/repos/:owner/:repo/collaborators?access_token=...
// Request Type: GET /repos/:owner/:repo/collaborators
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetForks(urlData map[string]string, getData map[string]string) (*Forks, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	urlStr := github.UrlDataConvert(getData)
	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/forks?" + urlStr)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		forks := &Forks{}
		forkJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(forkJson, forks); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return forks, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub  Docs: Create Fork
// Url: https://api.github.com/repos/:owner/:repo/collaborators?access_token=...
// Request Type: POST /repos/:owner/:repo/collaborators
// Access Token: REQUIRED
// 
func (github *GitHubClient) CreateFork(urlData map[string]string, org string) (*Fork, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/forks")
	var (
		res *http.Response
		err error
	)

	if org != "" {
		apiReader, err := github.CreateReader(map[string]string{"organization": "org"})
		if err != nil {
			return nil, err
		}

		res, err = github.Client.Post(apiUrl, "application/json", apiReader)
	} else {
		res, err = github.Client.Post(apiUrl, "application/json", nil)
	}

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 202 {
		fork := &Fork{}
		forkJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(forkJson, fork); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return fork, nil
	}

	return nil, errors.New("Didn't receive 202 status from Github: " + res.Status)
}

// 
// GitHub  Docs: Repo: Keys - List
// Url: https://api.github.com/repos/:owner/:repo/keys?access_token=...
// Request Type: GET /repos/:owner/:repo/keys
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetRepoKeys(urlData map[string]string) (*Keys, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/keys")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		keys := &Keys{}
		keyJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(keyJson, keys); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return keys, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub  Docs: Repo: Get a Key by ID
// Url: https://api.github.com/repos/:owner/:repo/keys/:id?access_token=...
// Request Type: GET /repos/:owner/:repo/keys/:id
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetRepoKey(urlData map[string]string) (*Key, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/keys/" + urlData["id"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		key := &Key{}
		keyJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(keyJson, key); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return key, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub  Docs: Repo: Create A Key
// Url: https://api.github.com/repos/:owner/:repo/keys/:id?access_token=...
// Request Type: POST /repos/:owner/:repo/keys
// Access Token: REQUIRED
// 
func (github *GitHubClient) CreateRepoKey(urlData map[string]string, key *map[string]string) (*Key, error) {
	if ok := github.AssertMapStrings([]string{"repo", "id"}, urlData); !ok {
		return nil, errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiReader := strings.NewReader(`{ "key": "` + (*key)["key"] + `", "title": "` + (*key)["title"] + `" }`)
	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/keys")
	res, err := github.Client.Post(apiUrl, "application/json", apiReader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		key := &Key{}
		keyJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(keyJson, key); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return key, nil
	}

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// 
// GitHub  Docs: Repo: Edit A Key
// Url: https://api.github.com/repos/:owner/:repo/keys/:id?access_token=...
// Request Type: PATCH /repos/:owner/:repo/keys/:id
// Access Token: REQUIRED
// 
func (github *GitHubClient) EditRepoKey(urlData map[string]string, key *map[string]string) (*Key, error) {
	if ok := github.AssertMapStrings([]string{"repo", "id"}, urlData); !ok {
		return nil, errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiReader := strings.NewReader(`{ "key": "` + (*key)["key"] + `", "title": "` + (*key)["title"] + `" }`)
	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/keys/" + urlData["id"])
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
		key := &Key{}
		keyJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(keyJson, key); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return key, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub  Docs: Repo: Delete A Key
// Url: https://api.github.com/repos/:owner/:repo/keys/:id?access_token=...
// Request Type: DELETE /repos/:owner/:repo/keys/:id
// Access Token: REQUIRED
// 
func (github *GitHubClient) DeleteRepoKey(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapStrings([]string{"repo", "id"}, urlData); !ok {
		return false, errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/keys/" + urlData["id"])
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

//Commits Section
// 
// GitHub  Docs: Repo: Commits -  List commits on a repository
// Url: https://api.github.com/repos/:owner/:repo/commits?access_token=...
// Request Type: GET /repos/:owner/:repo/commits
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetRepoCommits(urlData map[string]string, params map[string]string) (*Commits, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	paramUrl := github.UrlDataConvert(params)
	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/commits?" + paramUrl)
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
// GitHub  Docs: Repo: Commits - Get a single commit
// Url: https://api.github.com/repos/:owner/:repo/commits/:sha?access_token=...
// Request Type: /repos/:owner/:repo/commits/:sha
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetACommits(urlData map[string]string, params map[string]string) (*Commit, error) {
	if ok := github.AssertMapStrings([]string{"repo", "sha"}, urlData); !ok {
		return nil, errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	paramUrl := github.UrlDataConvert(params)
	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/commits/" + urlData["sha"] + "?" + paramUrl)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		commit := &Commit{}
		commitJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(commitJson, commit); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return commit, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

//Repo Merge Section
// 
// GitHub  Docs: Repo: Merge - Perform a Merge
// Url: https://api.github.com/repos/:owner/:repo/merges?access_token=...
// Request Type: POST /repos/:owner/:repo/merges
// Access Token: REQUIRED
// 
func (github *GitHubClient) Merge(urlData map[string]string, postData map[string]string) (*Commit, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("There is missing data in urlData")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}
	if ok := github.AssertMapStrings([]string{"head", "base"}, urlData); !ok {
		return nil, errors.New("There is missing data in postData")
	}

	apiReader, err := github.CreateReader(postData)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/merges")
	res, err := github.Client.Post(apiUrl, "application/json", apiReader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	//SUCCESS
	if res.StatusCode == 201 {
		commit := &Commit{}
		commitJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(commitJson, commit); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return commit, nil
	}
	// NO NEED TO MERGE
	if res.StatusCode == 204 {
		github.getLimits(res)
		return nil, nil
	}
	// CONFLICTS ETC
	if res.StatusCode == 409 || res.StatusCode == 404 {
		msg := Message{}
		msgJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(msgJson, msg); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return nil, errors.New(msg.Message)
	}

	github.getLimits(res)
	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

//Status Section
// 
// GitHub  Docs: Repo: List Statuses for a specific SHA
// Url: https://api.github.com/repos/:owner/:repo/merges?access_token=...
// Request Type: GET /repos/:owner/:repo/statuses/:sha
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetCommitStatus(urlData map[string]string) (*Statuses, error) {
	if ok := github.AssertMapStrings([]string{"repo", "sha"}, urlData); !ok {
		return nil, errors.New("There is missing data in urlData. Both 'repo' and 'sha' are required.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/statuses/" + urlData["sha"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		statuses := &Statuses{}
		statusesJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(statusesJson, statuses); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return statuses, nil
	}

	github.getLimits(res)
	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub  Docs: Repo: Statuses - Create a Status
// Url: https://api.github.com/repos/:owner/:repo/statuses/:sha?access_token=...
// Request Type: POST /repos/:owner/:repo/statuses/:sha
// Access Token: REQUIRED
// 
func (github *GitHubClient) CreateStatus(urlData, postData map[string]string) (*Status, error) {
	if ok := github.AssertMapStrings([]string{"repo", "sha"}, urlData); !ok {
		return nil, errors.New("There is data missing from urlData. 'repo' and 'sha' are both required values")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}
	if ok := github.AssertMapString("state", postData); !ok { //Post Data Validation
		return nil, errors.New("There is data for the merge missing.")
	}

	reader, err := github.CreateReader(postData)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/statuses/" + urlData["sha"])
	res, err := github.Client.Post(apiUrl, "application/json", reader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		status := &Status{}
		statusJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(statusJson, status); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return status, nil
	}

	github.getLimits(res)
	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

//Comments Section 
// @TODO Create Comments Section - Very low priority

//Web Hooks Section
// 
// GitHub  Docs: Repo: Hook - List Repo Hooks
// Url: https://api.github.com/repos/:owner/:repo/hooks?access_token=...
// Request Type: GET /repos/:owner/:repo/hooks
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetRepoHooks(urlData map[string]string) (*Hooks, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("The url data is missing the 'repo' value.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/hooks")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		hooks := &Hooks{}
		hooksJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(hooksJson, hooks); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return hooks, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub  Docs: Repo: Hook - Get single hook
// Url: https://api.github.com/repos/:owner/:repo/hooks/:id?access_token=...
// Request Type: GET /repos/:owner/:repo/hooks/:id
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetHookById(urlData map[string]string) (*Hook, error) {
	if ok := github.AssertMapStrings([]string{"repo", "id"}, urlData); !ok {
		return nil, errors.New("There are required parts of the urlData missing. 'repo' and 'id' are both required strings.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/hooks/" + urlData["id"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		hook := &Hook{}
		hookJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(hookJson, hook); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return hook, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub  Docs: Repo: Hook - Edit a hook
// Url: https://api.github.com/repos/:owner/:repo/hooks/:id?access_token=...
// Request Type: PATCH /repos/:owner/:repo/hooks/:id
// Access Token: REQUIRED
// 
func (github *GitHubClient) CreateHook(urlData map[string]string, postData map[string]interface{}) (*Hook, error) {
	if ok := github.AssertMapValues([]string{"config", "name"}, postData); !ok {
		return nil, errors.New("There is missing data in the post data, 'name' and 'config' are required values")
	}
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("There is missing data in the url data. 'repo' is a required value.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	reader, err := github.CreateReader(postData)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/hooks")
	res, err := github.Client.Post(apiUrl, "application/json", reader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		hook := &Hook{}
		hookJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(hookJson, hook); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return hook, nil
	}

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// 
// GitHub  Docs: Repo: Hook - Edit a hook
// Url: https://api.github.com/repos/:owner/:repo/hooks/:id?access_token=...
// Request Type: PATCH /repos/:owner/:repo/hooks/:id
// Access Token: REQUIRED
// 
func (github *GitHubClient) EditHook(urlData map[string]string, postData map[string]interface{}) (*Hook, error) {
	if ok := github.AssertMapValues([]string{"config", "name"}, postData); !ok {
		return nil, errors.New("There is missing data in the post data, 'name' and 'config' are required values")
	}
	if ok := github.AssertMapStrings([]string{"repo", "id"}, urlData); !ok {
		return nil, errors.New("There is missing data in the post data, 'name' and 'config' are required values")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	reader, err := github.CreateReader(postData)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/hooks/" + urlData["id"])

	apiRequest, err := http.NewRequest("PATCH", apiUrl, reader)
	if err != nil {
		return nil, err
	}

	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		hook := &Hook{}
		hookJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(hookJson, hook); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return hook, nil
	}

	github.getLimits(res)
	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub  Docs: Repo: Hook - Test a hook
// Url: https://api.github.com/repos/:owner/:repo/hooks/:id/tests?access_token=...
// Request Type: POST /repos/:owner/:repo/hooks/:id/tests
// Access Token: REQUIRED
// 
func (github *GitHubClient) TestHook(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapStrings([]string{"repo", "id"}, urlData); !ok {
		return false, errors.New("There is data missing for the url. Both 'repo' and 'id aree required fields.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/hooks/" + urlData["id"])
	res, err := github.Client.Post(apiUrl, "text/plain", nil)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode == 204 {
		github.getLimits(res)
		return true, nil
	}

	github.getLimits(res)
	return false, nil
}

// 
// GitHub  Docs: Repo: Hook - PubSubHubbub
// Url: NONE
// Request Type: NONE
// Access Token: REQUIRED
// 
// NO PubSubHubbub
