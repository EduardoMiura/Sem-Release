package semrelease

import (
	"mime/multipart"

	"github.com/Masterminds/semver"
)

// Repository ..
type Repository struct {
	Owner Owner  `json:"owner"`
	Name  string `json:"name"`
}

// PullRequest ..
type PullRequest struct {
	Merged bool `json:"merged"`
}

// Event ..
type Event struct {
	Repository  Repository   `json:"repository"`
	Action      string       `json:"action"`
	PullRequest *PullRequest `json:"pull_request"`
}

// Owner ..
type Owner struct {
	Login string `json:"login"`
}

// Commit ..
type Commit struct {
	SHA     string
	Raw     []string
	Type    string
	Scope   string
	Message string
	Change  Change
}

// Release ..
type Release struct {
	SHA     string
	Version *semver.Version
}
type Config struct {
	InstalationID string         `json:"InstalationID,omitempty"`
	IntegrationID string         `json:"IntegrationId,omitempty"`
	File          multipart.File `json:"File,omitempty"`
}

// TagRequest ...
type TagRequest struct {
	Tag     *string `json:"tag,omitempty"`
	Message *string `json:"message,omitempty"`
	Object  *string `json:"object,omitempty"`
	Type    *string `json:"type,omitempty"`
}

// Releases ..
type Releases []*Release

// Change ..
type Change struct {
	Major, Minor, Patch bool
}

// Less ..
func (r Releases) Less(i, j int) bool {
	return r[j].Version.LessThan(r[i].Version)
}

// Swap ..
func (r Releases) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
