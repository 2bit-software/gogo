package version

import (
	"fmt"
	"github.com/2bit-software/gogo/pkg/sh"
	"golang.org/x/mod/semver"
	"strings"
)

// MeetsGoVersion checks if the installed Go version is at least the passed version.
// It uses semantic versioning for proper version comparison.
// Returns true if the version is >= 1.24, false otherwise, and an error if something went wrong.
func MeetsGoVersion(required string) (bool, error) {
	// get the go version
	version, err := sh.Cmd("go version").StdOut()
	if err != nil {
		return false, err
	}

	// extract the version string
	versionStr, err := getGoVersionString(version)
	if err != nil {
		return false, err
	}

	// compare the versions
	return meetsGoVersionHelper(required, versionStr)
}

func meetsGoVersionHelper(required, current string) (bool, error) {
	// Compare version using semver library
	compareResult := semver.Compare(required, current)
	return compareResult <= 0, nil
}

func getGoVersionString(version string) (string, error) {
	// Parse the output to extract the version
	versionOutput := string(version)
	// Typical output format: "go version go1.xx.y os/arch"
	parts := strings.Split(versionOutput, " ")
	if len(parts) < 3 {
		return "", fmt.Errorf("unexpected 'go version' output format: %s", versionOutput)
	}

	// Extract version string (e.g., "go1.24.0" -> "1.24.0")
	versionStr := parts[2]
	if !strings.HasPrefix(versionStr, "go") {
		return "", fmt.Errorf("unexpected version string format: %s", versionStr)
	}

	versionStr = strings.TrimPrefix(versionStr, "go")
	return "v" + versionStr, nil
}
