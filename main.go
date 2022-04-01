package main

import (
	"kore-on/cmd"
	"kore-on/pkg/conf"
)
 
var Version = "unknown_version"
var CommitId = "unknown_commitid"
var BuildDate = "unknown_builddate"

func main() {
	conf.Version = Version
	conf.CommitId = CommitId
	conf.BuildDate = BuildDate

	cmd.Execute()
}
