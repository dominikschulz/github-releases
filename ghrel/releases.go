package ghrel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/blang/semver"
)

var (
	BaseURL = "https://api.github.com/repos/%s/%s/releases"
	sem     = regexp.MustCompile(`(?:^|\D)(\d+\.\d+\.\d+\S*)(?:$|\s)`)
)

type Asset struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
}

type Release struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	TagName     string    `json:"tag_name"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []Asset   `json:"assets"`
}

func (r Release) Version() semver.Version {
	match := sem.FindStringSubmatch(r.TagName)
	if len(match) < 2 {
		match = sem.FindStringSubmatch(r.Name)
	}
	if len(match) < 2 {
		return semver.Version{}
	}
	if sv, err := semver.ParseTolerant(match[1]); err == nil {
		return sv
	}
	return semver.Version{}
}

type Releases []Release

func (r Releases) Len() int {
	return len(r)
}

func (r Releases) Less(i, j int) bool {
	return r[j].Version().LT(r[i].Version())
}

func (r Releases) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func fetchReleases(user, project string) ([]Release, error) {
	url := fmt.Sprintf(BaseURL, user, project)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from %s: %s", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to fetch from %s: %d - %s", url, resp.StatusCode, resp.Status)
	}
	var rs []Release
	err = json.NewDecoder(resp.Body).Decode(&rs)
	if err != nil {
		return nil, err
	}
	sort.Sort(Releases(rs))

	return rs, nil
}

func filterStableReleases(rs []Release) []Release {
	out := make([]Release, 0, len(rs))
	for _, r := range rs {
		if strings.Contains(r.Name, "beta") || strings.Contains(r.Name, "rc") || r.Draft || r.Prerelease {
			continue
		}
		out = append(out, r)
	}
	return out
}

// FetchAllReleases will return all releases. The latest release will be at
// position 0.
func FetchAllReleases(user, project string) ([]Release, error) {
	rs, err := fetchReleases(user, project)
	if err != nil {
		return nil, err
	}
	if len(rs) < 1 {
		return rs, fmt.Errorf("No releases")
	}
	return rs, nil
}

// FetchLatestRelease will simply return the latested release, possibly a pre
// release.
func FetchLatestRelease(user, project string) (Release, error) {
	rs, err := FetchAllReleases(user, project)
	if err != nil {
		return Release{}, err
	}
	return rs[0], nil
}

// FetchAllStableReleases will return all stable releases. The latest release
// will be at position 0.
func FetchAllStableReleases(user, project string) ([]Release, error) {
	rs, err := fetchReleases(user, project)
	if err != nil {
		return []Release{}, err
	}
	if len(rs) < 1 {
		return []Release{}, fmt.Errorf("No releases")
	}
	return filterStableReleases(rs), nil
}

// FetchLatestStableRelease will return the latest stable release. This will
// exclude any releases marked as draft, prerelease or containing a pre-release
// marker in the name
func FetchLatestStableRelease(user, project string) (Release, error) {
	rs, err := FetchAllStableReleases(user, project)
	if err != nil {
		return Release{}, err
	}
	return rs[0], nil
}
