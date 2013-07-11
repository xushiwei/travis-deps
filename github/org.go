package github

// Org Section of the GitHub
// Includes the comments since I think they are much more useful in Gists, but rare for commits.
// GitHub API v3 Orgs Section
//
//	## Orgs API
//	-  List User Organizations
//	-  Get an Organization
//	-  Edit an Organization
//
//	Org Members API
//	-  Members list
//	-  Check membership
//	-  Add a member
//	-  Remove a member
//	-  Public members list
//	-  Check public membership
//	-  Publicize a user’s membership
//	-  Conceal a user’s membership
//
//	Org Teams API
//	-  List teams
//	-  Get team
//	-  Create team
//	-  Edit team
//	-  Delete team
//	-  List team members
//	-  Get team member
//	-  Add team member
//	-  Remove team member
//	-  List team repos
//	-  Get team repo
//	-  Add team repo
//	-  Remove team repo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type OrgPlan struct {
	Name         string `json:"name,omitempty"`
	Space        int    `json:"space,omitempty"`
	PrivateRepos int    `json:"private_repos,omitempty"`
}

type Org struct {
	Login             string  `json:"login"`
	ID                int     `json:"id"`
	Url               string  `json:"url"`
	AvatarUrl         Nstring `json:"avatar_url"`
	Name              string  `json:"name,omitempty"`
	Company           string  `json:"company,omitempty"`
	Blog              Nstring `json:"blog,omitempty"`
	Location          Nstring `json:"location,omitempty"`
	Email             string  `json:"email,omitempty"`
	PublicRepos       int     `json:"public_repos,omitempty"`
	PublicGists       int     `json:"public_gists,omitempty"`
	Followers         int     `json:"followers,omitempty"`
	Following         int     `json:"following,omitempty"`
	HtmlUrl           string  `json:"html_url,omitempty"`
	CreatedAt         string  `json:"created_at,omitempty"`
	Type              string  `json:"type,omitempty"`
	TotalPriavteRepos int     `json:"total_private_repos,omitempty"`
	OwnedPrivateRepos int     `json:"owned_private_repos,omitempty"`
	PrivateGists      int     `json:"private_gists,omitempty"`
	DiskUsage         int     `json:"disk_usage,omitempty"`
	Collaborators     int     `json:"collaborators,omitempty"`
	BillingEmail      Nstring `json:"billing_email,omitempty"`
	Plan              OrgPlan `json:"plan,omitempty"`
}

type Team struct {
	Url          string `json:"url"`
	Name         string `json:"name"`
	ID           int    `json:"id"`
	Permission   string `json:"permission,omitempty"`
	MembersCount int    `json:"members_count,omitempty"`
	ReposCount   int    `json:"repos_count,omitempty"`
}
type Teams []Team

type PostTeam struct {
	Name       string   `json:"name"`
	Permission string   `json:"permission,omitempty"`
	RepoNames  []string `json:"repo_names,omitempty"`
}

func (github *GitHubClient) getOrgs(res *http.Response) ([]Org, error) {
	orgJson, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	org := []Org{}
	if err = json.Unmarshal(orgJson, &org); err != nil {
		return nil, err
	}

	github.getLimits(res)
	return org, nil
}

//
// GitHub Doc - Orgs: List User Organizations
// Url: https://api.github.com/events?access_token=...
// Request Type: GET /user/orgs
// Access Token: REQUIRED
//
func (github *GitHubClient) GetUserOrgs(page int) ([]Org, error) {
	if page < 1 {
		return nil, errors.New("The page number is less then 1")
	}

	apiUrl := github.createUrl("/user/orgs?page=" + strconv.Itoa(page))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return github.getOrgs(res)
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

//
// GitHub Doc - Orgs: List User Organizations - Unauthorized user
// Url: https://api.github.com/users/:org/orgs
// Request Type: GET /users/:org/orgs
// Access Token: NONE
//
func (github *GitHubClient) GetPublicUserOrgs(user string, page int) ([]Org, error) {
	if page < 1 {
		return nil, errors.New("The page number is less then 1")
	}

	user = strings.TrimSpace(user)
	if len(user) == 0 {
		return nil, errors.New("The user data given does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/users/" + user + "/orgs?page=" + strconv.Itoa(page))
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return github.getOrgs(res)
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

func (github *GitHubClient) getOrg(res *http.Response) (*Org, error) {
	orgJson, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	org := &Org{}
	if err = json.Unmarshal(orgJson, org); err != nil {
		return nil, err
	}

	github.getLimits(res)
	return org, nil
}

//
// GitHub Doc - Orgs: Get an Organization
// Url: https://api.github.com/orgs/:org?access_token=...
// Request Type: GET /orgs/:org
// Access Token: REQUIRED
//
func (github *GitHubClient) GetOrgById(org string) (*Org, error) {
	org = strings.TrimSpace(org)
	if len(org) == 0 {
		return nil, errors.New("The org data given does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/orgs/" + org)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return github.getOrg(res)
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

//
// GitHub Doc - Orgs: Edit an Organization
// Url: https://api.github.com/orgs/:org?access_token=...
// Request Type: PATCH /orgs/:org
// Access Token: REQUIRED
//
func (github *GitHubClient) EditOrg(org string, orgData map[string]string) (*Org, error) {
	org = strings.TrimSpace(org)
	if len(org) == 0 {
		return nil, errors.New("The org given does not contain any non-whitespace content")
	}

	orgReader, err := github.CreateReader(orgData)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/orgs/" + org)
	apiRequest, err := http.NewRequest("PATCH", apiUrl, orgReader)
	if err != nil {
		return nil, err
	}

	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return github.getOrg(res)
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

func (github *GitHubClient) getUsers(res *http.Response) ([]GitUser, error) {
	usersJson, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	users := []GitUser{}
	if err = json.Unmarshal(usersJson, &users); err != nil {
		return nil, err
	}

	github.getLimits(res)
	return users, nil
}

// Org -  Members Section
//
// GitHub Doc - Orgs: Members list
// Url: https://api.github.com/orgs/:org/members?access_token=...
// Request Type: GET /orgs/:org/members
// Access Token: REQUIRED
//
func (github *GitHubClient) GetOrgMembers(org string) ([]GitUser, error) {
	org = strings.TrimSpace(org)
	if len(org) == 0 {
		return nil, errors.New("The org data given does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/orgs/" + org + "/members")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return github.getUsers(res)
	}

	// 302 means we're not a member of the organization and we should
	// check the public endpoint for members list
	// ref: http://developer.github.com/v3/orgs/members/#check-membership
	if res.StatusCode == 302 {
		github.GetPublicOrgMembers(org)
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

//
// GitHub Doc - Orgs: Check membership
// Url: https://api.github.com/orgs/:org/members/:user?access_token=...
// Request Type: GET /orgs/:org/members/:user
// Access Token: REQUIRED
// Returns: "member", "non-member", "unconfirmed"
//
func (github *GitHubClient) CheckOrgMembership(org, user string) (string, error) {
	org = strings.TrimSpace(org)
	if len(org) == 0 {
		return "", errors.New("The org given does not contain any non-whitespace content")
	}

	user = strings.TrimSpace(user)
	if len(user) == 0 {
		user = github.Login
	}

	apiUrl := github.createUrl("/orgs/" + org + "/members/" + user)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 204:
		return "member", nil
	case 404:
		return "non-member", nil
	case 302:
		res, err = github.Client.Get(res.Header.Get("Location"))
		if err != nil {
			return "unconfirmed", err
		}

		if res.StatusCode == 200 {
			return "member", nil
		} else if res.StatusCode == 404 {
			return "unconfirmed", nil
		}

	}

	return "", errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 	To add a member use the AddOrgTeamMember method created later in this code
// 	GitHub Doc - Orgs: Add a member - To add someone as a member to an org, you must add them to a team.
//
//
// GitHub Doc - Orgs: Remove a member
// Url: https://api.github.com/orgs/:org/members/:user?access_token=...
// Request Type: DELETE /orgs/:org/members/:user
// Access Token: REQUIRED
//
func (github *GitHubClient) RemoveOrgMember(org, user string) (bool, error) {
	org = strings.TrimSpace(org)
	if len(org) == 0 {
		return false, errors.New("The org given does not contain any non-whitespace content")
	}

	user = strings.TrimSpace(user)
	if len(user) == 0 {
		return false, errors.New("The org given does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/orgs/" + org + "/members/" + user)
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

//
// GitHub Doc - Orgs: Public members list - Members of an organization can choose to have their membership publicized or not.
// Url: https://api.github.com/orgs/:org/public_members?access_token=...
// Request Type: GET /orgs/:org/public_members
// Access Token: REQUIRED
//
func (github *GitHubClient) GetPublicOrgMembers(org string) ([]GitUser, error) {
	org = strings.TrimSpace(org)
	if len(org) == 0 {
		return nil, errors.New("The org data given does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/orgs/" + org + "/public_members")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return github.getUsers(res)
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

//
// GitHub Doc - Orgs: Check public membership
// Url: https://api.github.com/orgs/:org/public_members/:user?access_token=...
// Request Type: GET /orgs/:org/public_members/:user
// Access Token: REQUIRED
//
func (github *GitHubClient) CheckPublicOrgMembership(org, user string) (bool, error) {
	org = strings.TrimSpace(org)
	if len(org) == 0 {
		return false, errors.New("The org given does not contain any non-whitespace content")
	}

	user = strings.TrimSpace(user)
	if len(user) == 0 {
		user = github.Login
	}

	apiUrl := github.createUrl("/orgs/" + org + "/public_members/" + user)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	github.getLimits(res)
	if res.StatusCode == 204 {
		return true, nil
	} else if res.StatusCode == 404 {
		return false, nil
	}

	return false, errors.New("Didn't receive 204/404 status from Github: " + res.Status)
}

//
// GitHub Doc - Orgs: Publicize a user’s membership
// Url: https://api.github.com/orgs/:org/public_members/:user?access_token=...
// Request Type: PUT /orgs/:org/public_members/:user
// Access Token: REQUIRED
//
func (github *GitHubClient) PublishUserMembership(org, user string) (bool, error) {
	org = strings.TrimSpace(org)
	if len(org) == 0 {
		return false, errors.New("The org data given does not contain any non-whitespace content")
	}

	user = strings.TrimSpace(user)
	if len(user) == 0 {
		return false, errors.New("The user data given does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/orgs/" + org + "/public_members/" + user)
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
// GitHub Doc - Orgs: Conceal a user’s membership
// Url: https://api.github.com/orgs/:org/public_members/:user?access_token=...
// Request Type: DELETE /orgs/:org/public_members/:user
// Access Token: REQUIRED
//
func (github *GitHubClient) ConcealUserMembership(org, user string) (bool, error) {
	org = strings.TrimSpace(org)
	if len(org) == 0 {
		return false, errors.New("The org data given does not contain any non-whitespace content")
	}

	user = strings.TrimSpace(user)
	if len(user) == 0 {
		return false, errors.New("The user data given does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/orgs/" + org + "/public_members/" + user)
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

func (github *GitHubClient) getTeams(res *http.Response) ([]Team, error) {
	teamsJson, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	teams := &[]Team{}
	if err = json.Unmarshal(teamsJson, teams); err != nil {
		return nil, err
	}

	github.getLimits(res)
	return (*teams), nil
}

// Org - Team Section
//
// GitHub Doc - Orgs: List teams
// Url: https://api.github.com/orgs/:org/teams?access_token=...
// Request Type: GET /orgs/:org/teams
// Access Token: REQUIRED
//
func (github *GitHubClient) ListTeams(org string) ([]Team, error) {
	org = strings.TrimSpace(org)
	if len(org) == 0 {
		return nil, errors.New("The org data given does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/orgs/" + org + "/teams")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return github.getTeams(res)
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

func (github *GitHubClient) getTeam(res *http.Response) (*Team, error) {
	teamJson, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	team := &Team{}
	if err = json.Unmarshal(teamJson, team); err != nil {
		return nil, err
	}

	github.getLimits(res)
	return team, nil
}

//
// GitHub Doc - Orgs: Get Team
// Url: https://api.github.com/teams/:id?access_token=...
// Request Type: GET /teams/:id
// Access Token: REQUIRED
//
func (github *GitHubClient) GetTeam(id string) (*Team, error) {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return nil, errors.New("The org data given does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/teams/" + id)
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return github.getTeam(res)
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

//
// GitHub Doc: Gists: Create team
// Url: https://api.github.com/orgs/:org/teams?access_token=...
// Request Type: POST /orgs/:org/teams
// Access Token: REQUIRED
//
func (github *GitHubClient) CreateTeam(org string, postTeam *PostTeam) (*Team, error) {
	org = strings.TrimSpace(org)
	if len(org) == 0 {
		return nil, errors.New("The value of org does not contain any non-whitespace content")
	}

	postTeam.Name = strings.TrimSpace(postTeam.Name)
	if len(postTeam.Name) == 0 {
		return nil, errors.New("The value of postTeam.Name does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/orgs/" + org + "/teams")
	teamReader, err := github.CreateReader(postTeam)
	if err != nil {
		return nil, err
	}

	res, err := github.Client.Post(apiUrl, "application/json", teamReader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		return github.getTeam(res)
	}

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

//
// GitHub Doc: Orgs - Edit Team
// Url: https://api.github.com/gists:id?access_token=...
// Request Type: PATCH /gists/:id
// Access Token: REQUIRED
//
func (github *GitHubClient) EditTeam(id string, teamData map[string]string) (*Team, error) {
	id = strings.TrimSpace(id)
	if len(id) < 1 {
		return nil, errors.New("The id must have a length greater then zero.")
	}

	teamData["name"] = strings.TrimSpace(teamData["name"])
	if len(teamData["name"]) == 0 {
		return nil, errors.New("The value of postTeam.Name does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/teams/" + id)
	apiReader, err := github.CreateReader(teamData)
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
		return github.getTeam(res)
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

//
// GitHub Doc: Org - Delete team
// Url: https://api.github.com/teams/:id?access_token=...
// Request Type: DELETE /teams/:id
// Access Token: REQUIRED
//
func (github *GitHubClient) DeleteTeam(id string) (bool, error) {
	id = strings.TrimSpace(id)
	if len(id) < 1 {
		return false, errors.New("The id does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/teams/" + id)
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

//
// GitHub Doc - Orgs: List team members
// Url: https://api.github.com/teams/:id/members?access_token=...
// Request Type: GET /teams/:id/members
// Access Token: REQUIRED
//
func (github *GitHubClient) ListTeamMembers(id string) ([]GitUser, error) {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return nil, errors.New("The id value given does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/teams/" + id + "/members")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return github.getUsers(res)
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

//
// GitHub Doc - Orgs: Get team member
// Url: https://api.github.com/teams/:id/members/:user?access_token=...
// Request Type: GET /teams/:id/members/:user
// Access Token: REQUIRED
//
func (github *GitHubClient) GetTeamMember(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapStrings([]string{"id", "user"}, urlData); !ok {
		return false, errors.New("Data missing to create the url is missing. Both user and id are required fields for this map.")
	}

	apiUrl := github.createUrl("/teams/" + urlData["id"] + "/members/" + urlData["user"])
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

//
// GitHub Doc - Orgs: Add team member
// Url: https://api.github.com/teams/:id/members/:user?access_token=...
// Request Type: PUT /teams/:id/members/:user
// Access Token: REQUIRED
//
func (github *GitHubClient) AddTeamMember(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapStrings([]string{"id", "user"}, urlData); !ok {
		return false, errors.New("Data missing to create the url is missing. Both user and id are required fields for this map.")
	}

	apiUrl := github.createUrl("/teams/" + urlData["id"] + "/members/" + urlData["user"])
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
	} else if res.StatusCode == 422 {
		github.getLimits(res)
		return false, errors.New("Cannot add an organization to a Team.")
	}

	return false, errors.New("Didn't receive 204/422 status from Github: " + res.Status)
}

//
// GitHub Doc - Orgs: Remove team member
// Url: https://api.github.com/teams/:id/members/:user?access_token=...
// Request Type: DELETE /teams/:id/members/:user
// Access Token: REQUIRED
//
func (github *GitHubClient) RemoveTeamMember(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapStrings([]string{"id", "user"}, urlData); !ok {
		return false, errors.New("Data to create the url is missing. Both user and id are required fields for this map.")
	}

	apiUrl := github.createUrl("/teams/" + urlData["id"] + "/members/" + urlData["user"])
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

//
// GitHub Doc - Orgs: List team repos
// Url: https://api.github.com/teams/:id/repos?access_token=...
// Request Type: GET /teams/:id/repos
// Access Token: REQUIRED
//
func (github *GitHubClient) ListTeamRepos(id string) (*Repos, error) {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return nil, errors.New("The id value given does not contain any non-whitespace content")
	}

	apiUrl := github.createUrl("/teams/" + id + "/repos")
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
// GitHub Doc - Orgs: Get team repo
// Url: https://api.github.com/teams/:id/repos/:owner/:repo?access_token=...
// Request Type: GET /teams/:id/repos/:owner/:repo
// Access Token: REQUIRED
//
func (github *GitHubClient) GetTeamRepo(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapStrings([]string{"id", "repo"}, urlData); !ok {
		return false, errors.New("Data to create the url is missing. Both user and id are required fields for this map.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/teams/" + urlData["id"] + "/repos/" + urlData["owner"] + "/" + urlData["repo"])
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

	return false, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

//
// GitHub Doc - Orgs: Add team repo
// Url: https://api.github.com/teams/:id/repos/:owner/:repo?access_token=...
// Request Type: PUT /teams/:id/repos/:owner/:repo
// Access Token: REQUIRED
//
func (github *GitHubClient) AddTeamRepo(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapStrings([]string{"id", "repo"}, urlData); !ok {
		return false, errors.New("Data to create the url is missing. Both user and id are required fields for this map.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/teams/" + urlData["id"] + "/repos/" + urlData["owner"] + "/" + urlData["repo"])
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
	} else if res.StatusCode == 422 {
		github.getLimits(res)
		return false, errors.New("It isn't possible to add a organizai=tion ")
	}

	return false, errors.New("Didn't receive 204/422 status from Github: " + res.Status)
}

//
// GitHub Doc - Orgs: Remove team repo
// Url: https://api.github.com/teams/:id/members/:user?access_token=...
// Request Type: DELETE /teams/:id/members/:user
// Access Token: REQUIRED
//
func (github *GitHubClient) RemoveTeamRepo(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapStrings([]string{"id", "repo"}, urlData); !ok {
		return false, errors.New("Data to create the url is missing. Both user and id are required fields for this map.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/teams/" + urlData["id"] + "/repos/" + urlData["owner"] + "/" + urlData["repo"])
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
