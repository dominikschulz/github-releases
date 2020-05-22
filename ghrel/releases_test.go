package ghrel

import (
	"testing"

	"github.com/blang/semver"
)

func TestVersion(t *testing.T) {
	for _, tc := range []struct {
		rs []Release
		sv semver.Version
	}{
		{
			rs: []Release{
				{
					Name:    "1.2.0 / 2017-01-01",
					TagName: "v1.2.0",
				},
			},
			sv: semver.Version{
				Major: 1,
				Minor: 2,
				Patch: 0,
			},
		},
		{
			rs: []Release{
				{
					Name: "2017-01-01 / 1.2.0",
				},
			},
			sv: semver.Version{
				Major: 1,
				Minor: 2,
				Patch: 0,
			},
		},
		{
			rs: []Release{
				{
					Name:  "The great Emperor -  v1.2.0",
					Draft: true,
				},
				{
					Name:       "The Fall -  v0.2.0",
					Prerelease: true,
				},
				{
					Name: "The Empire strikes back -  v2.7.4",
				},
			},
			sv: semver.Version{
				Major: 2,
				Minor: 7,
				Patch: 4,
			},
		},
	} {
		rs := filterStableReleases(tc.rs)
		if len(rs) < 1 {
			t.Fatalf("Failed to find any stable release")
		}
		if !rs[0].Version().Equals(tc.sv) {
			t.Errorf("Version mismatch: %+v vs %+v", tc.rs, tc.sv)
		}
	}
}
