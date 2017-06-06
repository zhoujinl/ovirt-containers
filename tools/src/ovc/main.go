/*
Copyright (c) 2017 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

// This tool loads and builds all the images.

import (
	"fmt"
	"os"
	"path/filepath"

	"ovc/build"
	"ovc/log"
)

// ToolFunc is the type of functions that implement tools.
//
type ToolFunc func(project *build.Project) error

// This index contains the mapping from names to tool functions.
//
var toolsIndex = map[string]ToolFunc{
	"build":  buildTool,
	"clean":  cleanTool,
	"deploy": deployTool,
	"push":   pushTool,
	"save":   saveTool,
}

func main() {
	// Get the name of the tool:
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s TOOL [ARGS...]\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	toolName := os.Args[1]

	// Find the function that corresponds to the tool name:
	toolFunc := toolsIndex[toolName]
	if toolFunc == nil {
		fmt.Fprintf(os.Stderr, "Can't find tool named '%s'.\n", toolName)
		os.Exit(1)
	}

	// Open the log:
	log.Open(toolName)
	log.Info("Log file is '%s'", log.Path())

	// Load the project:
	path, _ := filepath.Abs("project.conf")
	log.Info("Loading project file '%s'", path)
	project, err := build.LoadProject(path)
	if err != nil {
		log.Error("%s", err)
		os.Exit(1)
	}

	// Call the tool function, and close the project regardless of
	// the result:
	log.Debug("Running tool '%s'", toolName)
	err = toolFunc(project)
	project.Close()

	// Report the result:
	var code int
	if err != nil {
		log.Error("%s", err)
		log.Error("Tool failed, check log file '%s' for details", log.Path())
		code = 1
	} else {
		log.Info("Tool finished successfully")
		code = 0
	}
	log.Debug("Exit code is %d", code)

	// Close the log and exit:
	log.Close()
	os.Exit(code)
}
