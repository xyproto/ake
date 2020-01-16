package main

import "strings"

// Command is a shell command, indented with "\t", belonging to a target
type Command struct {
	cmd         string // the command to be run, with variables replaced
	silent      bool   // the "@" prefix
	ignoreError bool   // the "-" prefix
	incomplete  bool   // ends with \ and is continued on the line below
}

// NewCommand interprets a line that starts with "\t" in a Makefile and
// returns a new Command struct.
func NewCommand(line string) *Command {
	trimmed := strings.TrimSpace(line)
	silent := false
	ignoreError := false
	incomplete := false
	if strings.HasPrefix(trimmed, "@") {
		silent = true
		trimmed = trimmed[1:]
	}
	if strings.HasPrefix(trimmed, "-") {
		ignoreError = true
		trimmed = trimmed[1:]
	}
	if strings.HasSuffix(trimmed, "\\") {
		incomplete = true
		trimmed = trimmed[:len(trimmed)-1]
	}
	return &Command{trimmed, silent, ignoreError, incomplete}
}
