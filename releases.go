package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexflint/go-arg"
)

var apiURL = "https://api.github.com/repos/%s/%s/releases/latest"

type Asset struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
}

type Release struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []Asset   `json:"assets"`
}

var args struct {
	User    string `arg:"required"`
	Project string `arg:"required"`
	Version string `arg:"required"`
	URL     bool   `arg:""`
}

func main() {
	arg.MustParse(&args)

	r, err := fetchRelease(args.User, args.Project)
	if err != nil {
		fmt.Printf("Failed to fetch releases for %s/%s: %s", args.User, args.Project, err)
		os.Exit(1)
	}

	if r.Name != args.Version {
		fmt.Printf("Latest: %s\n", r.Name)
		if len(r.Assets) > 0 && args.URL {
			fmt.Printf("URL: %s\n", r.Assets[0].URL)
		}
		os.Exit(1)
	}
	os.Exit(0)
}

func fetchRelease(user, project string) (*Release, error) {
	url := fmt.Sprintf(apiURL, user, project)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from %s: %s", url, err)
	}
	defer resp.Body.Close()

	var r Release
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
