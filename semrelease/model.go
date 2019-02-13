package semrelease

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

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
	SHA             string
	Version         *semver.Version
	PreviousVersion *semver.Version
	Change          CommitPriority
	ChangeLog       changeLog
	Branch          string
	IsPreRelease    bool
	Repository      string
	Owner           string
}

func (r Release) getReleaseNote() string {
	var b bytes.Buffer
	version := r.Version.String()

	if r.PreviousVersion != nil {
		b.WriteString(fmt.Sprintf("## [%s](https://github.com/%s/%s/compare/v%s...v%s) (%s)\n\n", version, r.Owner, r.Repository, r.PreviousVersion.String(), version, time.Now().UTC().Format("2006-01-02")))
	} else {
		b.WriteString(fmt.Sprintf("## %s (%s)\n\n", version, time.Now().UTC().Format("2006-01-02")))
	}
	b.WriteString(r.ChangeLog.String())
	return b.String()
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
