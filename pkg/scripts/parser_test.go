// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package scripts

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
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
			name:     "no imports",
			filePath: "scenarios/broken/.gogo/broken.go",
			err:      true,
		},
		{
			name:           "basic",
			filePath:       "scenarios/advanced/.gogo/advanced.go",
			expectedImport: "gogo",
		},
		{
			name:           "aliased",
			filePath:       "scenarios/aliased/.gogo/aliased.go",
			expectedImport: "gogo2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the source code
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, tt.filePath, nil, 0)
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

// This loads a file, and determines which functions are exported.
func TestFindGoGoFunctions(t *testing.T) {
	tests := []struct {
		name      string
		gogoAlias string
		filePath  string
		expected  any
		err       bool
	}{
		{
			name:      "basic",
			gogoAlias: "gogo",
			filePath:  "scenarios/basic/.gogo/basic.go",
		},
		{
			name:      "aliased",
			gogoAlias: "gogo2",
			filePath:  "scenarios/aliased/.gogo/aliased.go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the source code
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, tt.filePath, nil, 0)
			if err != nil {
				t.Fatalf("error parsing file: %v", err)
			}
			funcs := getGoGoFunctions(tt.gogoAlias, file)

			// Convert the result to a slice of function names
			funcDeclNames := make([]string, len(funcs))
			for i, funcDecl := range funcs {
				funcDeclNames[i] = funcDecl.Name.Name
			}

			cupaloy.SnapshotT(t, funcDeclNames)
		})
	}
}

// test when they are not chained together.
func TestNoChainedOptions(t *testing.T) {
	tests := []struct {
		name       string
		funcName   string
		optionName string
		expected   any
		err        bool
	}{
		{
			name:       "ctx description",
			funcName:   "BasicDescription",
			optionName: "SetDescription",
			expected:   "set a description",
		},
		{
			name:       "ctx set help",
			funcName:   "BasicSetHelp",
			optionName: "SetHelp",
			expected:   "is this a thing?",
		},
	}
	// Parse the source code
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "scenarios/basic/.gogo/basic.go", nil, 0)
	if err != nil {
		t.Fatalf("error parsing file: %v", err)
	}
	alias, _ := getGoGoImportName(astFile)
	funcs := getGoGoFunctions(alias, astFile)
	if funcs == nil {
		t.Fatal("expectedComment functions to be found")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// for each funcNames, retrieve that function
			funk := getFuncByName(funcs, tt.funcName)
			if funk == nil {
				t.Fatalf("expectedComment function %q to be found", tt.funcName)
			}
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
