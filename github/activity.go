package github

// GitHub API v3 Section - Activity
// Activities used to generate user streams - will be great to add a personal social and github.
// 
//	## Events API
//		-  List public events
//		-  List repository events
//		-  List issue events for a repository
//		-  List public events for a network of repositories
//		-  List public events for an organization
//		-  List events that a user has received
//		-  List public events that a user has received
//		-  List events performed by a user
//		-  List public events performed by a user
//		-  List events for an organization
//

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type EventRepo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type EventOrg struct {
	Login       string `json:"login"`
	ID          int    `json:"id"`
	AvatarUrl   string `json:"avatar_url"`
	GravatarUrl string `json:"gravatar_url"`
	Url         string `json:"url"`
}

type Event struct {
	Type      string                 `json:"type"`
	Public    bool                   `json:"public"`
	Payload   map[string]interface{} `json:"payload"`
	Repo      EventRepo              `json:"repository"`
	Actor     GitUser                `json:"actor"`
	Org       EventOrg               `json:"org"`
	CreatedAt string                 `json:"created_at"`
	ID        string                 `json:"id"`
}

type NotifyRepo struct {
	ID          int     `json:"id"`
	Owner       GitUser `json:"owner"`
	Name        string  `json:"name"`
	FullName    string  `json:"full_name"`
	Description string  `json:"description"`
	Private     bool    `json:"private"`
	Fork        bool    `json:"fork"`
	Url         string  `json:"url"`
	HtmlUrl     string  `json:"html_url"`
}

type Notification struct {
	ID         int               `json:"id"`
	Repository NotifyRepo        `json:"repository"`
	Subject    map[string]string `json:"subject"`
	Reason     string            `json:"reason"`
	Unread     bool              `json:"unread"`
	Url        string            `json:"url"`
	UpdatedAt  string            `json:"updated_at"`
	LastReadAt string            `json:"last_read_at"`
}

type Subscription struct {
	Subscribed bool    `json:"subscribed"`
	Ignored    bool    `json:"ignored"`
	Reason     Nstring `json:"reason"`
	CreatedAt  string  `json:"created_at"`
	Url        string  `json:"url"`
	ThreadUrl  string  `json:"thread_url"`
}

// Events Section

// 
// GitHub Doc - Events: List public events
// Url: https://api.github.com/events?access_token=...
// Request Type: GET /events
// Access Token: PUBLIC
// 
func (github *GitHubClient) ListPublicEvents(page int) ([]Event, error) {
	if page > 10 || page < 1 {
		return nil, errors.New("The page number is not between 1 and 10.")
	}

	apiUrl := github.createUrl("/events?page=" + string(page))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		events := &[]Event{}
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

// 
// GitHub Doc - Events: List repository events
// Url: https://api.github.com/repos/:owner/:repo/events?access_token=...
// Request Type: GET /repos/:owner/:repo/events
// Access Token: REQUIRED
// 
func (github *GitHubClient) ListRepoEvents(ownerAndRepo string, page int) ([]Event, error) {
	if page > 10 || page < 1 {
		return nil, errors.New("The page number is not between 1 and 10.")
	}

	ownerAndRepo = strings.TrimSpace(ownerAndRepo)
	if len(ownerAndRepo) < 1 && strings.Index(ownerAndRepo, "/") < -1 {
		return nil, errors.New("Your ownerAndRepo string value is not valid")
	}

	apiUrl := github.createUrl("/repos/" + ownerAndRepo + "/events?page=" + string(page))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		events := &[]Event{}
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

// 
// GitHub Doc - Events: List issue events for a repository
// Url: https://api.github.com/repos/:owner/:repo/issues/events?access_token=...
// Request Type: GET /repos/:owner/:repo/issues/events
// Access Token: REQUIRED
// 
func (github *GitHubClient) ListIssuesEvents(ownerAndRepo string, page int) ([]Event, error) {
	if page > 10 || page < 1 {
		return nil, errors.New("The page number is not between 1 and 10.")
	}

	ownerAndRepo = strings.TrimSpace(ownerAndRepo)
	if len(ownerAndRepo) < 1 && strings.Index(ownerAndRepo, "/") < -1 {
		return nil, errors.New("Your ownerAndRepo string value is not valid")
	}

	apiUrl := github.createUrl("/repos/" + ownerAndRepo + "/issues/events?page=" + string(page))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		events := &[]Event{}
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

// 
// GitHub Doc - Events: List public events for a network of repositories
// Url: https://api.github.com/networks/:owner/:repo/events?access_token=...
// Request Type: GET /networks/:owner/:repo/events
// Access Token: REQUIRED
// 
func (github *GitHubClient) ListNetworkEvents(ownerAndRepo string, page int) ([]Event, error) {
	if page > 10 || page < 1 {
		return nil, errors.New("The page number is not between 1 and 10.")
	}

	ownerAndRepo = strings.TrimSpace(ownerAndRepo)
	if len(ownerAndRepo) < 1 && strings.Index(ownerAndRepo, "/") < -1 {
		return nil, errors.New("Your ownerAndRepo string value is not valid")
	}

	apiUrl := github.createUrl("/networks/" + ownerAndRepo + "/events?page=" + string(page))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		events := &[]Event{}
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

// 
// GitHub Doc - Events: List public events for an organization
// Url: https://api.github.com/orgs/:org/events?access_token=...
// Request Type: GET /orgs/:org/events
// Access Token: REQUIRED
// 
func (github *GitHubClient) ListOrgEvents(org string, page int) ([]Event, error) {
	if page > 10 || page < 1 {
		return nil, errors.New("The page number is not between 1 and 10.")
	}

	org = strings.TrimSpace(org)
	if len(org) < 1 {
		return nil, errors.New("Your org string value is not valid")
	}

	apiUrl := github.createUrl("/org/" + org + "/events?page=" + string(page))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		events := &[]Event{}
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

// 
// GitHub Doc - Events: List events that a user has received
// Url: https://api.github.com/users/:user/received_events?access_token=...
// Request Type: GET /users/:user/received_events
// Access Token: REQUIRED
// 
func (github *GitHubClient) RecievedUserEvents(user string, page int) ([]Event, error) {
	if page > 10 || page < 1 {
		return nil, errors.New("The page number is not between 1 and 10.")
	}

	user = strings.TrimSpace(user)
	if len(user) < 1 {
		return nil, errors.New("Your user string value is not valid")
	}

	apiUrl := github.createUrl("/users/" + user + "/received_events?page=" + string(page))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		events := &[]Event{}
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

// 
// GitHub Doc - Events: List public events that a user has received
// Url: https://api.github.com/users/:user/received_events/public?access_token=...
// Request Type: GET /users/:user/received_events/public
// Access Token: REQUIRED
// 
func (github *GitHubClient) PublicRecievedUserEvents(user string, page int) ([]Event, error) {
	if page > 10 || page < 1 {
		return nil, errors.New("The page number is not between 1 and 10.")
	}

	user = strings.TrimSpace(user)
	if len(user) < 1 {
		return nil, errors.New("Your user string value is not valid")
	}

	apiUrl := github.createUrl("/users/" + user + "/received_events/public?page=" + string(page))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		events := &[]Event{}
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

// 
// GitHub Doc - Events: List events performed by a user
// Url: https://api.github.com/users/:user/events?access_token=...
// Request Type: GET /users/:user/events
// Access Token: REQUIRED
// 
func (github *GitHubClient) PreformedUserEvents(user string, page int) ([]Event, error) {
	if page > 10 || page < 1 {
		return nil, errors.New("The page number is not between 1 and 10.")
	}

	user = strings.TrimSpace(user)
	if len(user) < 1 {
		return nil, errors.New("Your user string value is not valid")
	}

	apiUrl := github.createUrl("/users/" + user + "/events?page=" + string(page))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		events := &[]Event{}
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

// 
// GitHub Doc - Events: List public events performed by a user
// Url: https://api.github.com/users/:user/events/public?access_token=...
// Request Type: GET /users/:user/events/public
// Access Token: REQUIRED
// 
func (github *GitHubClient) PublicPreformedUserEvents(user string, page int) ([]Event, error) {
	if page > 10 || page < 1 {
		return nil, errors.New("The page number is not between 1 and 10.")
	}

	user = strings.TrimSpace(user)
	if len(user) < 1 {
		return nil, errors.New("Your user string value is not valid")
	}

	apiUrl := github.createUrl("/users/" + user + "/events/public?page=" + string(page))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		events := &[]Event{}
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

// 
// GitHub Doc - Events: List events for an organization
// Url: https://api.github.com/users/:user/events/orgs/:org?access_token=...
// Request Type: GET /users/:user/events/orgs/:org
// Access Token: REQUIRED
// 
func (github *GitHubClient) ListUserOrgEvents(user, org string, page int) ([]Event, error) {
	if page > 10 || page < 1 {
		return nil, errors.New("The page number is not between 1 and 10.")
	}

	user = strings.TrimSpace(user)
	org = strings.TrimSpace(org)
	if len(user) < 1 || len(org) < 1 {
		return nil, errors.New("Your user and/or org string value is not long enough")
	}

	apiUrl := github.createUrl("/users/" + user + "/events/orgs/" + org + "?page=" + string(page))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		events := &[]Event{}
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

// Notifications Section

// 
// GitHub Doc - Notifications: List your notifications
// Url: https://api.github.com/notifications?access_token=...
// Request Type: GET /notifications
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetNotifications(urlData map[string]string) ([]Notification, error) {
	apiUrl := github.createUrl("/notifications?" + github.UrlDataConvert(urlData))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		notify := &[]Notification{}
		notifyJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(notifyJson, notify); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*notify), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - Notifications: List your notifications in a repository
// Url: https://api.github.com/repos/:owner/:repo/notifications?access_token=...
// Request Type: GET /repos/:owner/:repo/notifications
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetRepoNotifications(urlData, getData map[string]string) ([]Notification, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("One or more fields are missing and/or do not have content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/notifications?" + github.UrlDataConvert(getData))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		notify := &[]Notification{}
		notifyJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(notifyJson, notify); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*notify), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - Notifications: Mark As Read
// Url: https://api.github.com/repos/:owner/:repo/notifications?access_token=...
// Request Type: PUT /repos/:owner/:repo/notifications
// Access Token: REQUIRED
// 
func (github *GitHubClient) MarkNotificationsRead(read bool, lastRead string) (bool, error) {
	un := "read"
	if !read {
		un = "unread"
	}

	apiUrl := github.createUrl("/notifications?" + un + "=true&last_read_at=" + strings.TrimSpace(url.QueryEscape(lastRead)))
	apiRequest, err := http.NewRequest("PUT", apiUrl, nil)
	if err != nil {
		return false, err
	}

	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode == 205 {
		github.getLimits(res)
		return true, nil
	}

	return false, errors.New("Didn't receive 205 status from Github: " + res.Status)
}

// 
// GitHub Doc - Notifications: View a single thread
// Url: https://api.github.com/notifications/threads/:id?access_token=...
// Request Type:GET /notifications/threads/:id
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetNotification(id string) (*Notification, error) {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return nil, errors.New("The id given does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/notifications/threads/" + id)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		notify := &Notification{}
		notifyJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(notifyJson, notify); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return notify, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - Notifications: Mark a thread as read
// Url: https://api.github.com/notifications/threads/:id?access_token=...
// Request Type: PATCH /notifications/threads/:id
// Access Token: REQUIRED
// 
func (github *GitHubClient) MarkThreadRead(read bool, id string) (bool, error) {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return false, errors.New("The id given does not contain any non-whitespace content")
	}

	un := "read"
	if !read {
		un = "unread"
	}

	apiUrl := github.createUrl("/notifications/threads/" + url.QueryEscape(id) + "?" + un + "=true")
	apiRequest, err := http.NewRequest("PATCH", apiUrl, nil)
	if err != nil {
		return false, err
	}

	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode == 205 {
		github.getLimits(res)
		return true, nil
	}

	return false, errors.New("Didn't receive 205 status from Github: " + res.Status)
}

// 
// GitHub Doc - Notifications: Get a Thread Subscription
// Url: https://api.github.com/notifications/threads/:id?access_token=...
// Request Type: GET /notifications/threads/:id/subscription
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetThreadSub(id string) (*Subscription, error) {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return nil, errors.New("The id given does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/notifications/threads/" + url.QueryEscape(id) + "/subscription")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		sub := &Subscription{}
		subJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(subJson, sub); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return sub, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - Notifications: Set a Thread Subscription
// Url: https://api.github.com/notifications/threads/1/subscription?access_token=...
// Request Type: PUT /notifications/threads/1/subscription
// Access Token: REQUIRED
// 
func (github *GitHubClient) SubToThread(id string, subed, ignored bool) (*Subscription, error) {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return nil, errors.New("The id given does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/notifications/threads/" + url.QueryEscape(id) + "/subscription?subscribed=" + strconv.FormatBool(subed) + "&ignored=" + strconv.FormatBool(ignored))
	apiRequest, err := http.NewRequest("PUT", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		sub := &Subscription{}
		subJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(subJson, sub); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return sub, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - Notifications: Delete a Thread Subscription
// Url: https://api.github.com/notifications/threads/1/subscription?access_token=...
// Request Type: DELETE /notifications/threads/1/subscription
// Access Token: REQUIRED
// 
func (github *GitHubClient) DeleteThread(id string) (bool, error) {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return false, errors.New("The id given does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/notifications/threads/" + url.QueryEscape(id) + "/subscription")
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

// Starred Section

// 
// GitHub Doc - Starred: List Stargazers
// Url: https://api.github.com/repos/:owner/:repo/stargazers?access_token=...
// Request Type: GET /repos/:owner/:repo/stargazers
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetStargazers(urlData map[string]string, page int) ([]GitUser, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("One or more fields are missing and/or do not have content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	if page < 1 {
		page = 1
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/stargazers?page=" + string(page))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		user := &[]GitUser{}
		userJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(userJson, user); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*user), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// Starred Section

// 
// GitHub Doc - Starred: List repositories being starred
// Url: https://api.github.com/user/starred?access_token=...
// Request Type: GET /user/starred
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetStarredRepos(getData map[string]string) (*Repos, error) {
	urlStr := github.UrlDataConvert(getData)

	apiUrl := github.createUrl("/user/starred?" + urlStr)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		repos := &Repos{}
		reposJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(reposJson, repos); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return repos, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - Starred: List repositories being starred
// Url: https://api.github.com/user/starred/:owner/:repo?access_token=...
// Request Type: GET /user/starred/:owner/:repo
// Access Token: REQUIRED
// 
func (github *GitHubClient) AreStarringRepo(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return false, errors.New("One or more fields are missing and/or do not have content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/user/starred/" + strings.TrimSpace(urlData["owner"]) + "/" + strings.TrimSpace(urlData["repo"]))
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

	return false, errors.New("Didn't receive 204/404 status from Github: " + res.Status)
}

// 
// GitHub Doc - Starred: Star a repository - Requires for the user to be authenticated.
// Url: https://api.github.com/user/starred/:owner/:repo?access_token=...
// Request Type: PUT /user/starred/:owner/:repo
// Access Token: REQUIRED
// 
func (github *GitHubClient) StarRepo(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return false, errors.New("One or more fields are missing and/or do not have content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/user/starred/" + strings.TrimSpace(urlData["owner"]) + "/" + strings.TrimSpace(urlData["repo"]))
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
// GitHub Doc - Starred: Star a repository - Requires for the user to be authenticated.
// Url: https://api.github.com/user/starred/:owner/:repo?access_token=...
// Request Type: PUT /user/starred/:owner/:repo
// Access Token: REQUIRED
// 
func (github *GitHubClient) UnstarRepo(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return false, errors.New("One or more fields are missing and/or do not have content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/user/starred/" + strings.TrimSpace(urlData["owner"]) + "/" + strings.TrimSpace(urlData["repo"]))
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

// Watcher Section
// 
// GitHub Doc - Watchers: List watchers
// Url: https://api.github.com/repos/:owner/:repo/subscribers?access_token=...
// Request Type: GET /repos/:owner/:repo/subscribers
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetWatchers(urlData map[string]string, page int) ([]GitUser, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("One or more fields are missing and/or do not have content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	if page < 1 {
		page = 1
	}

	apiUrl := github.createUrl("/repos/" + strings.TrimSpace(urlData["owner"]) + "/" + strings.TrimSpace(urlData["repo"]) + "/subscribers?page=" + string(page))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		user := &[]GitUser{}
		userJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(userJson, user); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*user), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - Watchers: List repositories being watched
// Url: https://api.github.com/user/subscriptions?access_token=...
// Request Type: GET /user/subscriptions
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetWatchedRepos(page int) (*Repos, error) {
	apiUrl := github.createUrl("/user/subscriptions?page=" + string(page))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		repos := &Repos{}
		reposJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(reposJson, repos); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return repos, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - Watchers: Get a Repository Subscription
// Url: https://api.github.com/repos/:owner/:repo/subscription?access_token=...
// Request Type: GET /repos/:owner/:repo/subscription
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetRepoWatch(urlData map[string]string) (*Subscription, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("One or more fields are missing and/or do not have content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + strings.TrimSpace(urlData["owner"]) + "/" + strings.TrimSpace(urlData["repo"]) + "/subscription")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		sub := &Subscription{}
		subJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(subJson, sub); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return sub, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc - Watchers: Set a Repository Subscription
// Url: https://api.github.com/repos/:owner/:repo/subscription?access_token=...
// Request Type: PUT /repos/:owner/:repo/subscription
// Access Token: REQUIRED
// 
func (github *GitHubClient) WatchRepo(urlData map[string]string, subed, ignored bool) (*Subscription, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("One or more fields are missing and/or do not have content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + strings.TrimSpace(urlData["owner"]) + "/" + strings.TrimSpace(urlData["repo"]) + "/subscription?subscribed=" + strconv.FormatBool(subed) + "&ignored=" + strconv.FormatBool(ignored))
	apiRequest, err := http.NewRequest("PUT", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		sub := &Subscription{}
		subJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(subJson, sub); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return sub, nil
	}

	return nil, errors.New("Didn't receive 204 status from Github: " + res.Status)
}

// 
// GitHub Doc - Watchers: Set a Repository Subscription
// Url: https://api.github.com/repos/:owner/:repo/subscription?access_token=...
// Request Type: PUT /repos/:owner/:repo/subscription
// Access Token: REQUIRED
// 
func (github *GitHubClient) UnwatchRepo(urlData map[string]string, subed, ignored bool) (bool, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return false, errors.New("One or more fields are missing and/or do not have content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + strings.TrimSpace(urlData["owner"]) + "/" + strings.TrimSpace(urlData["repo"]) + "/subscription")
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
