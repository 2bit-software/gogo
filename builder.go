// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package gogo

import (
	"fmt"
	"go/format"
	"hash/fnv"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/2bit-software/gogo/pkg/sh"
)

type renderData struct {
	UseGoGoContext bool // if any of the commands use the gogo context, then include the context in the main file
	RootCmd        GoCmd
	SubCommands    []GoCmd
}

type GoCmd struct {
	Name           string // Name of the command
	Short          string // Short Description of the command. Comes from the comment, unless overridden.
	Long           string // Long Description of the command. This comes from the comment, if it exists.
	Example        string // An example of using this command
	GoFlags        []GoFlag
	ErrorReturn    bool // If true, the command returns an error
	UseGoGoContext bool // If true, the command uses the gogo context
}

type GoFlag struct {
	Type             string // string, int, bool, float64  This type is inferred from reading the code.
	Name             string // name of the flag
	Short            byte   // short name of the flag
	Default          any    // default value of the flag.
	HasDefault       bool   // if true, use the default value
	Help             string // help text for the flag
	AllowedValues    []any  // if provided, only these values are allowed, and are auto-completed in the shell
	RestrictedValues []any  // if provided, prohibits this flag from being set to these values. Panics if detected.
}

type RunOpts struct {
	BuildOpts
	Verbose          bool   `json:"GOGO_VERBOSE"`           // output verbose information when RUNNING gogo AND the sub-command
	GlobalSourceDir  string `json:"GOGO_GLOBAL_SOURCE_DIR"` // the global location for gogo functions
	GlobalBinDir     string `json:"GOGO_GLOBAL_BIN_DIR"`    // the output location for global binaries
	BuildLocalCache  bool   `json:"GOGO_BUILD_LOCAL"`       // When true, builds the local cache and exits
	BuildGlobalCache bool   `json:"GOGO_BUILD_GLOBAL"`      // When true, builds the global cache and exits
	ScreenWidth      int    // the width of the screen, if we know
}

type BuildOpts struct {
	KeepArtifacts      bool   `json:"GOGO_KEEP_ARTIFACTS"`      // When true, keeps the artifacts after the build. This includes the go src and the built binary
	IndividualBinaries bool   `json:"GOGO_INDIVIDUAL_BINARIES"` // When true, each function results in a binary of the same name
	DisableCache       bool   `json:"GOGO_DISABLE_CACHE"`       // When true, forces a rebuild of the binary
	Optimize           bool   `json:"GOGO_OPTIMIZE"`            // should the functions be compiled with optimization flags during this run
	BinaryFilepath     string `json:"GOGO_BINARY_FILEPATH"`     // the output location of the binary. If this is provided, then individual binaries is ignored.
	// The below properties are calculated by the build process
	SourceDir          string // the location of the directory where we are currently building the source
	OutputDir          string // the output location of the binaries
	OriginalWorkingDir string // the original working directory
}

func defaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"ToUpper": strings.ToUpper,
		"Add": func(a, b int) int {
			return a + b
		},
		"Capitalize": func(s string) string {
			return strings.Title(s)
		},
		"LowerFirstLetter": func(s string) string {
			return strings.ToLower(s[:1]) + s[1:]
		},
		// It's assumed this is only used/called when a default value is provided
		"QuoteDefault": func(hasDefault bool, v any) string {
			if s, ok := v.(string); ok {
				// if there's no default, return an empty string
				if !hasDefault {
					return `""`
				}
				if s != `""` && strings.TrimSpace(s) != "" {
					return fmt.Sprintf("\"%v\"", s)
				}
			}
			// it's a number, or float, or bool, so just return the value
			return fmt.Sprintf("%v", v)
		},
		"Subtract": func(a, b int) int {
			return a - b
		},
		"StripNewlines": func(s string) string {
			return strings.ReplaceAll(s, "\n", "")
		},
		"ByteToString": func(b byte) string {
			return string(b)
		},
		"BuildStringArgList": func(flags []GoFlag) string {
			var out = "[]string{"
			for i, f := range flags {
				if i > 0 {
					out += ", "
				}
				out += fmt.Sprintf("\"%v\"", f.Name)
			}
			out += "}"
			return out
		},
	}
}

// Build reads all the gogo files in the directory, applies their
// configuration options, and builds the resulting binary.
// The buildDir is the directory in which we are building the source FROM.
// The output of the binary can be specified in the buildOpts.OutputFilepath
func Build(log *log.Logger, buildOpts BuildOpts) error {
	mainFilePath := filepath.Join(buildOpts.SourceDir, MAIN_FILENAME)
	log.Printf("Building main go file: %v\n", mainFilePath)
	err := buildSource(buildOpts.SourceDir, mainFilePath)

	// delete the main.gogo.go file when we're done
	defer func(def error) {
		if buildOpts.KeepArtifacts {
			return
		}
		log.Printf("Removing main go file from %v\n", mainFilePath)
		defErr := os.Remove(mainFilePath)
		if defErr != nil {
			log.Println(defErr)
		}
	}(err)

	if err != nil {
		return err
	}

	// build binary
	return buildBinary(buildOpts.Optimize, buildOpts.SourceDir, buildOpts.BinaryFilepath)
}

// hashString hashes a string using SHA-256
func hashString(name string) (string, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(name))
	if err != nil {
		return "", err
	}
	hash := h.Sum32()
	return fmt.Sprintf("%08x", hash)[:8], nil
}

// buildSource reads all the gogo files in the directory, applies their
// configuration options, and builds the resulting main file.
func buildSource(inputDir, filePath string) error {
	// first we need to parse all functions in the directory that match our build requirements
	funcs, err := parseDirectory(inputDir)
	if err != nil {
		return err
	}
	// then we need to convert the functions into the renderData
	rd, err := convertToGoCmds(funcs)
	if err != nil {
		return err
	}

	cmd := rd[0]

	// TODO: this is wrong. We're passing a filePath, but it's possible we need to make multiple binaries.
	//  To support this we need to render the file, build the binary, and then delete the rendered file and repeat.
	// render from templates
	rendered, err := renderFromTemplates(cmd, defaultFuncMap())
	if err != nil {
		return err
	}
	// first write the file to disk unformatted
	err = os.WriteFile(filePath, []byte(rendered), 0644)
	if err != nil {
		return err
	}
	// format it
	formatted, err := format.Source([]byte(rendered))
	if err != nil {
		return err
	}
	// write to file
	err = os.WriteFile(filePath, formatted, 0644)
	if err != nil {
		return err
	}
	return nil
}

func renderFromTemplates(rd renderData, funcMap map[string]any) (string, error) {
	tmpl := template.New("main.go.tmpl")
	tmpl = tmpl.Funcs(funcMap)
	tmpl, err := tmpl.ParseFS(templates,
		"templates/main.go.tmpl",
		"templates/subCmd.go.tmpl",
		"templates/run.go.tmpl",
	)
	if err != nil {
		return "", err
	}
	outBuf := new(strings.Builder)
	err = tmpl.Execute(outBuf, rd)
	if err != nil {
		return "", err
	}
	return outBuf.String(), nil
}

// buildBinary formats, gets dependencies, and builds the binary
func buildBinary(optimize bool, sourceDir, to string) error {
	// go mod tidy
	err := sh.Cmd("go", "mod", "tidy").Dir(sourceDir).Run()
	if err != nil {
		return fmt.Errorf("failed to tidy go modules: %w", err)
	}

	formatOutput, err := sh.Cmd("go", "fmt", ".").Dir(sourceDir).String()
	if err != nil {
		return fmt.Errorf("failed to format rendered document: `%v` due to %w", formatOutput, err)
	}

	// go get
	err = sh.Cmd("go", "get").Dir(sourceDir).Run()
	if err != nil {
		return fmt.Errorf("failed to get dependencies: %w", err)
	}

	cmd := []string{
		"go", "build",
	}
	// optimize if requested
	if optimize {
		cmd = append(cmd, "-ldflags", "-s -w")
	}
	if to == "" {
		return fmt.Errorf("no output file specified")
	}

	// add build tags
	cmd = append(cmd, "-tags=gogo,mage")

	// add the output binary
	cmd = append(cmd, "-o", to)

	// add the source directory
	cmd = append(cmd, sourceDir)

	// build
	out, err := sh.Cmd(cmd...).Dir(sourceDir).String()
	if err != nil {
		_ = os.Remove(to)
		return fmt.Errorf("failed to build binary: `%v` due to: %w", out, err)
	}
	return nil
}
