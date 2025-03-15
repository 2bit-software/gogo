package gogo

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/2bit-software/gogo/pkg/mod"
	"github.com/2bit-software/gogo/pkg/sh"
)

func setupBinaries(t *testing.T, testFolder string) {
	// build the gogo binary
	// delete the /tmp/gogo binary if it exists
	_ = os.Remove("/tmp/gogo")
	cmd := sh.Cmd("go run cmd/gogo/main.go gadget CompileGo --binaryName=gogo --inputFolderPath=./cmd/gogo --outputFolderPath=/tmp --tags=gogo --versionPath=\"github.com/2bit-software/gogo/cmd/gogo/cmds\"\n")
	buildResult, err := cmd.String()
	require.NoErrorf(t, err, "failed to build gogo binary: %s", buildResult)

	// make sure the binary exists and runs
	_, err = os.Stat("/tmp/gogo")
	require.NoError(t, err)
	res, err := sh.Cmd("/tmp/gogo --version").String()
	require.NoError(t, err)
	require.Truef(t, strings.HasPrefix(res, "v"), "expected version string to start with 'v', got %s", res)

	root, err := mod.FindModuleRoot()
	require.NoError(t, err)

	// reset the working directory after the test
	cwd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(cwd)
	}()

	scenarioDir := path.Join(root, "scenarios")
	t.Logf("scenarioDir: %s", path.Join(scenarioDir, testFolder))

	// change to the scenario dir
	err = os.Chdir(path.Join(scenarioDir, testFolder))
	require.NoError(t, err)

	// see if we can build the scenario
	cmd = sh.Cmd("/tmp/gogo build --verbose -o /tmp/gadgets")
	gadgetBuildResult, err := cmd.String()
	require.NoErrorf(t, err, "error building scenario: %s", gadgetBuildResult)
	require.NotEmpty(t, gadgetBuildResult)
}

// TODO: try parsing args using a mix of --flags, positional args, and using defaults
// TODO: write tests that use arguments we expect to intentionally fail

func TestStandardArgParsing(t *testing.T) {
	testFolder := "standard"

	tests := []struct {
		command  string
		args     []string
		expected string
	}{
		{
			command:  "NoArgumentsNoReturns",
			args:     []string{},
			expected: "NoArgumentsNoReturns",
		},
		{
			command:  "ErrorReturn",
			args:     []string{},
			expected: "ErrorReturn",
		},
		{
			command:  "SingleArgument",
			args:     []string{"passedArg1"},
			expected: "SingleArgument with arg1: passedArg1",
		},
		{
			command:  "SingleArgumentAndErrorReturn",
			args:     []string{"passedArg9"},
			expected: "SingleArgumentAndErrorReturn with arg1: passedArg9",
		},
		{
			command:  "TwoDifferentArguments",
			args:     []string{"passedArg1", "true"},
			expected: "TwoDifferentArguments with arg1: passedArg1, arg2: true",
		},
		{
			command:  "TwoDifferentArgumentsAndErrorReturn",
			args:     []string{"passedArg1", "true"},
			expected: "TwoDifferentArgumentsAndErrorReturn with arg1: passedArg1, arg2: true",
		},
		{
			command:  "AdvancedFunction",
			args:     []string{"passedName", "true", "9"},
			expected: "name: passedName value: 9",
		},
		{
			command:  "AdvancedFunction",
			args:     []string{"passedName", "false", "9"},
			expected: "name: passedName",
		},
	}

	setupBinaries(t, testFolder)

	for _, test := range tests {
		t.Run(test.command, func(t *testing.T) {
			// run the scenario
			scenarioCmd := sh.Cmd("/tmp/gadgets " + test.command).SetEnv([]string{})
			scenarioCmd.SetEnv([]string{})
			if len(test.args) > 0 {
				scenarioCmd.SetArgs(test.args...)
			}
			result, err := scenarioCmd.String()
			require.NoError(t, err)
			require.Equal(t, test.expected, strings.TrimSpace(result))
		})
	}
}

// We don't need to test a lot of these, since we're just testing if using a go context
// still works. However, this is just testing arg parsing, not context config parsing,
// so there's not a lot to do here.
func TestAliasedCtxArgParsing(t *testing.T) {
	testFolder := "aliased"

	tests := []struct {
		command  string
		args     []string
		expected string
	}{
		{
			command:  "AliasedCtxDescription",
			args:     []string{},
			expected: "AliasedCtxDescription",
		},
		{
			command:  "AliasedCtxArgument",
			args:     []string{"passedArg1", "true"},
			expected: "AliasedCtxArgument with var1: passedArg1 and var2: true",
		},
	}

	setupBinaries(t, testFolder)

	for _, test := range tests {
		t.Run(test.command, func(t *testing.T) {
			// run the scenario
			scenarioCmd := sh.Cmd("/tmp/gadgets " + test.command).SetEnv([]string{})
			scenarioCmd.SetEnv([]string{})
			if len(test.args) > 0 {
				scenarioCmd.SetArgs(test.args...)
			}
			result, err := scenarioCmd.String()
			require.NoErrorf(t, err, "failed to run scenario %s: %s", test.command, result)
			require.Equal(t, test.expected, strings.TrimSpace(result))
		})
	}
}

// Same with this test as the aliased test. We're just testing arg parsing, not context
func TestUniqueGoModArgParsing(t *testing.T) {
	testFolder := "unique_gomod"

	tests := []struct {
		command  string
		args     []string
		expected string
	}{
		{
			command:  "BasicDescription",
			args:     []string{},
			expected: "BasicDescription",
		},
		{
			command:  "BasicArgument",
			args:     []string{"passedArg1", "true"},
			expected: "BasicArgument with var1: passedArg1 and var2: true",
		},
	}

	setupBinaries(t, testFolder)

	for _, test := range tests {
		t.Run(test.command, func(t *testing.T) {
			// run the scenario
			scenarioCmd := sh.Cmd("/tmp/gadgets " + test.command).SetEnv([]string{})
			scenarioCmd.SetEnv([]string{})
			if len(test.args) > 0 {
				scenarioCmd.SetArgs(test.args...)
			}
			result, err := scenarioCmd.String()
			require.NoErrorf(t, err, "failed to run scenario %s: %s", test.command, result)
			require.Equal(t, test.expected, strings.TrimSpace(result))
		})
	}
}
