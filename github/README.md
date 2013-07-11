GitHub API v3 for the Go Programming Language.
=============
Alpha stages right now, feel free to contribute.

TODO
  -  Tests
  -  More Organizing
  -  Better Docs

*Includes all parts of the GitHub API except:*
  -  OAuth Related Stuff
  -  Repo - Downloads  // These are DEPRECATED
  -  Repo - Comments
  -  PubSubHubbub Method(s)

## Install

Step 1: Get Repo<br>
```go get github.com/CodeHub-io/Go-GitHub-API```

Step 2: Set Application Data<br>
Constants - Set these to make the library work - one time thing<br>
```
github.go - line 14
BASEPATH    - File Path - This is where the zip/tar.gz files are stored
GITID       - This is your github application id
GITSECRET   - This is your github secret id
```

## Use Package

```import github "github.com/CodeHub-io/Go-GitHub-API"```

## Contributors

[franckcuny](https://github.com/franckcuny)


## Creating the Client 

```Go
package main

import (
    "fmt"
    github "github.com/CodeHub-io/Go-GitHub-API"
)

func main() {
    // Create A GitHub Client
	githubClient := github.NewGitHubClient("Super-Secret-GitHub-API-Token", "Username")
	urlData := map[string]string{"type": "all"}
    
	repos, err := githubClient.GetUserRepos(urlData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", repos)
}
```



