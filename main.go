package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rody/pmd-annotations-github-action/pmd"
	"github.com/sethvargo/go-githubactions"
)

var (
	dir              string
	reportfile       string
	minErrorPriority int
	failOnError      bool
)

func main() {
	flag.StringVar(&dir, "dir", "", "")
	flag.StringVar(&reportfile, "reportfile", "", "")
	flag.IntVar(&minErrorPriority, "min-error-priority", 0, "")
	flag.BoolVar(&failOnError, "fail-on-error", false, "")
	flag.Parse()

	githubactions.Debugf("minErrorPriority: %d", minErrorPriority)
	githubactions.Debugf("failOnError: %b", failOnError)

	if reportfile == "" {
		githubactions.Fatalf("missing input 'reportfile'")
	}

	report, err := parseReport(reportfile)
	if err != nil {
		githubactions.Fatalf("could not parse reportfile: %s", err)
	}

	violationCount := 0
	warningCount := 0
	errorCount := 0

	for _, f := range report.Files {
		relPath, err := relpath(f.Filename)
		if err != nil {
			githubactions.Fatalf("could not get path of file ''%s': %s", f.Filename, err)
		}

		githubactions.Group(fmt.Sprintf("%s: violations=%d", relPath, len(f.Violations)))
		for _, v := range f.Violations {
			violationCount += 1
			action := githubactions.WithFieldsMap(map[string]string{
				"file":   relPath,
				"line":   strconv.Itoa(v.BeginLine),
				"column": strconv.Itoa(v.BeginColumn),
				"title":  v.Description,
			})

			msg := fmt.Sprintf("%s (priority: %d)\n%s", v.Rule, v.Priority, v.ExternalInfoUrl)
			if v.Priority > minErrorPriority {
				warningCount += 1
				action.Warningf(msg)
			} else {
				errorCount += 1
				action.Errorf(msg)
			}
		}
		githubactions.EndGroup()
	}

	if violationCount == 0 {
		githubactions.Infof("no problem found")
	} else {
		githubactions.Infof("%d violations(s), errors=%d, warnings=%d", violationCount, errorCount, warningCount)
		if errorCount > 0 && failOnError {
			os.Exit(4)
		}
	}
}

func relpath(file string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	filename := filepath.Join(dir, file)

	if strings.HasPrefix(filename, "/") {
		return filepath.Rel(cwd, filename)
	}

	return filename, nil
}

func parseReport(filename string) (pmd.Report, error) {
	f, err := os.Open(filename)
	if err != nil {
		return pmd.Report{}, err
	}
	defer f.Close()
	return pmd.Parse(f)
}
