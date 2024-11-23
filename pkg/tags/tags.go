// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.
package tags

import (
	"fmt"
	"go/build/constraint"
	"strings"
)

// this file is for reading, finding, and updating build tags in source code

// HasBuildTag checks if the src contains any of the provided build tags.
func HasBuildTag(src string, tags []string) bool {
	if len(tags) == 0 {
		return true
	}
	lines := strings.Split(src, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//go:build") {
			tagLine := strings.TrimSpace(line[len("//go:build"):])
			tagExpr := strings.ReplaceAll(tagLine, " ", ",") // Convert space to comma for evaluation
			if matchesBuildTag(tagExpr, tags) {
				return true
			}
		}
		// Check for // +build (old format)
		if strings.HasPrefix(line, "// +build") {
			tagLine := strings.TrimSpace(line[len("// +build"):])
			tagExpr := strings.ReplaceAll(tagLine, " ", ",") // Convert space to comma for evaluation
			if matchesBuildTag(tagExpr, tags) {
				return true
			}
		}
	}
	return false
}

// matchesBuildTag parses the build expression and checks the given tags.
func matchesBuildTag(buildExpr string, tags []string) bool {
	expr, err := constraint.Parse(buildExpr) // Parse the build expression using go/build/constraint
	if err != nil {
		fmt.Println("Error parsing build expression:", err)
		return false
	}

	// Loop through the provided tags and check if any match the constraint expression.
	for _, tag := range tags {
		if expr.Eval(func(name string) bool { return name == tag }) {
			return true
		}
	}

	return false
}

// AddTag searches for build tags in the Go source (src) and adds
// the provided tag (tag) if it doesn't already exist. The function
// looks for both `//go:build` (new format) and `// +build` (old format)
// and ensures that the new tag is added under the appropriate format.
// If neither build tag format exists, it adds a `//go:build` directive at
// the top of the file. If the tag already exists, the original source is returned.
//
// The function doesn't attempt to rewrite the whole file but operates on the
// source text to add tags directly.
func AddTag(src string, tag string) string {
	// first check if the tag already exists
	if HasBuildTag(src, []string{tag}) {
		return src
	}

	lines := strings.Split(src, "\n")
	hasGoBuild := false
	hasPlusBuild := false
	// Track which lines are `//go:build` and `// +build`
	goBuildLineIdx := -1
	plusBuildLineIdx := -1

	// Variables to store new tag expressions
	goBuildExpr := tag
	plusBuildExpr := tag

	// Manually search the source code for existing build tags
	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Handle `//go:build` tags (new format)
		if strings.HasPrefix(line, "//go:build") {
			hasGoBuild = true
			goBuildLineIdx = i

			existingTags := strings.TrimSpace(line[len("//go:build"):])

			// Append the new tag to the existing expression
			goBuildExpr = existingTags + " && " + tag
		}

		// Handle `// +build` tags (old format)
		if strings.HasPrefix(line, "// +build") {
			hasPlusBuild = true
			plusBuildLineIdx = i

			existingTags := strings.TrimSpace(line[len("// +build"):])

			// Append the new tag to the existing expression
			plusBuildExpr = existingTags + " " + tag
		}

		// Stop looking after we pass comments and find actual code (package declaration, imports, etc).
		if !strings.HasPrefix(line, "//") && line != "" {
			break
		}
	}

	// Modify the source lines based on what was found
	if hasGoBuild {
		// Update the `//go:build` line with the new tag.
		lines[goBuildLineIdx] = "//go:build " + goBuildExpr
	}

	if hasPlusBuild {
		// Update the `// +build` line with the new tag.
		lines[plusBuildLineIdx] = "// +build " + plusBuildExpr
	}

	// If neither `//go:build` nor `// +build` was found, insert the new `//go:build` at the top.
	if !hasGoBuild && !hasPlusBuild {
		// Add the new build tags in the new format
		lines = append([]string{
			"//go:build " + tag,
			""}, lines...) // Add an empty line for proper formatting
	}

	// Return the modified source text.
	return strings.Join(lines, "\n")
}
