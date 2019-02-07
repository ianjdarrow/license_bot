package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/gosuri/uiprogress"
)

const rateLimit = time.Second / 20

// Client helps manage authenticated requests
type client struct {
	http.Client
	authToken string
	org       string
}

type response struct {
	statusCode int
	header     http.Header
	body       []byte
}

func (c *client) fetch(path string) response {
	url := fmt.Sprintf("https://api.github.com%s", path)
	req, _ := http.NewRequest("GET", url, nil)
	c.setRequestHeaders(req)
	resp, err := c.Do(req)
	if err != nil || resp.StatusCode > 299 {
		log.Printf("Error getting %s: %d", path, resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body for %s: %s", path, err.Error())
	}
	return response{
		statusCode: resp.StatusCode,
		header:     resp.Header,
		body:       body,
	}
}

func (c *client) fetchPaginated(path string) []response {
	initialResponse := c.fetch(path)
	pages := getPageCountFromHeader(initialResponse.header)

	responses := []response{initialResponse}

	throttle := time.Tick(rateLimit)
	mux := &sync.Mutex{}
	wg := sync.WaitGroup{}

	for i := 2; i < pages+1; i++ {
		wg.Add(1)
		<-throttle
		go func(page int) {
			defer wg.Done()
			pagePath := fmt.Sprintf("%s?page=%d", path, page)
			resp := c.fetch(pagePath)
			mux.Lock()
			responses = append(responses, resp)
			mux.Unlock()
		}(i)
	}
	wg.Wait()
	return responses
}

func (c *client) setRequestHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/vnd.github.inertia-preview+json")
	req.Header.Set("User-Agent", os.Getenv("GITHUB_USER"))
	req.Header.Set("Authorization", c.authToken)
}

func (c *client) setAuthToken() {
	credentials := fmt.Sprintf("%s:%s", os.Getenv("GITHUB_USER"), os.Getenv("GITHUB_TOKEN"))
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	c.authToken = fmt.Sprintf("Basic %s", string(encoded))
}

func getPageCountFromHeader(h http.Header) int {
	linkHeader, ok := h["Link"]
	if !ok {
		return 1
	}
	r := regexp.MustCompile(`.*(\d)>; rel="last"$`)
	pages := r.FindStringSubmatch(linkHeader[0])
	pageCount, err := strconv.Atoi(pages[1])
	if err != nil {
		log.Fatal("Error parsing page count header")
	}
	return pageCount
}

func getProgressBar(actions int) *uiprogress.Bar {
	bar := uiprogress.AddBar(actions - 1)
	bar.PrependElapsed()
	bar.AppendFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("%d/%d repos", b.Current(), b.Total)
	})
	return bar
}
