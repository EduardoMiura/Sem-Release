package semrelease

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/Masterminds/semver"
)

// Repo ..
type Repo struct {
	Owner Owner  `json:"owner"`
	Name  string `json:"name"`
}

// PullRequest ..
type PullRequest struct {
	Merged bool `json:"merged"`
}

// Event ..
type Event struct {
	Repository  Repo         `json:"repository"`
	Action      string       `json:"action"`
	PullRequest *PullRequest `json:"pull_request"`
}

// Owner ..
type Owner struct {
	Login string `json:"login"`
}

// Commit ..
type Commit struct {
	SHA              string     `json:"sha,omitempty"`
	AbbreviatedSHA   string     `json:"abbreviatedSHA,omitempty"`
	Raw              []string   `json:"raw,omitempty"`
	Type             CommitType `json:"type,omitempty"`
	Scope            string     `json:"scope,omitempty"`
	Message          string     `json:"message,omitempty"`
	SanitizedMessage string     `json:"sanitizeMessage,omitempty"`
	Change           Change     `json:"change,omitempty"`
}

type changeLog map[string][]string

// Release ..
type Release struct {
	SHA          string
	Version      *semver.Version
	Change       CommitPriority
	ChangeLog    changeLog //map[string][]string
	Branch       string
	IsPreRelease bool
}

func (c changeLog) String() string {
	var b bytes.Buffer
	for key, logMessage := range c {
		commitType := commitTypes[key]
		b.WriteString(fmt.Sprintf("#### %s \n\n %s\n", commitType.Description, strings.Join(logMessage, "\n")))
	}
	return b.String()
}

// Config ...
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

// // Less ..
// func (r Releases) Less(i, j int) bool {
// 	return r[j].Version.LessThan(r[i].Version)
// }

// // Swap ..
// func (r Releases) Swap(i, j int) {
// 	r[i], r[j] = r[j], r[i]
// }

// CommitPriority ...
type CommitPriority struct {
	Name  string `json:"name,omitempty"`
	value int
}

var (
	patch = CommitPriority{Name: "patch", value: 1}
	minor = CommitPriority{Name: "minor", value: 2}
	major = CommitPriority{Name: "major", value: 3}
)

// CommitType ...
type CommitType struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"descrition,omitempty"`
	priority    CommitPriority
}
