// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package gadgets

import (
	"github.com/2bit-software/gogo/pkg/mod"
	"github.com/stretchr/testify/require"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test parsing importing gogo in various ways
func TestParseImports(t *testing.T) {
	tests := []struct {
		name           string
		filePath       string
		expectedImport string
		err            bool
	}{
		{
			name:           "basic",
			filePath:       "scenarios/standard/.gogo/advanced.go",
			expectedImport: "gogo",
		},
		{
			name:           "aliased",
			filePath:       "scenarios/aliased/.gogo/aliased.go",
			expectedImport: "goCtx",
		},
	}
	root, err := mod.FindModuleRoot()
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the source code
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, path.Join(root, tt.filePath), nil, 0)
			if err != nil {
				t.Fatalf("error parsing file: %v", err)
			}

			importName, _ := getGoGoImportName(file)
			if importName != tt.expectedImport {
				t.Fatalf("expectedComment import name to be %q, got %q", tt.expectedImport, importName)
			}
		})
	}
}

// test when they are not chained together.
// TODO: what is this test good for? we need to update it.
// I think what we really want is to just test using every function signature
// and every ctx option.
func TestNoChainedOptions(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		expected any
		err      bool
	}{
		{
			name:     "ctx description",
			funcName: "DescriptionOnly",
		},
		{
			name:     "single argument",
			funcName: "SingleArgument",
		},
	}
	root, err := mod.FindModuleRoot()
	require.NoError(t, err)
	// Parse the source code
	funcs, err := parse(path.Join(root, "scenarios/standard/.gogo/basic.go"))
	require.NoError(t, err)
	require.NotNil(t, funcs)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// for each funcNames, retrieve that function
			isInside := slices.ContainsFunc(funcs, func(f function) bool {
				return strings.EqualFold(f.Name, tt.funcName)
			})
			assert.True(t, isInside)
		})
	}
}

func TestParsePlainFunc(t *testing.T) {
	tests := []struct {
		name            string
		src             string
		expectedName    string
		expectedComment string
		expectedArgs    []argument
		expectedError   bool
	}{
		{
			name:         "function name",
			src:          "package main\nfunc TestFunc() {}",
			expectedName: "TestFunc",
		},
		{
			name:            "function comment",
			src:             "package main\n // TestFunc is a test function\nfunc TestFunc() {}",
			expectedName:    "TestFunc",
			expectedComment: "TestFunc is a test function",
		},
		{
			name:         "function arguments",
			src:          "package main\nfunc TestFuncArgs(arg1 string, arg2 int) {}",
			expectedName: "TestFuncArgs",
			expectedArgs: []argument{
				{
					Name: "arg1",
					Type: "string",
				},
				{
					Name: "arg2",
					Type: "int",
				},
			},
		},
		{
			name:          "function with error return",
			src:           "package main\nfunc TestErrorFunc() error {}",
			expectedName:  "TestErrorFunc",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tt.src, parser.ParseComments)
			assert.NoError(t, err)

			funcDecl := file.Decls[0].(*ast.FuncDecl)
			pCtx, err := parsePlainFunc(funcDecl)
			assert.NoError(t, err)
			if tt.expectedName != "" {
				assert.Equal(t, tt.expectedName, pCtx.Name)
			}
			if tt.expectedComment != "" {
				assert.Equal(t, tt.expectedComment, pCtx.Comment)
			}
			// parse args
			pCtx = gatherDetails(pCtx, "gogo", funcDecl)
			if tt.expectedArgs != nil {
				assert.Equal(t, tt.expectedArgs, pCtx.Arguments)
			}
			assert.Equal(t, tt.expectedError, pCtx.ErrorReturn)
		})
	}
}
