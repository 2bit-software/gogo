// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package scripts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// this test assumes only one function exists in the file
func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected function
	}{
		{
			name: "simple",
			src: `package gogo
				func NewFunc() {}`,
			expected: function{
				Name: "NewFunc",
			},
		},
		{
			name: "with comment",
			src: `package gogo
				// This is a comment
				func NewFunc() {}`,
			expected: function{
				Name:    "NewFunc",
				Comment: "This is a comment",
			},
		},
		{
			name: "with single arg",
			src: `package gogo
				func NewFunc(arg1 string) {}`,
			expected: function{
				Name: "NewFunc",
				Arguments: []argument{
					{
						Name: "arg1",
						Type: "string",
					},
				},
			},
		},
		{
			name: "with multiple args",
			src: `package gogo
				func NewFunc(arg1 string, arg2 int) {}`,
			expected: function{
				Name: "NewFunc",
				Arguments: []argument{
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
		},
		{
			name: "with comment and args",
			src: `package gogo
				// This is a comment
				func NewFunc(arg1 string, arg2 int) {}`,
			expected: function{
				Name:    "NewFunc",
				Comment: "This is a comment",
				Arguments: []argument{
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
		},
		{
			name: "non-gogo context",
			src: `package gogo
				import "context"
				func NewFunc(ctx context.Context) {
					ctx.Value("key")
				}`,
			expected: function{
				Name:       "NewFunc",
				UseGoGoCtx: false,
				Arguments: []argument{
					{
						Name: "ctx",
						Type: "context.Context",
					},
				},
			},
		},
		{
			// a gogo context is a special type of argument,
			// so it should not show up in the arguments map
			// but it should be marked as a gogo context
			// with appropriate variable name
			name: "with gogo context",
			src: fmt.Sprintf(`package gogo
				import "%s"
				func NewFunc(ctx gogo.Context) {
ctx.ShortDescription("This is a description")
}`, GOGOIMPORTPATH),
			expected: function{
				Name:                "NewFunc",
				Description:         "This is a description",
				UseGoGoCtx:          true,
				GoGoCtxVariableName: "ctx",
				Arguments:           []argument(nil),
			},
		},
		{
			name: "with gogo context all function options",
			src: fmt.Sprintf(`package gogo
				import "%s"
// This is a long description
				func NewLongFunc(ctx gogo.Context) {
ctx.ShortDescription("This is a description").
Example("This is an example")
}`, GOGOIMPORTPATH),
			expected: function{
				Name:                "NewLongFunc",
				Description:         "This is a description",
				Comment:             "This is a long description",
				Example:             "This is an example",
				UseGoGoCtx:          true,
				GoGoCtxVariableName: "ctx",
				Arguments:           []argument(nil),
			},
		},
		{
			name: "with gogo context and argument",
			src: fmt.Sprintf(`package gogo
				import "%s"
				func NewFunc(ctx gogo.Context, arg1 string) {
ctx.ShortDescription("This is a description")
}`, GOGOIMPORTPATH),
			expected: function{
				Name:                "NewFunc",
				UseGoGoCtx:          true,
				Description:         "This is a description",
				GoGoCtxVariableName: "ctx",
				Arguments: []argument{
					{
						Name: "arg1",
						Type: "string",
					},
				},
			},
		},
		{
			name: "gogo context argument information",
			src: fmt.Sprintf(`package gogo
				import "%s"	
				func NewFunc(ctx gogo.Context, var1 string) {
					ctx.ShortDescription("This is a description").
					Argument(var1).
					Help("This is an argument help message")
}`, GOGOIMPORTPATH),
			expected: function{
				Name:                "NewFunc",
				UseGoGoCtx:          true,
				Description:         "This is a description",
				GoGoCtxVariableName: "ctx",
				Arguments: []argument{
					{
						Name: "var1",
						Type: "string",
						Help: "This is an argument help message",
					},
				},
			},
		},
		{
			name: "gogo context multiple argument information",
			src: fmt.Sprintf(`package gogo
				import "%s"
				func NewFuncAdvanced(ctx gogo.Context, var1 string, var2 int) {
					ctx.ShortDescription("This is a description").
					Argument(var1).
					Help("This is an argument help message").
					Default("default value").
					Argument(var2).
					Help("This is another argument help message").
					Short('v')
				}`, GOGOIMPORTPATH),
			expected: function{
				Name:                "NewFuncAdvanced",
				UseGoGoCtx:          true,
				Description:         "This is a description",
				GoGoCtxVariableName: "ctx",
				Arguments: []argument{
					{
						Name:    "var1",
						Type:    "string",
						Help:    "This is an argument help message",
						Default: "default value",
					},
					{
						Name:  "var2",
						Type:  "int",
						Help:  "This is another argument help message",
						Short: 'v',
					},
				},
			},
		},
		{
			name: "gogo context with alias",
			src: fmt.Sprintf(`package gogo
				import g2 "%s"
				func NewFuncAlias(ctx g2.Context) {
					ctx.ShortDescription("This is a description")
				}`, GOGOIMPORTPATH),
			expected: function{
				Name:                "NewFuncAlias",
				UseGoGoCtx:          true,
				Description:         "This is a description",
				GoGoCtxVariableName: "ctx",
				Arguments:           []argument(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			funcs, err := parseSource(tt.src)
			require.NoError(t, err)
			require.Equal(t, tt.expected, funcs[0])
		})
	}
}

// this test assumes many functions exist in the file
func TestParseMany(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected []function
	}{
		{
			name: "multiple functions no context",
			src: `package gogo
				func NewFunc1(func1Arg string) {}
				func NewFunc2(func2Arg bool) {}`,
			expected: []function{
				{
					Name: "NewFunc1",
					Arguments: []argument{
						{
							Name: "func1Arg",
							Type: "string",
						},
					},
				},
				{
					Name: "NewFunc2",
					Arguments: []argument{
						{
							Name: "func2Arg",
							Type: "bool",
						},
					},
				},
			},
		},
		{
			name: "multiple functions single context",
			src: fmt.Sprintf(`package gogo
				import "%s"
				func NewFunc1(ctx gogo.Context, func1Arg string) {
					ctx.ShortDescription("This is a description")
				}
				func NewFunc2(func2Arg bool) {}`, GOGOIMPORTPATH),
			expected: []function{
				{
					Name:                "NewFunc1",
					UseGoGoCtx:          true,
					Description:         "This is a description",
					GoGoCtxVariableName: "ctx",
					Arguments: []argument{
						{
							Name: "func1Arg",
							Type: "string",
						},
					},
				},
				{
					Name: "NewFunc2",
					Arguments: []argument{
						{
							Name: "func2Arg",
							Type: "bool",
						},
					},
				},
			},
		},
		{
			name: "multiple functions multiple context",
			src: fmt.Sprintf(`package gogo
				import "%s"
				func NewFunc1(ctx gogo.Context, func1Arg string) {
					ctx.ShortDescription("This is a description")
				}
				func NewFunc2(ctx gogo.Context, func2Arg bool) {
					ctx.ShortDescription("This is another description")
				}`, GOGOIMPORTPATH),
			expected: []function{
				{
					Name:                "NewFunc1",
					UseGoGoCtx:          true,
					Description:         "This is a description",
					GoGoCtxVariableName: "ctx",
					Arguments: []argument{
						{
							Name: "func1Arg",
							Type: "string",
						},
					},
				},
				{
					Name:                "NewFunc2",
					UseGoGoCtx:          true,
					Description:         "This is another description",
					GoGoCtxVariableName: "ctx",
					Arguments: []argument{
						{
							Name: "func2Arg",
							Type: "bool",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			funcs, err := parseSource(tt.src)
			require.NoError(t, err)
			require.Equal(t, tt.expected, funcs)
		})
	}
}
