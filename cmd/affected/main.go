package main

import (
	"encoding/json"
	"testing-infra/pkg/utils"
)

type StrategyMatrix struct {
	Packages []string `json:"packages"`
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

	println(string(matrixJson))
}

func affected(config utils.Config, diffs []string) StrategyMatrix {
	// TODO(dcavazos): Detect affected changes more granularly with the diffs.
	// TODO(dcavazos): If '.' (root diffs) in pkgs, return all packages.
	packages := make(map[string]bool)
	for _, diff := range diffs {
		if !config.Matches(diff) {
			continue
		}
		pkg := config.FindPackage(diff)
		packages[pkg] = true
	}

	var pkgs []string
	for pkg := range packages {
		pkgs = append(pkgs, pkg)
	}
	return StrategyMatrix{
		Packages: pkgs,
	}
}
