package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/davidcavazos/testing-infra/pkg/utils"
)

type Job struct {
	Package string `json:"package"`
}

func main() {
	usage := "usage: affected path/to/config.json [head-commit] [main-commit]"
	configPath := utils.ArgRequired(1, usage)
	headCommit := utils.ArgWithDefault(2, "HEAD")
	mainCommit := utils.ArgWithDefault(3, "origin/main")

	config, err := utils.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	diffs, err := utils.Diffs(headCommit, mainCommit)
	if err != nil {
		panic(err)
	}

	matrix := affected(config, diffs)
	matrixJson, err := json.Marshal(matrix)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(matrixJson))
}

func affected(config utils.Config, diffs []string) []Job {
	// TODO(dcavazos): Detect affected changes more granularly with the diffs.
	packages := make(map[string]bool)
	for _, diff := range diffs {
		if !config.Matches(diff) {
			continue
		}
		pkg := config.FindPackage(diff)
		packages[pkg] = true
	}

	if packages["."] {
		var jobs []Job
		err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if path == ".git" {
				return filepath.SkipDir
			}
			if d.IsDir() && config.IsPackageDir(path) {
				jobs = append(jobs, Job{Package: path})
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
		return jobs
	}

	jobs := make([]Job, 0, len(packages))
	for pkg := range packages {
		jobs = append(jobs, Job{Package: pkg})
	}
	return jobs
}
