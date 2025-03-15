package pkg

import (
	"fmt"
	"github.com/2bit-software/gogo/pkg/sh"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// CompileGo compiles a Go binary, assuming filePath contains a main package.
// It will output to [inputFolderPath|outputFolderPath]/binaryName-<GOOS>-<GOARCH>.
func CompileGo(inputFolderPath, outputFolderPath, tags, binaryName, versionPath, strOs, strArch string) error {
	fmt.Printf("Compiling Go binary with arguments: inputFolderPath=%s, outputFolderPath=%s, tags=%s, binaryName=%s, strOs=%s, strArch=%s\n", inputFolderPath, outputFolderPath, tags, binaryName, strOs, strArch)
	wd, err := os.Getwd()
	if err != nil {
		wd = "unknown"
	}
	fmt.Printf("Current directory: %s\n", wd)
	// Get the current HEAD SHA from git
	headSha, err := GetCurrentSha()
	if err != nil {
		return err
	}
	headSha = strings.TrimSpace(headSha)

	// Get the current time
	timeVar := time.Now().Format(time.RFC3339)

	// Get the current user
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	who := currentUser.Username

	// Determine the state of the git repository
	state := GetGitStatus()

	// Set GOOS and GOARCH if provided
	if strOs != "" {
		err = os.Setenv("GOOS", strOs)
		if err != nil {
			return err
		}
	}
	if strArch != "" {
		err = os.Setenv("GOARCH", strArch)
		if err != nil {
			return err
		}
	}

	// Get GOOS and GOARCH, defaulting to the current system's values
	goos := os.Getenv("GOOS")
	if goos == "" {
		goos = runtime.GOOS
	}
	goarch := os.Getenv("GOARCH")
	if goarch == "" {
		goarch = runtime.GOARCH
	}

	// Create the output filename
	filename := fmt.Sprintf("%v", binaryName)
	fmt.Printf("building %s...\n", filename)

	// get the full path to inputFolderPath
	inputFolderPath, err = filepath.Abs(inputFolderPath)
	if err != nil {
		return err
	}

	// if outputFolderPath is not provided, set the output path to the input+dist
	if outputFolderPath == "" {
		outputFolderPath = inputFolderPath
	}

	outputFolderPath, err = filepath.Abs(outputFolderPath)
	if err != nil {
		return err
	}

	// Set CGO_ENABLED environment variable
	err = os.Setenv("CGO_ENABLED", "0")
	if err != nil {
		fmt.Println("could not set CGO_ENABLED=0")
	}

	// Set the location
	location := "main"
	if versionPath != "" {
		location = versionPath
	}

	// Construct ldflags
	ldflags := fmt.Sprintf("-w -s -X '%s.BuildTime=%s' -X '%s.Who=%s' -X '%s.State=%s'",
		location, timeVar, location, who, location, state)

	// Build the command string
	cmdStr := fmt.Sprintf("go build -a -ldflags \"%s\"", ldflags)
	if tags != "" {
		// strip any extraneous quotes and whitespace from the tags
		tags = strings.TrimSpace(strings.ReplaceAll(tags, "\"", ""))
		cmdStr += fmt.Sprintf(" -tags \"%s\"", tags)
	}

	cmdStr += fmt.Sprintf(" -o \"%s\" \"%s\"", filepath.Join(outputFolderPath, filename), inputFolderPath)
	// Execute the build command
	if out, err := sh.Cmd(cmdStr).String(); err != nil {
		return fmt.Errorf("could not build binary: %s: %w", out, err)
	}

	fmt.Println("done")
	return nil
}
