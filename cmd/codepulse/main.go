package main

import "codepulse/internal/codepulse"

// filled by ldflags during build
var Version string

func main() {
	codepulse.StartApplication(Version)
}
