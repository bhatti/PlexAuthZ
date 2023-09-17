package main

import (
	cmd "github.com/bhatti/PlexAuthZ/cmd/commands"
)

var (
	// Version of the application.
	Version = "dev"
	// Commit for latest change.
	Commit = "dirty"
	// Date for latest change.
	Date = ""
)

func main() {
	if err := cmd.Execute(Version, Commit, Date); err != nil {
		panic(err)
	}
}
