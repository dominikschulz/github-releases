package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
)

var apiURL = "https://api.github.com/repos/%s/%s/releases"

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
	Version string `arg:""`
	URL     bool   `arg:""`
}

func main() {
	arg.MustParse(&args)

	r, err := fetchLatestStableRelease(args.User, args.Project)
	if err != nil {
		fmt.Printf("Failed to fetch releases for %s/%s: %s", args.User, args.Project, err)
		os.Exit(1)
	}

	if len(args.Version) < 1 {
		fmt.Println(r.Name)
		os.Exit(0)
	}
	args.Version = strings.TrimPrefix(args.Version, "v")
	r.Name = strings.TrimPrefix(r.Name, "v")
	if r.Name != args.Version {
		fmt.Printf("Not latest. Your Version %s - Latest: %s\n", args.Version, r.Name)
		if len(r.Assets) > 0 && args.URL {
			fmt.Printf("URL: %s\n", r.Assets[0].URL)
		}
		os.Exit(1)
	}
	os.Exit(0)
}

func fetchLatestStableRelease(user, project string) (Release, error) {
	url := fmt.Sprintf(apiURL, user, project)
	resp, err := http.Get(url)
	if err != nil {
		return Release{}, fmt.Errorf("Failed to fetch from %s: %s", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Release{}, fmt.Errorf("Failed to fetch from %s: %d - %s", url, resp.StatusCode, resp.Status)
	}
	var rs []Release
	err = json.NewDecoder(resp.Body).Decode(&rs)
	if err != nil {
		return Release{}, err
	}
	if len(rs) < 1 {
		return Release{}, fmt.Errorf("No releases")
	}
	for _, r := range rs {
		if strings.Contains(r.Name, "beta") || strings.Contains(r.Name, "rc") {
			continue
		}
		return r, nil
	}
	return Release{}, fmt.Errorf("No stable release found")
}
