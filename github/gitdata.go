package github

// Git DB API
//
//	## Blobs API
//		-  Get a Blob
//		-  Create a Blob
//		-  Custom media types
//
//	## Commits API
//		-  Get a Commit
//		-  Create a Commit
//
//	## References API
//		-  Get a Reference
//		-  Get all References
//		-  Create a Reference
//		-  Update a Reference
//		-  Delete a Reference
//
//	## Tags API
//		-  Get a Tag
//		-  Create a Tag Object
//
//	## Trees API
//		-  Get a Tree
//		-  Get a Tree Recursively
//		-  Create a Tree

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Structs for Git Data
type Blob struct {
	Content  string `json:"content,omitempty"`
	Encoding string `json:"encoding,omitempty"`
	SHA      string `json:"sha"`
	Size     uint   `json:"size,omitempty"`
}

// Commit returned from the API- Get Commit Data
type DataCommit struct {
	SHA       string              `json:"sha"`
	Url       string              `json:"url"`
	Author    map[string]string   `json:"author"`
	Committer map[string]string   `json:"committer"`
	Message   string              `json:"message"`
	Tree      map[string]string   `json:"tree"`
	Parents   []map[string]string `json:"parents"`
}

// Commit sent to the API - Create Commit Data
type CreateDataCommit struct {
	Message   string            `json:"message"`
	Author    map[string]string `json:"author"`
	Committer map[string]string `json:"committer"`
	Parents   []string          `json:"parents"`
	Tree      string            `json:"tree"`
}

type TreeNode struct {
	Path string `json:"path,omitempty"`
	Mode string `json:"mode,omitempty"`
	Type string `json:"type,omitempty"`
	Size int    `json:"size,omitempty"`
	SHA  string `json:"sha"`
	Url  string `json:"url"`
}

type Tree struct {
	SHA  string     `json:"sha"`
	Url  string     `json:"url"`
	Tree []TreeNode `json:"tree,omitempty"`
}

type CreateTreeNode struct {
	Path string `json:"path,omitempty"`
	Mode string `json:"mode,omitempty"`
	SHA  string `json:"sha"`
	Type string `json:"type,omitempty"`
}

type CreateTree struct {
	BaseTree string            `json:"base_tree"`
	Tree     []*CreateTreeNode `json:"tree"`
}

type DataTag struct {
	Tag     string            `json:"tag"`
	SHA     string            `json:"sha,omitempty"`
	Url     string            `json:"url,omitempty"`
	Message string            `json:"message"`
	Tagger  map[string]string `json:"tagger"`
	Object  map[string]string `json:"object"`
}

type Reference struct {
	Ref    string            `json:"ref"`
	Url    string            `json:"url"`
	Object map[string]string `json:"object"`
}

type UpdateRef struct {
	SHA   string `json:"sha"`
	Force bool   `json:"force"`
}

// GitData Functions

// Blobs Section
// 
// GitHub Doc - GitData: Blobs - Get a Blob
// Url: https://api.github.com/repos/:owner/:repo/git/blobs/:sha?access_token=...
// Request Type: GET /repos/:owner/:repo/git/blobs/:sha 
// Access Token: REQUIRED
// 

func (github *GitHubClient) GetBlob(urlData map[string]string) (*Blob, error) {
	if ok := github.AssertMapStrings([]string{"repo", "sha"}, urlData); !ok {
		return nil, errors.New("One or more fields are missing and/or do not have content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/git/blobs/" + urlData["sha"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		blob := &Blob{}
		blobJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(blobJson, blob); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return blob, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Docs - GitData: Blobs - Create a Blob
// Url: https://api.github.com/repos/:owner/:repo/git/blobs?access_token=...
// Request Type: POST /repos/:owner/:repo/git/blobs
// Access Token: REQUIRED
// 

func (github *GitHubClient) CreateBlob(urlData, postData map[string]string) (*Blob, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("Your repo in your urlData is either missing or has a length of zero.")
	}
	if ok := github.AssertMapStrings([]string{"encoding", "content"}, postData); !ok {
		return nil, errors.New("One or more fields are missing and/or do not have content in  your post content.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	blobReader, err := github.CreateReader(postData)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/git/blobs")
	res, err := github.Client.Post(apiUrl, "application/json", blobReader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		blob := &Blob{}
		blobJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(blobJson, blob); err != nil {
			return nil, err
		}

		blob.Content = postData["content"]
		blob.Encoding = postData["encoding"]
		github.getLimits(res)
		return blob, nil
	}

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// GitData - Commits Section
// 
// GitHub Doc: GitData: Commits - Get a Commit
// Url: https://api.github.com/repos/:owner/:repo/git/blobs/:sha?access_token=...
// Request Type: GET /repos/:owner/:repo/git/blobs/:sha 
// Access Token: REQUIRED
// urlData{ "owner": string, "repo": string, "sha": string}
// 

func (github *GitHubClient) GetCommit(urlData map[string]string) (*DataCommit, error) {
	if ok := github.AssertMapStrings([]string{"repo", "sha"}, urlData); !ok {
		return nil, errors.New("urlData has insufficient data to make a request of the GitHub API.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/git/commits/" + urlData["sha"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		commit := &DataCommit{}
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

// 
// GitHub Doc: GitData: Commits - Create a Commit
// Url: https://api.github.com/repos/:owner/:repo/git/blobs?access_token=...
// Request Type: POST /repos/:owner/:repo/git/blobs
// Access Token: REQUIRED
// 

func (github *GitHubClient) CreateCommit(urlData map[string]string, commitData *CreateDataCommit) (*DataCommit, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("Your repo in your urlData is either missing or has a length of zero.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	blobReader, err := github.CreateReader(commitData)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/git/commits")
	res, err := github.Client.Post(apiUrl, "application/json", blobReader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		commit := &DataCommit{}
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

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// GitData - Tree Section
// 
// GitHub Doc: GitData: Trees - Get a Tree
// Url: https://api.github.com/repos/:owner/:repo/git/trees/:sha?access_token=...
// Request Type: GET /repos/:owner/:repo/git/trees/:sha
// Access Token: REQUIRED
// urlData{ "owner": string, "repo": string, "sha": string}
// 

func (github *GitHubClient) GetTree(urlData map[string]string) (*Tree, error) {
	if ok := github.AssertMapStrings([]string{"repo", "sha"}, urlData); !ok {
		return nil, errors.New("urlData has insufficient data to make a request of the GitHub API.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/git/trees/" + urlData["sha"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		tree := &Tree{}
		treeJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(treeJson, tree); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return tree, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 5f3a81f2aba703c00ef3341360300afe84ae895e
// 
// GitHub Doc: GitData: Trees - Get a Tree Recursively
// Url: https://api.github.com/repos/:owner/:repo/git/trees/:sha?recursive=1&access_token=...
// Request Type: GET /repos/:owner/:repo/git/trees/:sha?recursive=1
// Access Token: REQUIRED
// urlData{ "owner": string, "repo": string, "sha": string}
// 

func (github *GitHubClient) GetRecursiveTree(urlData map[string]string) (*Tree, error) {
	if ok := github.AssertMapStrings([]string{"repo", "sha"}, urlData); !ok {
		return nil, errors.New("urlData has insufficient data to make a request of the GitHub API.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/git/trees/" + urlData["sha"] + "?recursive=1")
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		tree := &Tree{}
		treeJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(treeJson, tree); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return tree, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: GitData: Trees - Create a Tree
// Url: https://api.github.com/repos/:owner/:repo/git/trees?access_token=...
// Request Type: POST /repos/:owner/:repo/git/trees
// Access Token: REQUIRED
// 

func (github *GitHubClient) CreateTree(urlData map[string]string, treeData *CreateTree) (*Tree, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("Your repo in your urlData is either missing or has a length of zero.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	treeReader, err := github.CreateReader(treeData)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/git/trees")
	res, err := github.Client.Post(apiUrl, "application/json", treeReader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		tree := &Tree{}
		treeJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(treeJson, tree); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return tree, nil
	}

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// GitData - Tag Section
// 
// GitHub Doc: GitData: Trees - Get a Tree Recursively
// Url: https://api.github.com/repos/:owner/:repo/git/tags/:sha?access_token=...
// Request Type: GET /repos/:owner/:repo/git/tags/:sha
// Access Token: REQUIRED
// urlData{ "owner": string, "repo": string, "sha": string}
// 

func (github *GitHubClient) GetTag(urlData map[string]string) (*DataTag, error) {
	if ok := github.AssertMapStrings([]string{"repo", "sha"}, urlData); !ok {
		return nil, errors.New("urlData has insufficient data to make a request of the GitHub API.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/git/tags/" + urlData["sha"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		tag := &DataTag{}
		tagJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(tagJson, tag); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return tag, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: GitData: Tags - Create a Tag Object
// Url: https://api.github.com/repos/:owner/:repo/git/tags?access_token=...
// Request Type: POST /repos/:owner/:repo/git/tags
// Access Token: REQUIRED
// 

func (github *GitHubClient) CreateTag(urlData map[string]string, tag *DataTag) (*DataTag, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("Your repo in your urlData is either missing or has a length of zero.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	tagReader, err := github.CreateReader(tag)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/git/tag")
	res, err := github.Client.Post(apiUrl, "application/json", tagReader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		tag := &DataTag{}
		tagJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(tagJson, tag); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return tag, nil
	}

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// GitData - Reference Section
// 
// GitHub Doc: GitData: Reference - Get a Reference
// Url: https://api.github.com/repos/:owner/:repo/git/refs/:ref?access_token=...
// Request Type: GET /repos/:owner/:repo/git/refs/:ref
// Access Token: REQUIRED
// urlData{ "owner": string, "repo": string, "ref": string}
// 

func (github *GitHubClient) GetRef(urlData map[string]string) (*Reference, error) {
	if ok := github.AssertMapStrings([]string{"repo", "ref"}, urlData); !ok {
		return nil, errors.New("urlData has insufficient data to make a request of the GitHub API.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/git/refs/" + urlData["ref"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		ref := &Reference{}
		refJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(refJson, ref); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return ref, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: GitData: Reference - Get all References
// Url: https://api.github.com/repos/:owner/:repo/git/refs/:ref?access_token=...
// Request Type: GET /repos/:owner/:repo/git/refs
// Access Token: REQUIRED
// urlData{ "owner": string, "repo": string, "ref": string}
// 

func (github *GitHubClient) GetAllRefs(urlData map[string]string) ([]Reference, error) {
	if ok := github.AssertMapStrings([]string{"repo", "ref"}, urlData); !ok {
		return nil, errors.New("urlData has insufficient data to make a request of the GitHub API.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/git/refs/" + urlData["ref"])
	res, err := github.Client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		refs := &[]Reference{}
		refJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(refJson, refs); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return (*refs), nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: GitData: Tags - Create a Tag Object
// Url: https://api.github.com/repos/:owner/:repo/git/tags?access_token=...
// Request Type: POST /repos/:owner/:repo/git/tags
// Access Token: REQUIRED
// 

func (github *GitHubClient) CreateRef(urlData map[string]string, refData map[string]string) (*Reference, error) {
	if ok := github.AssertMapString("repo", urlData); !ok {
		return nil, errors.New("Your repo in your urlData is either missing or has a length of zero.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	refReader, err := github.CreateReader(refData)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/git/refs")
	res, err := github.Client.Post(apiUrl, "application/json", refReader)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		ref := &Reference{}
		refJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(refJson, ref); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return ref, nil
	}

	return nil, errors.New("Didn't receive 201 status from Github: " + res.Status)
}

// 
// GitHub Doc: GitData: Tags - Create a Tag Object
// Url: https://api.github.com/repos/:owner/:repo/git/tags?access_token=...
// Request Type: POST /repos/:owner/:repo/git/tags
// Access Token: REQUIRED
// 

func (github *GitHubClient) EditRef(urlData map[string]string, refData *UpdateRef) (*Reference, error) {
	if ok := github.AssertMapStrings([]string{"repo", "ref"}, urlData); !ok {
		return nil, errors.New("Your repo in your urlData is either missing or has a length of zero.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	refReader, err := github.CreateReader(refData)
	if err != nil {
		return nil, err
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/git/refs/" + urlData["ref"])
	apiRequest, err := http.NewRequest("PATCH", apiUrl, refReader)
	if err != nil {
		return nil, err
	}

	res, err := github.Client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		ref := &Reference{}
		refJson, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(refJson, ref); err != nil {
			return nil, err
		}

		github.getLimits(res)
		return ref, nil
	}

	return nil, errors.New("Didn't receive 200 status from Github: " + res.Status)
}

// 
// GitHub Doc: GitData: Tags - Create a Tag Object
// Url: https://api.github.com/repos/:owner/:repo/git/tags?access_token=...
// Request Type: POST /repos/:owner/:repo/git/tags
// Access Token: REQUIRED
// 

func (github *GitHubClient) DeleteRef(urlData map[string]string) (bool, error) {
	if ok := github.AssertMapStrings([]string{"repo", "ref"}, urlData); !ok {
		return false, errors.New("Your repo in your urlData is either missing or has a length of zero.")
	}
	if ok := github.AssertMapString("owner", urlData); !ok {
		urlData["owner"] = github.Login
	}

	apiUrl := github.createUrl("/repos/" + urlData["owner"] + "/" + urlData["repo"] + "/git/refs/" + urlData["ref"])
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
