package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type contributor struct {
	Login         string `json:"login"`
	Email         string `json:"email"`
	Contributions int    `json:"contributions"`
}

func (c *client) getContributorsByRepo(repo string) []contributor {
	path := fmt.Sprintf("/repos/%s/%s/contributors", c.org, repo)
	pages := c.fetchPaginated(path)
	allContributors := []contributor{}
	for _, page := range pages {
		if page.statusCode == 204 {
			continue
		}
		contributors := []contributor{}
		if err := json.Unmarshal(page.body, &contributors); err != nil {
			log.Println(page.statusCode)
			log.Fatalf("Error retrieving %s: %s\n", path, err.Error())
		}
		allContributors = append(allContributors, contributors...)
	}
	return allContributors
}
