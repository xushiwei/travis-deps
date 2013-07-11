package github

// GitHub API v3 Section - Users
// Allows you to manage the user related data as well as add keys, emails etc
//
//	##Users API
//		-  Get a single user
//		-  Get the authenticated user
//		-  Update the authenticated user
//		-  Get all users
//
//	##User Emails API
//		-  List email addresses for a user
//		-  Add email address(es)
//		-  Delete email address(es)
//
//	##User Followers API
//		-  List followers of a user
//		-  List users followed by another user
//		-  Check if you are following a user
//		-  Follow a user
//		-  Unfollow a user
//
//	##User Public Keys API
//		-  List public keys for a user
//		-  List your public keys
//		-  Get a single public key
//		-  Create a public key
//		-  Update a public key
//		-  Delete a public key

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type User struct {
	ID           int     `json:"id"`
	Login        string  `json:"login"`
	Url          string  `json:"url"`
	Name         Nstring `json:"name"`
	Company      Nstring `json:"company,omitempty"`
	HtmlURL      string  `json:"html_url,omitempty"`
	Blog         Nstring `json:"blog,omitempty"`
	Avatar       string  `json:"avatar_url"`
	GravatarID   string  `json:"gravatar_id"`
	Email        string  `json:"email,omitempty"`
	Location     Nstring `json:"location,omitempty"`
	Hireable     bool    `json:"hireable,omitempty"`
	CreatedAt    string  `json:"created_at,omitempty"`
	RepoUrl      string  `json:"repo_url,omitempty"`
	GistUrl      Nstring `json:"gist_url,omitempty"`
	EventUrl     string  `json:"events_url,omitempty"`
	StarUrl      string  `json:"starred_url,omitempty"`
	FollowersUrl string  `json:"followers_url,omitempty"`
	FollowingUrl string  `json:"following_url,omitempty"`
}

type Users []User

type Email string

type Emails []string

type GitKey struct {
	ID       int    `json:"id"`
	Key      string `json:"key"`
	Url      string `json:"url"`
	Verified bool   `json:"verfied"`
	Title    string `json:"title"`
}

type GitKeys []GitKey

type Follower GitUser

type Followers []Follower

// ******************
//// User Section  *
// ******************
// 
// GitHub Doc: "Get the authenticated user"
// Url: https://api.github.com/user?access_token=...
// Request Type: GET 
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetUser() (*User, error) {
	apiUrl := github.createUrl("/user")

	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		userJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		user := &User{}
		err = json.Unmarshal(userJson, user)
		if err != nil {
			return nil, err
		}

		github.getLimits(res)
		return user, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// ************************
// * START: Email Section  *
// **************************
// 
// GitHub Docs: List email addresses for a user - This endpoint is accessible with the user:email scope.
// Url: https://api.github.com/user/emails?access_token=...
// Request Type: GET /user/emails
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetEmails() (*Emails, error) {
	apiUrl := github.createUrl("/user/emails")

	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		emailsJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		emails := &Emails{}
		err = json.Unmarshal(emailsJson, emails)
		if err != nil {
			return nil, err
		}

		github.getLimits(res)
		return emails, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Docs: Add email address(es) - You can post a single email address or an array of addresses
// Url: https://api.github.com/user/emails?access_token=...
// Request Type: GET /user/keys
// Access Token: REQUIRED
// 
func (github *GitHubClient) AddEmail(email string) (*Emails, error) {
	apiUrl := github.createUrl("/user/emails")
	reader := strings.NewReader(`"` + email + `"`)

	res, err := github.Client.Post(apiUrl, "text/plain", reader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		emailsJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		emails := &Emails{}
		err = json.Unmarshal(emailsJson, emails)
		if err != nil {
			return nil, err
		}

		github.getLimits(res)
		return emails, nil
	}

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// 
// GitHub Docs: Delete email address(es) - You can post a single email address or an array of addresses
// Url: https://api.github.com/user/emails?access_token=...
// Request Type: DELETE /user/emails
// Access Token: REQUIRED
// 
func (github *GitHubClient) DeleteEmail(email string) error {
	apiUrl := github.createUrl("/user/emails")
	reader := strings.NewReader(`"` + email + `"`)

	apiRequest, err := http.NewRequest("DELETE", apiUrl, reader)
	if err != nil {
		return err
	}

	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 204 {
		github.getLimits(res)
		return nil
	}

	return errors.New("Didn't receive 204 status from Github: " + res.Status)
}

// ***********************
//  END: Email Section  *
// *************************
// ***********************
//  START: Key Section  *
// *************************
// 
// GitHub Docs: List public keys for a user - Lists the verified public keys for a user. This is accessible by anyone.
// Url: https://api.github.com/user/keys
// Request Type: GET /users/:user/keys
// Access Token: OPTIONAL
// 
// NOT NEEDED??
// 
// GitHub Docs: Get a single public key - Lists the current user’s keys. 
// 		Management of public keys via the API requires that you are 
//		authenticated through basic auth, or OAuth with the ‘user’ scope.
// Url: https://api.github.com/user/keys?access_token=...
// Request Type: GET /user/keys
// Access Token: REQUIRED
// 
func (github *GitHubClient) GetUserKeys() (*GitKeys, error) {
	apiUrl := github.createUrl("/user/keys")

	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		keysJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		keys := &GitKeys{}
		err = json.Unmarshal(keysJson, keys)
		if err != nil {
			return nil, err
		}

		github.getLimits(res)
		return keys, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Docs: Get a single public key
// Request Type: GET /user/keys/:id
// Access Token: REQUIRED
// Url: https://api.github.com/user/keys/:id?access_token=...
//
// id {int} - id of the key as noted in the struct of GitHubKey
// 
func (github *GitHubClient) GetKeyById(id int) (*GitKey, error) {
	apiUrl := github.createUrl("/user/keys/" + string(id))

	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		keyJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		key := &GitKey{}
		err = json.Unmarshal(keyJson, key)
		if err != nil {
			return nil, err
		}

		github.getLimits(res)
		return key, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Docs: Create a public key
// Request Type: POST /user/keys
// Access Token: REQUIRED
// Url: https://api.github.com/user/keys/:id?access_token=...
//
// id {int} - id of the key as noted in the struct of GitHubKey
// 
func (github *GitHubClient) CreateKey(key, title string) (*GitKey, error) {
	if key == "" {
		return nil, errors.New("No data for the key")
	}
	if title == "" {
		title = "CodeHub"
	}

	reader := strings.NewReader(`{ "key": "` + key + `", "title": "` + title + `" }`)
	apiUrl := github.createUrl("/user/keys")
	res, err := github.Client.Post(apiUrl, "application/json", reader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		keyJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		key := &GitKey{}
		err = json.Unmarshal(keyJson, key)
		if err != nil {
			return nil, err
		}

		github.getLimits(res)
		return key, nil
	}

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// 
// GitHub Docs: Update a public key
// Request Type: PATCH /user/keys/:id
// Access Token: REQUIRED
// Url: https://api.github.com/user/keys?access_token=...
//
// key {string} - the contents of the key - (Required)
// title {string} - the title of the key to help identify it - Defaults to CodeHub
// 

func (github *GitHubClient) UpdateKey(id int, key, title string) (*GitKey, error) {
	if id < 1 {
		return nil, errors.New("Ids cannot be less than 1")
	}
	if key == "" {
		return nil, errors.New("No data for the key")
	}
	if title == "" {
		title = "CodeHub"
	}

	reader := strings.NewReader(`{ "key": "` + key + `", "title": "` + title + `" }`)
	apiUrl := github.createUrl("/user/keys/" + string(id))
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
		keyJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		key := &GitKey{}
		err = json.Unmarshal(keyJson, key)
		if err != nil {
			return nil, err
		}

		github.getLimits(res)
		return key, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Docs: Update a public key
// Request Type: PATCH /user/keys/:id
// Access Token: REQUIRED
// Url: https://api.github.com/user/keys?access_token=...
//
// key {string} - the contents of the key - (Required)
// title {string} - the title of the key to help identify it - Defaults to "CodeHub"  
// 

func (github *GitHubClient) DeleteKey(id int) error {
	if id < 1 {
		return errors.New("Ids cannot be less than 1")
	}

	apiUrl := github.createUrl("/user/keys/" + string(id))
	apiRequest, err := http.NewRequest("DELETE", apiUrl, nil)
	if err != nil {
		return err
	}

	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 204 {
		github.getLimits(res)
		return nil
	}

	return errors.New("Didn't receive 204 status from Github: " + res.Status)
}

// 
// GitHub Docs: Get a single user
// Request Type: GET /user/followers
// Access Token: REQUIRED
// Url: https://api.github.com/user/followers
// 

func (github *GitHubClient) GetFollowers() (*Followers, error) {
	apiUrl := github.createUrl("/user/followers")

	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		followJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		followers := &Followers{}
		err = json.Unmarshal(followJson, followers)
		if err != nil {
			return nil, err
		}

		github.getLimits(res)
		return followers, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Docs: List users followed by another user
// Request Type: GET /user/following
// Access Token: REQUIRED
// Url: https://api.github.com/user/following
// 

func (github *GitHubClient) GetFollowing() (*Followers, error) {
	apiUrl := github.createUrl("/user/following")

	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		followJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		followers := &Followers{}
		err = json.Unmarshal(followJson, followers)
		if err != nil {
			return nil, err
		}

		github.getLimits(res)
		return followers, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Docs: Check if you are following a user
// Request Type: GET /user/following/:user
// Access Token: REQUIRED
// Url: https://api.github.com/user/following
// 

func (github *GitHubClient) AreFollowing(user string) (bool, error) {
	apiUrl := github.createUrl("/user/following/" + user)

	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode == 204 {
		return true, nil
	} else if res.StatusCode == 404 {
		return false, nil
	}

	return false, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Docs: Follow a user - Following a user requires the user to be logged in and authenticated with basic auth or OAuth with the user:follow scope.
// Request Type: PUT /user/following/:user
// Access Token: REQUIRED
// Url: https://api.github.com/user/following/:user
// 

func (github *GitHubClient) FollowUser(user string) (bool, error) {
	apiUrl := github.createUrl("/user/following/" + user)

	apiRequest, err := http.NewRequest("PUT", apiUrl, nil)

	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	if res.StatusCode == 204 {
		return true, nil
	}

	return false, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Docs: Unfollow a user - Unfollowing a user requres the user to be logged in and authenticated with basic auth or OAuth with the user:follow scope.
// Request Type: DELETE /user/following/:user
// Access Token: REQUIRED
// Url: https://api.github.com/user/following/:user
// 

func (github *GitHubClient) UnfollowUser(user string) (bool, error) {
	apiUrl := github.createUrl("/user/following/" + user)
	apiRequest, err := http.NewRequest("DELETE", apiUrl, nil)
	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode == 204 {
		return true, nil
	}

	return false, errors.New("Didn't receive 200 status from Github: " + res.Status)
}
