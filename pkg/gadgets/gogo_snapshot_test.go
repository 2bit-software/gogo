// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package gadgets

import (
	"github.com/2bit-software/gogo/pkg/mod"
	"go/format"
	"os"
	"path"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/require"
)

var (
	SCENARIOS = []string{
		"standard",
		"aliased",
		"unique_gomod",
	}
)

// Tests that the GoGo functions are correctly parsed from the source code,
// and each scenario renders as expected. This is a very coupled test, but it
// is useful to detect changes in the output of GoGo.
func TestRenderTemplates(t *testing.T) {
	tests := []struct {
		name       string
		renderData renderData
	}{
		{
			name: "empty",
			renderData: renderData{
				RootCmd: GoCmd{
					Name:    "rootFlag",
					GoFlags: nil,
				},
			},
		},
		{
			name: "rootcmd with flags",
			renderData: renderData{
				RootCmd: GoCmd{
					Name: "rootFlag",
					GoFlags: []GoFlag{
						{
							Type:    "string",
							Name:    "stringFlag",
							Short:   's',
							Default: "default",
							Help:    "help text",
						},
					},
				},
			},
		},
		{
			name: "subCmd with flags",
			renderData: renderData{
				RootCmd: GoCmd{
					Name:    "rootFlag",
					GoFlags: nil,
				},
				SubCommands: []GoCmd{
					{
						Name: "subCmd",
						GoFlags: []GoFlag{
							{
								Type:    "string",
								Name:    "stringFlag",
								Short:   's',
								Default: "default",
								Help:    "help text",
							},
						},
					},
				},
			},
		},
		{
			name: "cmd with error return",
			renderData: renderData{
				RootCmd: GoCmd{
					Name: "rootFlag",
					GoFlags: []GoFlag{
						{
							Type:    "string",
							Name:    "stringFlag",
							Short:   's',
							Default: "default",
							Help:    "help text",
						},
					},
					ErrorReturn: true,
				},
			},
		},
	}

	templateNames := []string{
		"templates/main.go.tmpl",
		"templates/subCmd.go.tmpl",
		"templates/function.go.tmpl",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.renderData
			// set the import path
			data.GoGoImportPath = GOGOIMPORTPATH
			res, err := renderFromTemplates(tt.renderData, defaultFuncMap(), templateNames)
			if err != nil {
				t.Fatal(err)
			}
			// format it
			formatted, err := format.Source([]byte(res))
			if err != nil {
				t.Fatalf("error formatting source: %v", err)
			}
			cupaloy.SnapshotT(t, formatted)
		})
	}
}

// Test whether the BuildFuncList function correctly parses the GoGo functions
// and then lists them out as expected. This test is very coupled to the output rendering
// of GoGo, so it is not great other than detecting changes. It is not a good signal
// that GoGo is broken.
func TestBuildFuncList(t *testing.T) {
	// save current directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	root, err := mod.FindModuleRoot()
	require.NoError(t, err)

	snapshotDir := path.Join(root, "pkg", "gadgets", ".snapshots")
	t.Logf("snapshot dir: %s", snapshotDir)
	scenarioDir := path.Join(root, "scenarios")
	t.Logf("scenarioDir: %s", scenarioDir)

	for _, scenario := range SCENARIOS {
		t.Run(scenario, func(t *testing.T) {
			err := os.Chdir(path.Join(scenarioDir, scenario))
			require.NoError(t, err)
			opts := RunOpts{
				Verbose: true,
			}
			funcList, err := BuildFuncList(opts, "")
			require.NoError(t, err)

			// return, so that we can snapshot the output in the correct directory
			err = os.Chdir(originalDir)
			require.NoError(t, err)
			cfg := cupaloy.NewDefaultConfig()
			cfg = cfg.WithOptions(cupaloy.SnapshotSubdirectory(snapshotDir))
			cfg.SnapshotT(t, funcList)
		})
	}
}

func TestPrintFuncList(t *testing.T) {
	// save current directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	root, err := mod.FindModuleRoot()
	require.NoError(t, err)

	snapshotDir := path.Join(root, "pkg", "gadgets", ".snapshots")
	t.Logf("snapshot dir: %s", snapshotDir)
	scenarioDir := path.Join(root, "scenarios")
	t.Logf("scenarioDir: %s", scenarioDir)

	for _, scenario := range SCENARIOS {
		t.Run(scenario, func(t *testing.T) {
			err := os.Chdir(path.Join(scenarioDir, scenario))
			require.NoError(t, err)
			opts := RunOpts{
				Verbose: true,
			}
			funcList, err := BuildFuncList(opts, "")
			require.NoError(t, err)
			output := generateFuncListOutput(funcList, 300)

			// return, so that we can snapshot the output in the correct directory
			err = os.Chdir(originalDir)
			require.NoError(t, err)

			cfg := cupaloy.NewDefaultConfig()
			cfg = cfg.WithOptions(cupaloy.SnapshotSubdirectory(snapshotDir))
			cfg.SnapshotT(t, output)
		})
	}
}
