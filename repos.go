package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"
)

type repo struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	FullName        string `json:"full_name"`
	Private         bool   `json:"private"`
	BlobsURL        string `json:"blobs_url"`
	BranchesURL     string `json:"branches_url"`
	CommitsURL      string `json:"commits_url"`
	ContentsURL     string `json:"contents_url"`
	ContributorsURL string `json:"contributors_url"`
	DefaultBranch   string `json:"default_branch"`
	Permissions     struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permissions"`
	License struct {
		Key    string `json:"key"`
		Name   string `json:"name"`
		SpdxID string `json:"spdx_id"`
		URL    string `json:"url"`
		NodeID string `json:"node_id"`
	} `json:"license"`
	ObservedLicenses []string
}

func (c *client) getRepos() []repo {
	s := time.Now()
	log.Printf("Fetching list of %s repos... ", c.org)
	path := fmt.Sprintf("/orgs/%s/repos", c.org)
	pages := c.fetchPaginated(path)
	allRepos := []repo{}
	for _, page := range pages {
		repos := []repo{}
		if err := json.Unmarshal(page.body, &repos); err != nil {
			log.Fatal(err)
		}
		for _, repo := range repos {
			allRepos = append(allRepos, repo)
		}
	}
	log.Printf("OK\nGot %d repos (%d requests, %s)\n", len(allRepos), len(pages), getTimeSinceMs(s))
	return allRepos
}

type commit struct {
	Sha string `json:"sha"`
}

func (c *client) getLastCommit(repo string) (commit, error) {
	var commits []commit
	path := fmt.Sprintf("/repos/%s/%s/commits", c.org, repo)
	result := c.fetch(path)
	if err := json.Unmarshal(result.body, &commits); err != nil {
		log.Printf("Couldn't get last commit hash for %s/%s", c.org, repo)
		log.Println("This is likely an empty or nonstandard repo")
		return commit{}, err
	}
	return commits[0], nil
}

type blob struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	Sha  string `json:"sha"`
}

type tree struct {
	Sha  string `json:"sha"`
	URL  string `json:"url"`
	Tree []blob `json:"tree"`
}

func (c *client) getCurrentTree(repo string) tree {
	last, err := c.getLastCommit(repo)
	if err != nil {
		return tree{}
	}
	path := fmt.Sprintf("/repos/%s/%s/git/trees/%s", c.org, repo, last.Sha)
	var t tree
	result := c.fetch(path)
	if err := json.Unmarshal(result.body, &t); err != nil {
		log.Printf("Couldn't get tree for %s/%s", c.org, repo)
	}
	return t
}

func (c *client) getRepoLicense(repo string) []string {
	licenseRegexp, _ := regexp.Compile(`^(?i)licen[c/s]e(.?)(txt|md|-.*)?$`)
	t := c.getCurrentTree(repo)
	licenses := []string{}
	for _, b := range t.Tree {
		pathLower := strings.ToLower(b.Path)
		if licenseRegexp.MatchString(pathLower) || strings.HasPrefix(pathLower, "copyright") {
			licenses = append(licenses, b.Path)
		}
	}
	if len(licenses) == 0 {
		licenses = []string{"NONE"}
	}
	return licenses
}

func licensesAreGood(licenses []string) bool {
	if len(licenses) != 3 {
		return false
	}
	sort.Strings(licenses)
	if !strings.HasPrefix(strings.ToLower(licenses[0]), "copyright") {
		return false
	}
	if !strings.HasPrefix(strings.ToLower(licenses[1]), "license-apache") {
		return false
	}
	if !strings.HasPrefix(strings.ToLower(licenses[2]), "license-mit") {
		return false
	}
	return true
}
