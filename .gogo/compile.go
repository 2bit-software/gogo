// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package main

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/2bit-software/gogo"
	"github.com/2bit-software/gogo/pkg/cmdr"
)

// PrintShortSha prints the sha of the current/latest commit using only the first 6 characters
func PrintShortSha() {
	sha, err := GetCurrentShortSha()
	if err != nil {
		fmt.Println("could not get current SHA")
		return
	}
	fmt.Print(sha)
}

// GetCurrentShortSha returns the first 6 characters of the current SHA.
func GetCurrentShortSha() (string, error) {
	sha, err := GetCurrentSha()
	if err != nil {
		return "", err
	}
	return sha[:6], nil
}

// GetCurrentSha returns the current SHA of the git repository. It first tries using the git cli tool,
// then falls back to reading the .git/HEAD file.
func GetCurrentSha() (string, error) {
	headSha, err := cmdr.New("git", "rev-parse", "--verify", "HEAD").StdOut()
	if err == nil {
		return strings.TrimSpace(headSha), nil
	}
	// try the fallback method
	sha, err := tryReadingShaFromGitHead()
	if err != nil {
		return "", fmt.Errorf("could not get current SHA: %w", err)
	}
	fmt.Printf("Detected SHA: %s\n", sha)
	return sha, nil
}

// tryReadingShaFromGitHead attempts to read the current SHA from the .git/HEAD file
// It recursively walks up the tree to first find the .git folder, then reads the SHA from the reference file.
// It errors out if it reaches the root and can't find the .git folder.
// This is basically equivalent to the bash: "cat .git/HEAD | awk '{print ".git/"$2}' | xargs cat"
func tryReadingShaFromGitHead() (string, error) {
	gitFolder, err := findGitFolder(".")
	if err != nil {
		return "", fmt.Errorf("could not find .git folder: %w", err)
	}

	// Read the contents of .git/HEAD
	headContent, err := os.ReadFile(fmt.Sprintf("%v/HEAD", gitFolder))
	if err != nil {
		return "", fmt.Errorf("could not read .git/HEAD: %w", err)
	}

	// Extract the reference path
	refPath := strings.TrimSpace(strings.Split(string(headContent), " ")[1])
	refPath = ".git/" + refPath

	// Read the contents of the reference file
	refContent, err := os.ReadFile(path.Join(path.Dir(gitFolder), refPath))
	if err != nil {
		return "", fmt.Errorf("could not read reference file: %w", err)
	}
	return strings.TrimSpace(string(refContent)), nil
}

func findGitFolder(dir string) (string, error) {
	// get the absolute path to the dir
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	gitFolder := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitFolder); err == nil {
		fmt.Printf("Found git folder: %s\n", gitFolder)
		return gitFolder, nil
	}
	if isRootPath(absPath) {
		return "", fmt.Errorf("could not find .git folder")
	}
	return findGitFolder(filepath.Dir(absPath))
}

// isRootPath checks if the given path is the root directory.
func isRootPath(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	if runtime.GOOS == "windows" {
		// On Windows, the root path is a drive letter followed by a colon and a backslash (e.g., C:\)
		return len(absPath) == 3 && absPath[1] == ':' && absPath[2] == '\\'
	}

	// On Unix-like systems, the root path is just a single forward slash (/)
	return absPath == "/"
}

// GetGitStatus returns the state of the git repository, either "Clean" or "Dirty".
func GetGitStatus() string {
	state := "Unknown"
	statusOutput, err := cmdr.New("git", "status", "-s").StdOut()
	if err != nil {
		return state
	}
	if strings.TrimSpace(statusOutput) == "" {
		state = "Clean"
	} else {
		state = "Dirty"
	}
	return state
}

// CompileGo compiles a Go binary, assuming filePath contains a main package.
// It will output to [inputFolderPath|outputFolderPath]/binaryName-<GOOS>-<GOARCH>.
func CompileGo(inputFolderPath, outputFolderPath, tags, binaryName, strOs, strArch string) error {
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
	shortSha := headSha
	if len(headSha) >= 6 {
		shortSha = headSha[:6]
	}

	// If REF environment variable is set, overwrite shortSha
	ref := os.Getenv("REF")
	if ref != "" {
		fmt.Printf("Overwriting SHORTSHA with REF from ENV: %s\n", ref)
		shortSha = ref
	}

	// Get the current time
	timeVar := time.Now().Format(time.RFC3339)

	// Get the current user
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	who := currentUser.Username

	// Set the location
	location := "main"

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
	filename := fmt.Sprintf("%v-%s-%s", binaryName, goos, goarch)
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

	// Construct ldflags
	ldflags := fmt.Sprintf("-w -s -X '%s.BuildSha=%s' -X '%s.BuildTime=%s' -X '%s.Who=%s' -X '%s.State=%s'",
		location, shortSha, location, timeVar, location, who, location, state)

	// Build the command string
	cmdStr := fmt.Sprintf("go build -a -ldflags \"%s\"", ldflags)
	if tags != "" {
		// strip any extraneous quotes and whitespace from the tags
		tags = strings.TrimSpace(strings.ReplaceAll(tags, "\"", ""))
		cmdStr += fmt.Sprintf(" -tags \"%s\"", tags)
	}

	cmdStr += fmt.Sprintf(" -o \"%s\" \"%s\"", filepath.Join(outputFolderPath, filename), inputFolderPath)
	// Execute the build command
	if out, err := cmdr.New(cmdStr).String(); err != nil {
		return fmt.Errorf("could not build binary: %s: %w", out, err)
	}

	fmt.Println("done")
	return nil
}

func PrintShaCtx(ctx gogo.Context) {
	ctx.SetLongDescription("Prints the SHA of the current/latest commit using only the first 6 characters, this test uses the gogo context").Example("gogo PrintShaCtx")
	PrintShortSha()
}

// Hello prints a greeting to the console
func Hello(input string, count int) {
	fmt.Printf("Hello, %v. Can you count to %v?", input, count)
}
