package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WIP: This is in progress and is an optimization

// PathInfo contains the common parent directory and relative paths
type PathInfo struct {
	CommonParent  string
	RelativePaths []string
}

// ParentDirWithRelatives finds the common parent directory and relative paths
func ParentDirWithRelatives(paths []string) (PathInfo, error) {
	switch len(paths) {
	case 0:
		return PathInfo{}, nil
	case 1:
		cleanPath := filepath.Clean(paths[0])
		return PathInfo{
			CommonParent:  filepath.Dir(cleanPath),
			RelativePaths: []string{filepath.Base(cleanPath)},
		}, nil
	}

	// Get absolute paths for all inputs
	absPaths := make([]string, len(paths))
	for i, p := range paths {
		absPath, err := filepath.Abs(p)
		if err != nil {
			return PathInfo{}, fmt.Errorf("failed to get absolute path for %s: %w", p, err)
		}
		absPaths[i] = absPath
	}

	// Split all paths into components
	var pathComponents [][]string
	for _, p := range absPaths {
		// Handle volume/drive name for Windows
		vol := filepath.VolumeName(p)
		// Remove volume name if present and split remaining path
		p = strings.TrimPrefix(p, vol)
		// determine if the path is a file or a folder, if it's a file, just get the directory
		fileInfo, err := os.Stat(p)
		if err == nil && !fileInfo.IsDir() {
			p = filepath.Dir(p)
		}
		components := strings.Split(filepath.Clean(p), string(filepath.Separator))
		// Prepend volume name as first component if present
		if vol != "" {
			components = append([]string{vol}, components...)
		}
		pathComponents = append(pathComponents, components)
	}

	// Find common prefix
	commonComponents := pathComponents[0]
	for _, components := range pathComponents[1:] {
		if len(components) < len(commonComponents) {
			commonComponents = commonComponents[:len(components)]
		}
		for i := 0; i < len(commonComponents); i++ {
			if components[i] != commonComponents[i] {
				commonComponents = commonComponents[:i]
				break
			}
		}
	}

	// Construct common parent path
	commonParent := filepath.Join(commonComponents...)
	if !strings.HasPrefix(commonParent, filepath.VolumeName(commonParent)) {
		commonParent = string(filepath.Separator) + commonParent
	}

	// Calculate relative paths
	relatives := make([]string, len(paths))
	for i, absPath := range absPaths {
		rel, err := filepath.Rel(commonParent, absPath)
		if err != nil {
			return PathInfo{}, fmt.Errorf("failed to get relative path: %w", err)
		}
		relatives[i] = rel
	}

	return PathInfo{
		CommonParent:  commonParent,
		RelativePaths: relatives,
	}, nil
}
