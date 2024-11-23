package fs

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

type FS interface {
	Open(name string) (fs.File, error)
}

// CompareTimes determines if the target file is newer than the source files.
// Source files may be glob patterns with shell wildcards and expansion.
func CompareTimes(sources []string, target string) (bool, error) {
	if len(sources) == 0 {
		return false, fmt.Errorf("no source files provided")
	}
	if target == "" {
		return false, fmt.Errorf("no target file provided")
	}
	// expand and get abs paths
	sources = expand(sources)
	target = os.ExpandEnv(target)
	for i, p := range sources {
		absPath, err := filepath.Abs(p)
		if err != nil {
			return false, fmt.Errorf("failed to get absolute path for %s: %w", p, err)
		}
		sources[i] = absPath
	}
	targetAbsPath, err := filepath.Abs(target)
	if err != nil {
		return false, fmt.Errorf("failed to get absolute path for %s: %w", target, err)
	}

	return compareTimes(sources, targetAbsPath)
}

// GlobMany performs a glob search in many locations.
func GlobMany(locations, patterns []string) ([]string, error) {
	var matches []string
	for _, location := range locations {
		m, err := Glob(location, patterns)
		if err != nil {
			return matches, err
		}
		matches = append(matches, m...)
	}
	return matches, nil
}

func Glob(location string, patterns []string) ([]string, error) {
	return GlobFS(os.DirFS(location), location, patterns)
}

// GlobFS searches many patterns in a single path/FS
func GlobFS(fsys FS, prefix string, patterns []string) ([]string, error) {
	var matches []string
	for _, p := range patterns {
		m, err := fs.Glob(fsys, p)
		if err != nil {
			return matches, err
		}
		for _, match := range m {
			matches = append(matches, path.Join(prefix, match))
		}
	}
	return matches, nil
}

func expand(patterns []string) []string {
	for i, pattern := range patterns {
		// expand pattern
		patterns[i] = os.ExpandEnv(pattern)
	}
	return patterns
}

func compareTimes(sources []string, target string) (bool, error) {
	// get the target file info
	targetInfo, err := os.Stat(target)
	if err != nil {
		return false, fmt.Errorf("failed to get target file info: %w", err)
	}
	// now for each file in the sources, check if it's newer than the target
	for _, source := range sources {
		sourceInfo, err := os.Stat(source)
		if err != nil {
			return false, fmt.Errorf("failed to get source file info: %w", err)
		}
		if sourceInfo.ModTime().After(targetInfo.ModTime()) {
			return true, nil
		}
	}
	return false, nil
}
