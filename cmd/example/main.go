package main

import "github.com/z3f1r/grpc-ecosystems-example/cmd"

// GitCommit will contain current commit hash
var gitCommit string

func main() {
	cmd.Execute(gitCommit)
}
