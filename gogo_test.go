// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package gogo

import (
	"go/format"
	"os"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/require"
)

var (
	SCENARIOS = []string{
		"basic",
		"aliased",
		"advanced",
		"many_files",
		"gogo_unique_gomod",
		"gogo_most_advanced",
	}
)

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := renderFromTemplates(tt.renderData, defaultFuncMap())
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

func TestBuildFuncList(t *testing.T) {
	// save current directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)

	for _, scenario := range SCENARIOS {
		t.Run(scenario, func(t *testing.T) {
			err := os.Chdir("scenarios/" + scenario)
			require.NoError(t, err)
			opts := RunOpts{
				Verbose: true,
			}
			funcList, err := BuildFuncList(opts)
			require.NoError(t, err)

			// return, so that we can snapshot the output in the correct directory
			err = os.Chdir(originalDir)
			require.NoError(t, err)
			cupaloy.SnapshotT(t, funcList)
		})
	}
}

func TestPrintFuncList(t *testing.T) {
	// save current directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)

	for _, scenario := range SCENARIOS {
		t.Run(scenario, func(t *testing.T) {
			err := os.Chdir("scenarios/" + scenario)
			require.NoError(t, err)
			opts := RunOpts{
				Verbose: true,
			}
			funcList, err := BuildFuncList(opts)
			require.NoError(t, err)
			output := generateFuncListOutput(funcList, 300)

			// return, so that we can snapshot the output in the correct directory
			err = os.Chdir(originalDir)
			require.NoError(t, err)
			cupaloy.SnapshotT(t, output)
		})
	}
}
