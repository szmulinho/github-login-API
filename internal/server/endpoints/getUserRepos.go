package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getUserRepoAccesses(githubUsername string) []RepoAccess {
	// Get the user's access to all repositories.
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100", githubUsername)
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// Parse the JSON response.
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	var repos []Repo
	err = json.Unmarshal(body, &repos)
	if err != nil {
		panic(err)
	}

	// Filter the repositories to only include those that the user has write access to.
	repoAccesses := []RepoAccess{}
	for _, repo := range repos {
		if repo.Permissions.Push {
			repoAccesses = append(repoAccesses, RepoAccess{
				RepoName:    repo.Name,
				Permissions: repo.Permissions,
			})
		}
	}

	return repoAccesses
}
