package github

// Issues Section of the GitHub
// Includes the comments since I think they are much more useful in Gists, but rare for commits.
// GitHub API v3 Orgs Section
//
//	##  Issues API
//		-  List issues X3
//		-  List issues for a repository
//		-  Get a single issue
//		-  Create an issue
//		-  Edit an issue
//		-  Custom media types
//
//	##  Assignees API
//		-  List assignees
//		-  Check assignee
//
//	##  Issue Comments API
//		-  List comments on an issue
//		-  List comments in a repository
//		-  Get a single comment
//		-  Create a comment
//		-  Edit a comment
//		-  Delete a comment
//
//	##  Issue Events API
//		-  List events for an issue
//		-  List events for a repository
//		-  Get a single event
//
//	##  Labels API
//		-  List all labels for this repository
//		-  Get a single label
//		-  Create a label
//		-  Update a label
//		-  Delete a label
//		-  List labels on an issue
//		-  Add labels to an issue
//		-  Remove a label from an issue
//		-  Replace all labels for an issue
//		-  Remove all labels from an issue
//		-  Get labels for every issue in a milestone
//
//	##  Milestones API
//		-  List milestones for a repository
//		-  Get a single milestone
//		-  Create a milestone
//		-  Update a milestone
//		-  Delete a milestone

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type Milestone struct {
	Url          string  `json:"url"`
	Number       int     `json:"number"`
	State        Nstring `json:"state,omitempty"`
	Title        string  `json:"title"`
	Description  Nstring `json:"description,omitempty"`
	Creator      GitUser `json:"creator,omitempty"`
	OpenIssues   int     `json:"open_issues,omitempty"`
	ClosedIssues int     `json:"closed_issues,omitempty"`
	CreatedAt    Nstring `json:"created_at,omitempty"`
	DueOn        Nstring `json:"due_on,omitempty"`
}

type CreateMilestone struct {
	Title       string `json:"title,omitempty"`
	State       string `json:"state,omitempty"`
	Description string `json:"description,omitempty"`
	DueOn       string `json:"due_on,omitempty"`
}

func (n *Milestone) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*Milestone)(n))
}

type IssueLinks struct {
	Self        map[string]Nstring `json:"self,omitempty"`
	Html        map[string]Nstring `json:"html,omitempty"`
	PullRequest map[string]Nstring `json:"pull_request,omitempty"`
}

type Comment struct {
	ID        int        `json:"id"`
	Url       string     `json:"url,omitempty"`
	Body      Nstring    `json:"body,omitempty"`
	User      GitUser    `json:"user,omitempty"`
	CreatedAt string     `json:"created_at,omitempty"`
	UpdatedAt Nstring    `json:"updated_at,omitempty"`
	Links     IssueLinks `json:"_links,omitempty"`
}

func (n *Comment) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*Comment)(n))
}

type IssueLabel struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Color string `json:"color"`
}

type Issue struct {
	Url         string              `json:"url"`
	HtmlUrl     string              `json:"html_url"`
	Number      int                 `json:"number"`
	State       Nstring             `json:"state,omitempty"`
	Title       string              `json:"title"`
	Body        Nstring             `json:"body,omitempty"`
	User        GitUser             `json:"user"`
	Labels      []map[string]string `json:"labels,omitempty"`
	Assignee    GitUser             `json:"assignee,omitempty"`
	Milestone   Milestone           `json:"milestone,omitempty"`
	Comments    int                 `json:"comments,omitempty"`
	PullRequest map[string]Nstring  `json:"pull_request,omitempty"`
	ClosedAt    Nstring             `json:"closed_at,omitempty"`
	CreatedAt   string              `json:"created_at,omitempty"`
	UpdatedAt   Nstring             `json:"updated_at,omitempty"`
}

type CreateIssue struct {
	Title     string   `json:"title,omitempty"`
	Body      string   `json:"body,omitempty"`
	Assignee  string   `json:"assignee,omitempty"`
	Milestone int      `json:"milestone,omitempty"`
	Labels    []string `json:"labels,omitempty"`
}

type IssueEvent struct {
	Event     string  `json:"event"`
	Url       string  `json:"url"`
	Actor     GitUser `json:"actor"`
	CommitId  string  `json:"commit_id"`
	CreatedAt string  `json:"created_At"`
}

// 
// GitHub Doc - Issues: List issues
// Url: https://api.github.com/issues?access_token=...
// Request Type: GET /issues
// Access Token: REQUIRED
// 
// List all issues across all the authenticated user’s visible repositories including owned repositories, 
// member repositories, and organization repositories:
func (github *GitHubClient) ListAllIssues(getData map[string]string) ([]Issue, error) {
	if ok := github.AssertMapString("filter", getData); !ok {
		return nil, errors.New(getData["filter"] + `The getData["filter"] value is either empty or doesn't contain any non-whitespace content`)
	}

	urlStr := github.UrlDataConvert(getData)
	apiUrl := github.createUrl("/issues?" + urlStr)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		issues := &[]Issue{}
		issuesJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(issuesJson, issues); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*issues), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - Issues: List issues
// Url: https://api.github.com/user/issues?access_token=...
// Request Type: GET /user/issues
// Access Token: REQUIRED
// 
//  List all issues across owned and member repositories for the authenticated user:  
func (github *GitHubClient) ListUserIssues(getData map[string]string) ([]Issue, error) {
	if ok := github.AssertMapString("filter", getData); !ok {
		return nil, errors.New(`The getData["filter"] value is either empty or doesn't contain any non-whitespace content`)
	}

	urlStr := github.UrlDataConvert(getData)
	apiUrl := github.createUrl("/user/issues?" + urlStr)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		issues := &[]Issue{}
		issuesJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(issuesJson, issues); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*issues), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - Issues: List issues
// Url: https://api.github.com/orgs/:org/issues?access_token=...
// Request Type: GET /orgs/:org/issues
// Access Token: REQUIRED
// 
//  List all issues for a given organization for the authenticated user: 
func (github *GitHubClient) ListOrgIssues(org string, getData map[string]string) ([]Issue, error) {
	if ok := github.AssertMapString("filter", getData); !ok {
		return nil, errors.New(`The getData["filter"] value is either empty or doesn't contain any non-whitespace content`)
	}

	org = strings.TrimSpace(org)
	if len(org) == 0 {
		return nil, errors.New("The org data given does not contain any non-whitespace content")
	}

	urlStr := github.UrlDataConvert(getData)
	apiUrl := github.createUrl("/orgs/" + org + "/issues?" + urlStr)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		issues := &[]Issue{}
		issuesJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(issuesJson, issues); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*issues), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - List issues for a repository
// Url: https://api.github.com/repos/:owner/:repo/issues?access_token=...
// Request Type: GET /repos/:owner/:repo/issues
// Access Token: REQUIRED
// 
//  List all issues for a given organization for the authenticated user: 
func (github *GitHubClient) ListRepoIssues(urlData, getData map[string]string) ([]Issue, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("filter", getData); !ok {
		return nil, errors.New(`The getData["filter"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	urlStr := github.UrlDataConvert(getData)
	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/issues?" + urlStr)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		issues := &[]Issue{}
		issuesJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(issuesJson, issues); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*issues), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - Get a single issue
// Url: https://api.github.com/repos/:owner/:repo/issues/:number?access_token=...
// Request Type: GET /repos/:owner/:repo/issues/:number
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetIssue(urlData map[string]string) (*Issue, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value and/or urlData["number"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/issues/" + urlData["number"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		issue := &Issue{}
		issueJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(issueJson, issue); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return issue, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - Create an issue
// Url: https://api.github.com/repos/:owner/:repo/issues?access_token=...
// Request Type: POST /repos/:owner/:repo/issues
// Access Token: REQUIRED
// 
func (github *GitHubClient) CreateIssue(urlData map[string]string, issueData *CreateIssue) (*Issue, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiReader, err := github.CreateReader(issueData)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/issues")
	res, err := github.Client.Post(apiUrl, "application/json", apiReader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		issue := &Issue{}
		issueJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(issueJson, issue); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return issue, nil
	}

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// 
// GitHub Doc - Edit an issue
// Url: https://api.github.com/repos/:owner/:repo/issues/:number?access_token=...
// Request Type: PATCH /repos/:owner/:repo/issues/:number
// Access Token: REQUIRED
// 
func (github *GitHubClient) EditIssue(urlData map[string]string, issueData *CreateIssue) (*Issue, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value and/or urlData["number"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiReader, err := github.CreateReader(issueData)
	if err != nil {
		return nil, err
	}
	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/issues/" + urlData["number"])
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
		issue := &Issue{}
		issueJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(issueJson, issue); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return issue, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

//  Issues - Assignee Section
// 
// GitHub Doc - Issues: List assignees
// Url: https://api.github.com/repos/:owner/:repo/assignees?access_token=...
// Request Type: GET /repos/:owner/:repo/assignees
// Access Token: REQUIRED
// 
func (github *GitHubClient) ListAssignees(urlData map[string]string) ([]GitUser, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/assignees")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		users := &[]GitUser{}
		usersJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(usersJson, users); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*users), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// GitHub Doc - Issues: Check assignee
// Url: https://api.github.com/repos/:owner/:repo/assignees/:assignee?access_token=...
// Request Type: GET /repos/:owner/:repo/assignees/:assignee
// Access Token: REQUIRED
// 
func (github *GitHubClient) CheckAssignees(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapStrings([]string{"repo", "assignee"}, urlData); !ok {
		return false, errors.New(`The urlData["repo"] value and/or urlData["assignee"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/assignees/" + urlData["assignee"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode == 204 {
		github.getLimits(res)
		return true, nil
	} else if res.StatusCode == 404 {
		github.getLimits(res)
		return false, nil
	}

	return false, errors.New("Didn't receive 204/404 status from Github: " + res.Status)
}

//  Issues - Events Section
// 
// GitHub Doc - Issues: List events for an issue
// Url: https://api.github.com/repos/:owner/:repo/issues/:issue_number/events?access_token=...
// Request Type: GET /repos/:owner/:repo/issues/:issue_number/events
// Access Token: REQUIRED
// 
func (github *GitHubClient) ListIssueEvents(urlData map[string]string, page int) ([]IssueEvent, error) {
	if ok := github.AssertMapStrings([]string{"repo", "issueNumber"}, urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value and/org urlData["issueNumber"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/issues/" + urlData["issueNumber"] + "/events?page=" + string(page))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		events := &[]IssueEvent{}
		eventsJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(eventsJson, events); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*events), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// GitHub Doc - Issues: List events for a repository
// Url: https://api.github.com/repos/:owner/:repo/issues/events?access_token=...
// Request Type: GET /repos/:owner/:repo/issues/events
// Access Token: REQUIRED
// 
func (github *GitHubClient) ListRepoIssueEvents(urlData map[string]string, page int) ([]IssueEvent, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/issues/events?page=" + string(page))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		events := &[]IssueEvent{}
		eventsJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(eventsJson, events); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*events), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// GitHub Doc - Issues: Get a single event
// Url: https://api.github.com/repos/:owner/:repo/issues/events/:id?access_token=...
// Request Type: GET /repos/:owner/:repo/issues/events/:id
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetIssueEvent(urlData map[string]string) (*IssueEvent, error) {
	if ok := github.AssertMapStrings([]string{"repo", "id"}, urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/issues/events/" + urlData["id"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		event := &IssueEvent{}
		eventJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(eventJson, event); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return event, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// GitHub Doc - Issues: List milestones for a repository
// Url: https://api.github.com/repos/:owner/:repo/milestones?access_token=...
// Request Type: GET /repos/:owner/:repo/milestones
// Access Token: REQUIRED
// getData map[string]string -> included page as a string
// 
func (github *GitHubClient) ListRepoMilestones(urlData, getData map[string]string) ([]Milestone, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	urlStr := github.UrlDataConvert(getData)

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/milestones?" + urlStr)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		milestones := &[]Milestone{}
		milestonesJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(milestonesJson, milestones); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*milestones), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// GitHub Doc - Issues: Get a single milestone
// Url: https://api.github.com/repos/:owner/:repo/milestones/:number?access_token=...
// Request Type: GET /repos/:owner/:repo/milestones/:number
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetRepoMilestone(urlData map[string]string) (*Milestone, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/milestones/" + urlData["number"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		milestone := &Milestone{}
		milestoneJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(milestoneJson, milestone); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return milestone, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// GitHub Doc - Issues: Create a milestone
// Url: https://api.github.com/repos/:owner/:repo/milestones?access_token=...
// Request Type: POST /repos/:owner/:repo/milestones
// Access Token: REQUIRED
// 
func (github *GitHubClient) CreateMilestone(urlData, msData map[string]string) (*Milestone, error) {
	if len(strings.TrimSpace(msData["title"])) == 0 {
		return nil, errors.New(`The msData["title"] value doesn't containt any non-whitespace content`)
	}
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/milestones")
	apiReader, err := github.CreateReader(msData)
	if err != nil {
		return nil, err
	}

	res, err := github.Client.Post(apiUrl, "application/json", apiReader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		milestone := &Milestone{}
		milestoneJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(milestoneJson, milestone); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return milestone, nil
	}

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// GitHub Doc - Issues: Update a milestone
// Url: https://api.github.com/repos/:owner/:repo/milestones/:number?access_token=...
// Request Type: PATCH /repos/:owner/:repo/milestones/:number
// Access Token: REQUIRED
// 
func (github *GitHubClient) UpdateMilestone(urlData, msData map[string]string) (*Milestone, error) {
	if len(strings.TrimSpace(msData["title"])) == 0 {
		return nil, errors.New(`The msData["title"] value doesn't containt any non-whitespace content`)
	}
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/milestones/" + urlData["number"])
	apiReader, err := github.CreateReader(msData)
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
		milestone := &Milestone{}
		milestoneJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(milestoneJson, milestone); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return milestone, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// GitHub Doc - Issues: Update a milestone
// Url: https://api.github.com/repos/:owner/:repo/milestones/:number?access_token=...
// Request Type: PATCH /repos/:owner/:repo/milestones/:number
// Access Token: REQUIRED
// 
func (github *GitHubClient) DeleteMilestone(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return false, errors.New(`The urlData["repo"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/milestones/" + urlData["number"])
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

// Issues - Comments Section
// 
// GitHub Doc - Issues: List comments on an issue
// Url: https://api.github.com/repos/:owner/:repo/issues/:number/comments?access_token=...
// Request Type: GET /repos/:owner/:repo/issues/:number/comments
// Access Token: REQUIRED
// 
func (github *GitHubClient) ListIssueComments(urlData map[string]string) ([]Comment, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value and/org urlData["number"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData[""] + "/" + urlData["repo"] + "/issues/" + urlData["number"] + "/comments")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		comments := &[]Comment{}
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

// GitHub Doc - Issues: List comments in a repository
// Url: https://api.github.com/repos/:owner/:repo/issues/comments?access_token=...
// Request Type: GET /repos/:owner/:repo/issues/comments
// Access Token: REQUIRED
// 
func (github *GitHubClient) ListRepoIssueComments(urlData, getData map[string]string) ([]Comment, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value and/org urlData["number"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	urlStr := github.UrlDataConvert(getData)
	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/issues/comments?" + urlStr)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		comments := &[]Comment{}
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

// GitHub Doc - Issues: Get a single comment
// Url: https://api.github.com/repos/:owner/:repo/issues/comments/:id?access_token=...
// Request Type: GET /repos/:owner/:repo/issues/comments/:id
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetIssueComment(urlData map[string]string) (*Comment, error) {
	if ok := github.AssertMapStrings([]string{"repo", "id"}, urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value and/org urlData["id"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/issues/comments/" + urlData["id"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		comment := &Comment{}
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

// GitHub Doc - Issues: Create a comment
// Url: https://api.github.com/repos/:owner/:repo/issues/:number/comments?access_token=...
// Request Type: POST /repos/:owner/:repo/issues/:number/comments
// Access Token: REQUIRED
// 
func (github *GitHubClient) CreateIssueComment(urlData map[string]string, commentBody string) (*Comment, error) {
	commentBody = strings.TrimSpace(commentBody)
	if len(commentBody) == 0 {
		return nil, errors.New("The comment body does not contain any non-whitespace content.")
	}

	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value and/org urlData["number"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	jsonText := `{ "body": "` + commentBody + `" }`
	apiReader := strings.NewReader(jsonText)
	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/issues/" + urlData["number"] + "/comments")
	res, err := github.Client.Post(apiUrl, "application/json", apiReader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		comment := &Comment{}
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

// GitHub Doc - Issues: Edit a comment
// Url: https://api.github.com/repos/:owner/:repo/issues/comments/:id?access_token=...
// Request Type: PATCH /repos/:owner/:repo/issues/comments/:id
// Access Token: REQUIRED
// 
func (github *GitHubClient) EditIssueComment(urlData map[string]string, commentBody string) (*Comment, error) {
	commentBody = strings.TrimSpace(commentBody)
	if len(commentBody) == 0 {
		return nil, errors.New("The comment body does not contain any non-whitespace content.")
	}

	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value and/org urlData["number"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	jsonText := `{ "body": "` + commentBody + `" }`
	apiReader := strings.NewReader(jsonText)
	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/issues/comments/" + urlData["id"])
	apiRequest, err := http.NewRequest("PATCH", apiUrl, apiReader)
	if err != nil {
		return nil, err
	}

	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		comment := &Comment{}
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

// GitHub Doc - Issues: Edit a comment
// Url: https://api.github.com/repos/:owner/:repo/issues/comments/:id?access_token=...
// Request Type: PATCH /repos/:owner/:repo/issues/comments/:id
// Access Token: REQUIRED
// 
func (github *GitHubClient) DeleteIssueComment(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return false, errors.New(`The urlData["repo"] value and/org urlData["number"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/issues/comments/" + urlData["id"])
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

	return false, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// Issues - Labels Section  
// 
// GitHub Doc - Issues: List comments in a repository
// Url: https://api.github.com/repos/:owner/:repo/labels?access_token=...
// Request Type: GET /repos/:owner/:repo/labels
// Access Token: REQUIRED
// 
func (github *GitHubClient) ListRepoLabels(urlData map[string]string) ([]IssueLabel, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/labels")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		labels := &[]IssueLabel{}
		labelsJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(labelsJson, labels); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*labels), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// GitHub Doc - Issues: List comments in a repository
// Url: https://api.github.com/repos/:owner/:repo/labels?access_token=...
// Request Type: GET /repos/:owner/:repo/labels
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetRepoLabel(urlData map[string]string) (*IssueLabel, error) {
	if ok := github.AssertMapStrings([]string{"repo", "name"}, urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/labels/" + urlData["name"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		label := &IssueLabel{}
		labelJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(labelJson, label); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return label, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// GitHub Doc - Issues:Create a label
// Url: https://api.github.com/repos/:owner/:repo/labels?access_token=...
// Request Type: POST /repos/:owner/:repo/labels
// Access Token: REQUIRED
// 
func (github *GitHubClient) CreateRepoLabel(urlData, labelData map[string]string) (*IssueLabel, error) {
	if ok := github.AssertMapStrings([]string{"repo", "name"}, urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiReader, err := github.CreateReader(labelData)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/labels")
	res, err := github.Client.Post(apiUrl, "application/json", apiReader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		label := &IssueLabel{}
		labelJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(labelJson, label); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return label, nil
	}

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// GitHub Doc - Issues: Create a label
// Url: https://api.github.com/repos/:owner/:repo/labels?access_token=...
// Request Type: POST /repos/:owner/:repo/labels
// Access Token: REQUIRED
// 
func (github *GitHubClient) UpdateRepoLabel(urlData, labelData map[string]string) (*IssueLabel, error) {
	if ok := github.AssertMapStrings([]string{"repo", "name"}, urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiReader, err := github.CreateReader(labelData)
	if err != nil {
		return nil, err
	}
	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/labels/" + urlData["name"])
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
		label := &IssueLabel{}
		labelJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(labelJson, label); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return label, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// GitHub Doc - Issues: Delete a label
// Url: https://api.github.com/repos/:owner/:repo/labels?access_token=...
// Request Type: DELETE /repos/:owner/:repo/labels/:name
// Access Token: REQUIRED
// 
func (github *GitHubClient) DeleteRepoLabel(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapStrings([]string{"repo", "name"}, urlData); !ok {
		return false, errors.New(`The urlData["repo"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/labels/" + urlData["name"])
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

	return false, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// GitHub Doc - Issues:List labels on an issue
// Url: https://api.github.com/repos/:owner/:repo/labels?access_token=...
// Request Type: GET /repos/:owner/:repo/issues/:number/labels
// Access Token: REQUIRED
// 
func (github *GitHubClient) ListIssueLabels(urlData map[string]string) ([]IssueLabel, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/issues/" + urlData["number"] + "/labels")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		labels := &[]IssueLabel{}
		labelsJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(labelsJson, labels); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*labels), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// GitHub Doc - Issues: Add labels to an issue
// Url: https://api.github.comrepos/:owner/:repo/issues/:number/labels?access_token=...
// Request Type: POST /repos/:owner/:repo/issues/:number/labels
// Access Token: REQUIRED
// 
func (github *GitHubClient) CreateIssueLabel(urlData, labelData map[string]string, labels []string) ([]IssueLabel, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New(`The urlData["repo"] value is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiReader, err := github.CreateReader(labelData)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/issues/" + urlData["number"] + "/labels")
	res, err := github.Client.Post(apiUrl, "application/json", apiReader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		labels := &[]IssueLabel{}
		labelsJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(labelsJson, labels); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*labels), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// GitHub Doc - Issues: Remove a label from an issue
// Url: https://api.github.comrepos/:owner/:repo/issues/:number/labels?access_token=...
// Request Type: POST /repos/:owner/:repo/issues/:number/labels
// Access Token: REQUIRED
// 
func (github *GitHubClient) RemoveIssueLabel(urlData, labelData map[string]string) (bool, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number", "name"}, urlData); !ok {
		return false, errors.New(`The urlData -> repo, number and name values is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/issues/" + urlData["number"] + "/labels/" + urlData["name"])
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

// GitHub Doc - Issues: Replace all labels for an issue
// Url: https://api.github.comrepos/repos/:owner/:repo/issues/:number/labelss?access_token=...
// Request Type: PUT /repos/:owner/:repo/issues/:number/labels
// Access Token: REQUIRED
// 
func (github *GitHubClient) ReplaceeIssueLabels(urlData map[string]string, labels []string) ([]IssueLabel, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New(`The urlData -> repo and/or number values is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiReader, err := github.CreateReader(labels)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/issues/" + urlData["number"] + "/labels")
	apiRequest, err := http.NewRequest("PUT", apiUrl, apiReader)
	if err != nil {
		return nil, err
	}

	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		labels := &[]IssueLabel{}
		labelsJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(labelsJson, labels); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*labels), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// GitHub Doc - Issues: Get labels for every issue in a milestone
// Url: https://api.github.comrepos/:owner/:repo/issues/:number/labels?access_token=...
// Request Type: GET /repos/:owner/:repo/milestones/:number/labels
// Access Token: REQUIRED
// 
func (github *GitHubClient) RemoveIssueLabels(urlData map[string]string, labels []string) ([]IssueLabel, error) {
	if ok := github.AssertMapStrings([]string{"repo", "number"}, urlData); !ok {
		return nil, errors.New(`The urlData -> repo and/or number values is either empty or doesn't contain any non-whitespace content`)
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/issues/" + urlData["number"] + "/labels")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		labels := &[]IssueLabel{}
		labelsJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(labelsJson, labels); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*labels), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}
