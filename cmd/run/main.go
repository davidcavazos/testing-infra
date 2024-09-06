package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing-infra/pkg/utils"
)

func main() {
	usage := "usage: run path/to/config.json [PACKAGE] [ACTION] [OPTIONS..]"
	configPath := utils.ArgRequired(1, usage)
	packagePath := utils.ArgRequired(2, usage)
	action := utils.ArgRequired(3, usage)
	options := os.Args[4:]

	config, err := utils.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	steps, ok := config.Actions[action]
	if !ok {
		panic("unknown action '" + action + "'")
	}

	for _, step := range steps {
		args := utils.InterpolateArgs(step.Args, options)
		fmt.Printf(">> %v %v\n", step.Command, args)
		cmd := exec.Command(step.Command, args...)
		cmd.Dir = packagePath
		output, err := cmd.CombinedOutput()
		println(string(output))
		if err != nil {
			panic(err)
		}
	}
}
