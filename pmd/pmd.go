package pmd

import (
	"encoding/json"
	"io"
)

type Report struct {
	FormatVersion int8   `json:"formatVersion"`
	PMDVersion    string `json:"pmdVersion"`
	Timestamp     string `json:"timestamp"`
	Files         []File `json:"files"`
}

type File struct {
	Filename   string      `json:"filename"`
	Violations []Violation `json:"violations"`
}

type Violation struct {
	BeginLine       int32  `json:"beginLine"`
	BeginColumn     int32  `json:"beginColumn"`
	EndLine         int32  `json:"endLine"`
	EndColumn       int32  `json:"endColumn"`
	Description     string `json:"description"`
	Rule            string `json:"rule"`
	Ruleset         string `json:"ruleset"`
	Priority        int   `json:"priority"`
	ExternalInfoUrl string `json:"externalInfoUrl"`
}

// Parse reads a PMD report in JSON format
// and returns a Report.
func Parse(pmdstring io.Reader) (Report, error) {
	report := Report{}
	err := json.NewDecoder(pmdstring).Decode(&report)
	if err != nil {
		return report, err
	}
	return report, nil
}
