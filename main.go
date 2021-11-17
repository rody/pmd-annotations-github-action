package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rody/pmd-annotations-github-action/pmd"
	"github.com/sethvargo/go-githubactions"
)

var (
	dir string
	reportfile string
	minErrorPriority int
)

func main() {
	flag.StringVar(&dir, "dir", "", "")
	flag.StringVar(&reportfile, "reportfile", "", "")
	flag.IntVar(&minErrorPriority, "min-error-priority", 1, "")

	flag.Parse()

	if reportfile == "" {
		githubactions.Fatalf("missing input 'reportfile'")
	}

	report, err := parseReport(reportfile)
	if err != nil {
		githubactions.Fatalf("could not parse reportfile: %s", err)
	}

	violationCount := 0
	for _, f := range report.Files {
		relPath, err := relpath(f.Filename)
		if err != nil {
			githubactions.Fatalf("could not get path of file ''%s': %s", f.Filename, err)
		}

		githubactions.Group(fmt.Sprintf("PMD %s: errors=%d", relPath, len(f.Violations)))
		for _, v := range f.Violations {
			violationCount += 1
			annotation := fmt.Sprintf("file=%s,line=%d::[%s] %s", relPath, v.BeginLine, v.Rule, v.Description)
			if (v.Priority > minErrorPriority) {
				githubactions.Warningf(annotation)
			} else {
				githubactions.Errorf(annotation)
			}
		}
		githubactions.EndGroup()
	}

	if violationCount == 0 {
		githubactions.Infof("no problem found")
	} else {
		githubactions.Infof("%d problem(s)", violationCount)
	}
}

func relpath(file string) (string, error) {
	cwd, err := os.Getwd();
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
