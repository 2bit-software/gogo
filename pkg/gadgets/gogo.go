// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package gadgets

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fatih/color"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/wordwrap"

	"github.com/2bit-software/gogo/pkg/fs"
	"github.com/2bit-software/gogo/pkg/sh"
)

const MAIN_FILENAME = "main.gogo.go"

var (
	//go:embed templates/*
	templates   embed.FS
	gogoFolders = []string{".gogo", "gogofiles", "magefiles"}
	gogoTags    = []string{"gogo", "mage"}
)

// TODO: I think this run function needs a cache of parsed files, passed throughout
//  so that we don't have to reparse the files every time we run a function

// Run this is a simplified version of the Run function in cmd/gogo/main.go
// For now it only searches for the local gogo files, and does not try to
// determine if the function exists in the global cache.
func Run(opts RunOpts, args []string) error {
	debug := opts.GetLogger()
	debug.Printf("Running with %+v\n", opts)
	// detect if we're requesting to build the local cache
	if opts.BuildLocalCache {
		return BuildLocal(opts)
	}
	if opts.BuildGlobalCache {
		opts.SourceDir = opts.GlobalSourceDir
		opts.OutputDir = opts.GlobalBinDir
		return Build(debug, opts.BuildOpts)
	}

	// determine the outputFilePath if not provided
	if opts.OutputDir == "" {
		opts.OutputDir = "/tmp"
	}
	if len(args) == 0 {
		return fmt.Errorf("no function provided")
	}

	debug.Printf("Running function: %s\n", args[0])
	funcToRun := args[0]
	// search for gogo files to run in local namespaces
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	gogoFile, found, err := findLocalFunc(cwd, funcToRun, gogoTags, gogoFolders)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("function %s not found", funcToRun)
	}
	gogoFolder := path.Dir(gogoFile)

	opts.BinaryFilepath, err = getBinaryFilepath(opts)
	if err != nil {
		return err
	}

	opts.SourceDir = gogoFolder

	err = getBuiltBinary(debug, opts.BuildOpts)
	if err != nil {
		return err
	}
	debug.Printf("Running built binary: %s with args %v\n", opts.BinaryFilepath, args)
	// run the binary with the desire target func and arguments, unless it exists in the cache
	ex := sh.Cmd(opts.BinaryFilepath).SetArgs(args...)
	if opts.Verbose {
		ex = ex.SetPrintFinalCommand(true)
	}
	err = ex.RunAndStream()
	if err != nil {
		errString := err.Error()
		if errString == "exit status 1" {
			errString = ""
		}
		return fmt.Errorf("%v", errString)
	}

	// if found, we should check the global cache bin, rebuild if necessary, and run the binary
	return nil
}

// BuildLocal searches for the local gogo files, and builds the binary
func BuildLocal(opts RunOpts) error {
	debug := opts.GetLogger()
	debug.Println("Building local cache...")
	var err error
	opts.BinaryFilepath, err = getBinaryFilepath(opts)
	debug.Println("Output file path:", opts.BinaryFilepath)
	if err != nil {
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	gogoFiles, err := findLocalFiles(cwd, gogoFolders)
	if err != nil {
		return err
	}
	gogoFolder := path.Dir(gogoFiles[0])
	opts.SourceDir = gogoFolder
	return Build(debug, opts.BuildOpts)
}

func BuildGlobal(opts RunOpts) error {
	return nil
}

// ShowFuncList lists all the available functions in the local and global namespaces
func ShowFuncList(opts RunOpts) (int, error) {
	wd, err := os.Getwd()
	if err != nil {
		return 0, err
	}
	// then we are listing the available functions
	funcList, err := BuildFuncList(opts, wd)
	if err != nil {
		return 0, err
	}
	// sort the functions
	slices.SortFunc(funcList, func(e, e2 function) int {
		return strings.Compare(e.Name, e2.Name)
	})
	printFuncList(generateFuncListOutput(funcList, opts.ScreenWidth))
	return len(funcList), nil
}

// BuildFuncList builds a list of functions that can be run. It combines
// both local and global functions. If there are name collisions, the local one
// takes precedence, and the global one can be used with a prefix.
// e.g. `gogo g:funcName` would run the global function `funcName`
func BuildFuncList(opts RunOpts, dir string) ([]function, error) {
	if dir == "" {
		dir = "."
	}
	localFiles, err := listLocalFiles(dir, opts)
	if err != nil {
		return nil, err
	}
	globalFuncs, err := listGlobalFuncs(opts)
	if err != nil {
		return nil, err
	}
	// merge them
	files := append(localFiles, globalFuncs...)
	return parseAll(files)
}

func getBinaryFilepath(opts RunOpts) (string, error) {
	if opts.BinaryFilepath != "" {
		return opts.BinaryFilepath, nil
	}
	// generate filename for this binary
	// get the name of the current directory
	dirName := path.Base(opts.OriginalWorkingDir)
	// hash the directory name
	hashedDirName, err := hashString(dirName)
	if err != nil {
		return "", fmt.Errorf("failed to hash directory name: %w", err)
	}
	filename := fmt.Sprintf("%v-%v", dirName, hashedDirName)
	opts.GetLogger().Printf("Building binary in: %v with filename:%v\n", opts.OutputDir, filename)
	return filepath.Join(opts.OutputDir, filename), nil
}

// printFuncList formats the output and prints it to the console
func printFuncList(output []string) {
	for _, line := range output {
		if line == "" {
			continue
		}
		fmt.Println(line)
	}
}

// generateFuncListOutput creates a formatted list of functions that can be run based on
// an already parsed set of functions.
func generateFuncListOutput(funcs []function, width int) []string {
	// Leave some margin
	width = width - 4

	// Find max name length
	maxNameLen := 0
	for _, f := range funcs {
		if len(f.Name) > maxNameLen {
			maxNameLen = len(f.Name)
		}
	}

	// Calculate description width
	descWidth := width - maxNameLen - 4
	if descWidth < 20 {
		descWidth = 20 // Minimum description width
	}

	// Create color schemes for alternating rows
	evenRow := color.New(color.FgHiWhite)
	//evenRow := color.Cmd(color.BgHiBlack, color.FgWhite)
	//oddRow := color.Cmd(color.FgHiWhite)
	oddRow := color.New(color.FgHiBlue)

	var lines []string
	for i, f := range funcs {
		// Choose color based on row index
		rowColor := evenRow
		if i%2 == 1 {
			rowColor = oddRow
		}

		// Get description or comment
		var description string
		if f.Description != "" {
			description = strings.ReplaceAll(f.Description, "\n", " ")
		} else if f.Comment != "" {
			description = strings.ReplaceAll(f.Comment, "\n", " ")
		} else {
			description = "-"
		}

		if description == "-" {
			lines = append(lines, rowColor.Sprintf("%-*s  %s", maxNameLen, f.Name, description))
			continue
		}

		// Wrap the description
		wrapped := wordwrap.String(description, descWidth)
		wrappedLines := strings.Split(wrapped, "\n")

		// First line with function name
		lines = append(lines, rowColor.Sprintf("%-*s  %s", maxNameLen, f.Name, wrappedLines[0]))

		// Subsequent lines indented (using same color)
		for _, line := range wrappedLines[1:] {
			indented := indent.String(line, uint(maxNameLen+2))
			lines = append(lines, rowColor.Sprint(indented))
		}
	}

	return lines
}

func listLocalFiles(cwd string, opts RunOpts) ([]string, error) {
	// this returns a list of files that match our local search
	localFiles, err := findLocalFiles(cwd, gogoFolders)
	if err != nil {
		return nil, err
	}
	if len(localFiles) == 0 {
		// convert each
		return nil, nil
	}
	// then parse all the funcs, and determine their information
	return localFiles, nil
}

// findLocalFiles searches for the gogo files directory in the given directory.
// If it finds it, it returns a list of all the .go files in that directory.
// If it does not, it walks up the tree until it either finds it, or:
// It detects a .git folder, which signifies it's a git root. We assume we don't want to search beyond that
// It reaches the root of the filesystem
// We've described other cases in the NOTES.md file, which we may add here.
func findLocalFiles(dir string, searchFolders []string) ([]string, error) {
	// search for the gogo files directory in the current directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	foundGitRoot := false
	for _, file := range files {
		// if the file is a directory, and it matches one of the search folders, use it
		if file.IsDir() && slices.Contains(searchFolders, file.Name()) {
			return findGoFiles(filepath.Join(dir, file.Name()))
		}
		if file.IsDir() && file.Name() == ".git" {
			foundGitRoot = true
		}
	}
	// TODO: this is not cross-platform compatible, fix that
	// detect if we're at the root of the filesystem
	if dir == "/" {
		return nil, nil
	}
	// if we've detected a git root, we don't want to search beyond that
	if foundGitRoot {
		return nil, nil
	}
	// search up the tree
	return findLocalFiles(filepath.Dir(dir), searchFolders)
}

// findLocalFunc searches the local environment for the function,
// and returns the binary to run if it is found.
// What this really means is that:
// 1. Search for files in the local folder with a +mage tag
// 2. If none found, search for the following folders in this order: .gogo, gogofiles, magefiles
// 3. If none found, walk up the local tree and try 3
// 4. Stop at either a .git folder, or the root of the filesystem
// 5. Or stop when a .gogo_no_further or some .gogobuild config file, something
func findLocalFunc(cwd, funcToRun string, tags, codeFolders []string) (string, bool, error) {
	// this returns a list of files that match our local search
	localFiles, err := findLocalFiles(cwd, codeFolders)
	if err != nil {
		return "", false, err
	}
	if len(localFiles) == 0 {
		return "", false, nil
	}

	// determine if target exists in local namespace and use it
	exists := filesHaveFunc(localFiles, funcToRun)
	if !exists {
		return "", false, nil
	}
	return path.Join(localFiles[0]), true, nil
}

func listGlobalFuncs(opts RunOpts) ([]string, error) {
	return nil, nil
}

// TODO: implement this
func findGlobalFunc(opts RunOpts, funcToRun string) (string, bool, error) {
	return "", false, nil
}

// findGoFiles searches for all .go files in the given directory.
func findGoFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var goFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) == ".go" {
			goFiles = append(goFiles, path.Join(dir, file.Name()))
		}
	}
	return goFiles, nil
}

// getBuiltBinary returns the latest built binary for the given function.
// It first decides if a cached binary exists, and based on timestamp comparisons, decides if we should re-use it.
// The timestamp comparison is based on the go.mod, go.sum, and source .go files vs the timestamp of the binary
// If the binary is out of date or does not exist it gets built.
// We basically just enter the directory and use 'go build' to build the binary.
func getBuiltBinary(log *log.Logger, buildOpts BuildOpts) error {
	log.Printf("Checking for cached binary: %s\n", buildOpts.BinaryFilepath)
	rebuild := decideToRebuild(log, buildOpts)
	if !rebuild {
		return nil
	}
	err := Build(log, buildOpts)
	if err != nil {
		return err
	}
	return nil
}

// decideToRebuild determines if we should rebuild the binary based on the source files and the binary file
func decideToRebuild(debug *log.Logger, buildOpts BuildOpts) bool {
	sourceFiles, modified := gatherFilesToCompare(debug, buildOpts.SourceDir)
	debug.Printf("Found the following source files: %v\n", sourceFiles)
	if modified {
		return true
	}
	modified, err := fs.CompareTimes(sourceFiles, buildOpts.BinaryFilepath)
	if err != nil {
		debug.Printf("Error comparing timestamps: %v\n", err)
		return true
	}
	debug.Printf("Changes detected: %v\n", modified)
	// if the file is not modified, and we're not forcing a rebuild, return the path to the binary
	if !modified && !buildOpts.DisableCache {
		debug.Printf("Re-using binary `%s` from cache\n", buildOpts.BinaryFilepath)
		return false
	}
	if buildOpts.DisableCache {
		debug.Printf("Forcing rebuild of binary `%s`\n", buildOpts.BinaryFilepath)
		return true
	}
	if modified {
		debug.Printf("Rebuilding binary `%s`\n", buildOpts.BinaryFilepath)
		return true
	}
	return false
}

func gatherFilesToCompare(debug *log.Logger, dir string) ([]string, bool) {
	sourceFiles, err := fs.GlobMany([]string{dir}, []string{"*.go", "go.mod", "go.sum"})
	// if there's an error with the comparison, just build it
	if err != nil {
		debug.Printf("Error finding files to glob: %v\n", err)
		return nil, true
	}

	// Walk all subdirectories recursively
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		// Skip the root directory itself
		if path == dir {
			return nil
		}

		// Only process directories
		if info.IsDir() {
			// Glob files in this subdirectory
			subFiles, err := fs.GlobMany([]string{path}, []string{"*.go", "go.mod", "go.sum"})
			if err != nil {
				return fmt.Errorf("error globbing subdirectory %s: %w", path, err)
			}

			// Add files to the main list
			sourceFiles = append(sourceFiles, subFiles...)
		}

		return nil
	})

	if err != nil {
		debug.Printf("Error walking directories: %v\n", err)
		return nil, true
	}

	return sourceFiles, false
}

// TODO: This needs option overrides to determine if we should build individual binaries for each function
func convertToGoCmds(funcs []function) ([]renderData, error) {
	if len(funcs) == 0 {
		return nil, nil
	}
	rd := renderData{}
	// there's only one, so the root command can be the function
	//if len(funcs) == 1 {
	//	rd.RootCmd = convertToGoCmd(funcs[0])
	//	return []renderData{rd}, nil
	//}

	// if there are multiple functions, we need to create a root command
	// and then add the functions as subcommands
	rd.SubCommands = make([]GoCmd, len(funcs))
	for i, funk := range funcs {
		rd.SubCommands[i] = convertToGoCmd(funk)
		if rd.SubCommands[i].UseGoGoContext {
			rd.UseGoGoContext = true
		}
	}
	return []renderData{rd}, nil
}

// convertToGoCmd converts a function to a GoCmd
func convertToGoCmd(funk function) GoCmd {
	cleanup := func(s string) string {
		// remove newlines
		s = strings.ReplaceAll(s, "\n", " ")
		// escape quotes
		s = strings.ReplaceAll(s, "\"", "\\\"")
		return s
	}
	cmd := GoCmd{
		Name:           funk.Name,
		Short:          cleanup(funk.Description),
		Long:           cleanup(funk.Comment),
		Example:        funk.Example,
		GoFlags:        nil,
		ErrorReturn:    funk.ErrorReturn,
		UseGoGoContext: funk.UseGoGoCtx,
	}
	// now for each of the flags, convert them to GoFlags
	for _, argProperties := range funk.Arguments {
		// if the argumentsi the gogo context, skip adding it, since that's injected via the funk.UseGoGoCtx
		if argProperties.Type == "gogo.Context" {
			continue
		}
		flag := GoFlag{
			Type: argProperties.Type,
			Name: argProperties.Name,
		}
		flag.Default = argProperties.Default
		flag.HasDefault = argProperties.Default != nil
		if argProperties.Default == nil {
			switch argProperties.Type {
			case "bool":
				flag.Default = false
			case "float64":
				fallthrough
			case "int":
				flag.Default = 0
			case "string":
				flag.Default = `""`
			}
		}
		if argProperties.AllowedValues != nil {
			flag.AllowedValues = argProperties.AllowedValues
		}
		if argProperties.RestrictedValues != nil {
			flag.RestrictedValues = argProperties.RestrictedValues
		}
		if argProperties.Help != "" {
			flag.Help = argProperties.Help
		}
		if argProperties.Short != byte(0) {
			flag.Short = argProperties.Short
		}
		cmd.GoFlags = append(cmd.GoFlags, flag)
	}
	return cmd
}
