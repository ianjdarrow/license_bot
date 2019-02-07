package main

import (
	"log"
	"sort"
	"sync"
	"time"

	"github.com/gosuri/uiprogress"
)

func (c *client) getAllContributors() []contributor {
	repos := c.getRepos()
	log.Printf("Getting all %s contributors...", c.org)
	uiprogress.Start()
	bar := getProgressBar(len(repos))

	throttle := time.Tick(rateLimit)
	mux := &sync.Mutex{}
	wg := sync.WaitGroup{}

	allContributors := []contributor{}
	for _, repo := range repos {
		wg.Add(1)
		<-throttle
		go func(r string) {
			defer wg.Done()
			contributors := c.getContributorsByRepo(r)
			mux.Lock()
			for _, con := range contributors {
				allContributors = addToContributorSet(allContributors, con)
			}
			mux.Unlock()
			bar.Incr()
		}(repo.Name)
	}
	wg.Wait()

	sort.Slice(allContributors, func(i, j int) bool {
		return allContributors[i].Contributions > allContributors[j].Contributions
	})
	return allContributors
}

func addToContributorSet(s []contributor, con contributor) []contributor {
	for i, el := range s {
		if el.Login == con.Login {
			s[i] = contributor{
				Email:         con.Email,
				Contributions: con.Contributions + el.Contributions,
				Login:         con.Login,
			}
			return s
		}
	}
	return append(s, con)
}

func (c *client) getAllLicenses() []repo {
	repos := c.getRepos()
	log.Printf("Getting all %s licenses...", c.org)

	withLicenses := []repo{}

	throttle := time.Tick(rateLimit)
	mux := &sync.Mutex{}
	wg := sync.WaitGroup{}
	uiprogress.Start()
	bar := getProgressBar(len(repos))

	for _, i := range repos {
		go func(r repo) {
			wg.Add(1)
			defer wg.Done()
			<-throttle
			bar.Incr()
			r.ObservedLicenses = c.getRepoLicense(r.Name)
			mux.Lock()
			withLicenses = append(withLicenses, r)
			mux.Unlock()
		}(i)
	}
	wg.Wait()

	return withLicenses
}
