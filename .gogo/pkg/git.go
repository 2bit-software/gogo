// Copyright (C) 2024  Morgan Stewart Hein
//
// This Source Code Form is subject to the terms
// of the Mozilla Public License, v. 2.0. If a copy
// of the MPL was not distributed with this file, You
// can obtain one at https://mozilla.org/MPL/2.0/.

package pkg

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/2bit-software/gogo/pkg/sh"
)

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
	headSha, err := sh.Cmd("git", "rev-parse", "--verify", "HEAD").StdOut()
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
	statusOutput, err := sh.Cmd("git", "status", "-s").StdOut()
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
