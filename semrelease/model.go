package semrelease

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/semver"
)

// Commit ..
type Commit struct {
	SHA              string     `json:"sha,omitempty"`
	AbbreviatedSHA   string     `json:"abbreviatedSHA,omitempty"`
	Raw              []string   `json:"raw,omitempty"`
	Type             CommitType `json:"type,omitempty"`
	Scope            string     `json:"scope,omitempty"`
	Message          string     `json:"message,omitempty"`
	SanitizedMessage string     `json:"sanitizeMessage,omitempty"`
}

func (c Commit) String() string {
	return fmt.Sprintf("* **%s:** %s (%s)", c.Scope, c.Message, c.AbbreviatedSHA)
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

var (
	patch = CommitPriority{Name: "patch", value: 1}
	minor = CommitPriority{Name: "minor", value: 2}
	major = CommitPriority{Name: "major", value: 3}
)

// CommitPriority ...
type CommitPriority struct {
	Name  string `json:"name,omitempty"`
	value int
}

// CommitType ...
type CommitType struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"descrition,omitempty"`
	priority    CommitPriority
}
