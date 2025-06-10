package version

import (
	"runtime"
)

var (
	// Version is the version of the CLI, set at build time
	Version = "dev"

	// BuildTime is the time the CLI was built, set at build time
	BuildTime = "unknown"

	// GoVersion is the version of Go used to build the CLI
	GoVersion = runtime.Version()

	// Platform is the platform the CLI is running on
	Platform = runtime.GOOS + "/" + runtime.GOARCH
)
