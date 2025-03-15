// Copyright (C) 2024  Morgan Stewart Hein
//
// This Source Code Form is subject to the terms
// of the Mozilla Public License, v. 2.0. If a copy
// of the MPL was not distributed with this file, You
// can obtain one at https://mozilla.org/MPL/2.0/.

package gogo

import (
	"fmt"
	"io"
	"runtime/debug"
	"strings"
)

// These attributes should only be changed by the build script, do not change it manually.
var (
	// BuildTime is the timestamp when the binary was built, injected by -ldflags
	BuildTime string = "..."
	// Who identifies the builder of this binary, injected by -ldflags
	Who string = "..."
	// State indicates whether the build was from a clean working directory, injected by -ldflags
	State string = "..."
	// VersionTag is the semantic version of the application, injected by -ldflags
	VersionTag = "0.1.0"
)

// BuildInfo contains comprehensive version and build information
// from both custom ldflags injection and Go's runtime/debug package.
type BuildInfo struct {
	// Manual build information (from ldflags)

	// BuildSha is the Git commit SHA, injected by -ldflags during build
	BuildSha string
	// BuildTime is the timestamp when the binary was built, injected by -ldflags
	BuildTime string
	// Who identifies the builder of this binary, injected by -ldflags
	Who string
	// State indicates whether the build was from a clean working directory, injected by -ldflags
	State string
	// Version is the semantic version of the application, injected by -ldflags
	Version string

	// Go runtime information (from debug.BuildInfo)

	// GoVersion is the Go version used to build the binary (e.g., "go1.18.3")
	GoVersion string
	// Path is the main module path
	Path string
	// Main contains information about the main module
	Main Module
	// Dependencies contains information about all dependencies
	Dependencies []Module

	// Build settings (from debug.BuildInfo.Settings)

	// BuildMode indicates how the binary was built (e.g., "exe")
	BuildMode string
	// Compiler indicates which compiler was used (e.g., "gc")
	Compiler string
	// CGOEnabled indicates whether CGO was enabled during build
	CGOEnabled string
	// GOOS is the target operating system (e.g., "darwin", "linux")
	GOOS string
	// GOARCH is the target architecture (e.g., "amd64", "arm64")
	GOARCH string
	// GOARMVersion specifies the ARM version if GOARCH is "arm" or "arm64"
	GOARMVersion string

	// Version control information (from debug.BuildInfo.Settings)

	// VCS identifies the version control system (e.g., "git")
	VCS string
	// VCSRevision is the full commit hash from version control
	VCSRevision string
	// VCSTime is the commit timestamp from version control
	VCSTime string
	// VCSModified indicates whether the working directory had uncommitted changes
	VCSModified string

	// Raw build settings for any custom or additional settings
	RawSettings []BuildSetting
}

// Module represents a Go module with its version information
type Module struct {
	// Path is the import path of the module
	Path string
	// Version is the module version
	Version string
	// Sum is the checksum of the module
	Sum string
	// Replace points to a replacement if this module is replaced
	Replace *Module
}

// BuildSetting represents an individual build setting key-value pair
type BuildSetting struct {
	// Key is the name of the build setting
	Key string
	// Value is the value of the build setting
	Value string
}

// returns the version in the format "[vX.Y.Z:]commit"
func Version() string {
	versionString := ""
	b := GetBuildInfo()
	if b.Version != "" {
		versionString = fmt.Sprintf("v%v", b.Version)
	}
	if b.VCSRevision != "" {
		if versionString != "" {
			versionString += ":" // separate version and commit with colon when both exist
		}
		versionString += b.VCSRevision[:7]
	}

	return versionString
}

// GetBuildInfo returns comprehensive version information by combining
// manually injected build variables with Go's runtime build information
func GetBuildInfo() BuildInfo {
	info := BuildInfo{
		// Manual build info
		BuildTime: BuildTime,
		Who:       Who,
		State:     State,
		Version:   VersionTag,
	}

	// Get Go's built-in build information
	if bi, ok := debug.ReadBuildInfo(); ok {
		info.GoVersion = bi.GoVersion
		info.Path = bi.Path

		// Copy main module info
		info.Main = Module{
			Path:    bi.Main.Path,
			Version: bi.Main.Version,
			Sum:     bi.Main.Sum,
		}

		// Copy dependencies
		for _, dep := range bi.Deps {
			module := Module{
				Path:    dep.Path,
				Version: dep.Version,
				Sum:     dep.Sum,
			}

			if dep.Replace != nil {
				module.Replace = &Module{
					Path:    dep.Replace.Path,
					Version: dep.Replace.Version,
					Sum:     dep.Replace.Sum,
				}
			}

			info.Dependencies = append(info.Dependencies, module)
		}

		// Process build settings
		for _, setting := range bi.Settings {
			// Store raw setting
			info.RawSettings = append(info.RawSettings, BuildSetting{
				Key:   setting.Key,
				Value: setting.Value,
			})

			// Extract specific settings into dedicated fields
			switch setting.Key {
			case "-buildmode":
				info.BuildMode = setting.Value
			case "-compiler":
				info.Compiler = setting.Value
			case "CGO_ENABLED":
				info.CGOEnabled = setting.Value
			case "GOOS":
				info.GOOS = setting.Value
			case "GOARCH":
				info.GOARCH = setting.Value
			case "GOARM64", "GOARM":
				info.GOARMVersion = setting.Value
			case "vcs":
				info.VCS = setting.Value
			case "vcs.revision":
				info.VCSRevision = setting.Value
			case "vcs.time":
				info.VCSTime = setting.Value
			case "vcs.modified":
				info.VCSModified = setting.Value
			}
		}
	}

	return info
}

// PrettyPrint formats build information in a human-readable way for CLI applications
// It returns a string containing formatted build information with proper indentation and grouping
func (info BuildInfo) PrettyPrint() string {
	var output strings.Builder

	// Format app version information first - the most important details for users
	output.WriteString("Application:\n")
	output.WriteString(fmt.Sprintf("  Version:    %s\n", info.Version))
	output.WriteString(fmt.Sprintf("  Build Time: %s\n", info.BuildTime))
	output.WriteString(fmt.Sprintf("  Builder:    %s\n", info.Who))

	// Version control information
	output.WriteString("\nVersion Control:\n")
	output.WriteString(fmt.Sprintf("  System:      %s\n", info.VCS))
	output.WriteString(fmt.Sprintf("  Commit:      %s\n", info.VCSRevision))
	output.WriteString(fmt.Sprintf("  Commit Time: %s\n", info.VCSTime))
	modified := "No"
	if info.VCSModified == "true" {
		modified = "Yes"
	}
	output.WriteString(fmt.Sprintf("  Modified:    %s\n", modified))

	// Build environment
	output.WriteString("\nBuild Environment:\n")
	output.WriteString(fmt.Sprintf("  Go Version: %s\n", info.GoVersion))
	output.WriteString(fmt.Sprintf("  OS/Arch:    %s/%s\n", info.GOOS, info.GOARCH))
	if info.GOARMVersion != "" {
		output.WriteString(fmt.Sprintf("  ARM Version: %s\n", info.GOARMVersion))
	}
	output.WriteString(fmt.Sprintf("  CGO Enabled: %s\n", info.CGOEnabled))
	output.WriteString(fmt.Sprintf("  Compiler:    %s\n", info.Compiler))
	output.WriteString(fmt.Sprintf("  Build Mode:  %s\n", info.BuildMode))

	// Module information (condensed)
	output.WriteString("\nMain Module:\n")
	output.WriteString(fmt.Sprintf("  Path:    %s\n", info.Main.Path))
	output.WriteString(fmt.Sprintf("  Version: %s\n", info.Main.Version))

	// For verbose output, can be conditionally disabled if too much info
	if len(info.Dependencies) > 0 {
		output.WriteString(fmt.Sprintf("\nDependencies: (%d total)\n", len(info.Dependencies)))
		// List just some key dependencies to avoid flooding the output
		for i, dep := range info.Dependencies {
			if i >= 5 {
				output.WriteString(fmt.Sprintf("  ... and %d more\n", len(info.Dependencies)-5))
				break
			}
			output.WriteString(fmt.Sprintf("  %s@%s\n", dep.Path, dep.Version))
		}
	}

	return output.String()
}

// PrintVersion prints the build information to the specified writer
// This is useful for CLI applications where you want to output version info
// to stdout, stderr, or a log file
func PrintVersion(w io.Writer) {
	fmt.Fprintln(w, GetBuildInfo().PrettyPrint())
}
