// Package build holds version information injected at build time via -ldflags.
//
// Build with:
//   go build -ldflags "-X github.com/SparkssL/Midaz-cli/internal/build.Version=0.2.0" ./cmd/seer-q/
package build

import "runtime"

// Version is the CLI version, set at build time.
var Version = "dev"

// GoVersion returns the Go runtime version.
func GoVersion() string {
	return runtime.Version()
}

// OS returns the operating system.
func OS() string {
	return runtime.GOOS
}

// Arch returns the architecture.
func Arch() string {
	return runtime.GOARCH
}
