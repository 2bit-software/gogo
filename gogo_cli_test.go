package gogo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/2bit-software/gogo/pkg/mod"
	"github.com/2bit-software/gogo/pkg/sh"
)

const (
	GOGO_FILENAME = "gogo"
)

// returns the path to the gadgets binary
func setupBinaries(t *testing.T, testFolder string) string {
	// make a new temporary directory for the test
	dir, err := os.MkdirTemp("", "gogo_test")
	require.NoError(t, err)

	gogoFilePath := path.Join(dir, GOGO_FILENAME)
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// list the funcs
	funcs, err := sh.Cmd("go run cmd/gogo/main.go").String()
	require.NoError(t, err)
	fmt.Println(funcs)

	// build the gogo binary
	// delete the /tmp/gogo binary if it exists
	_ = os.Remove(gogoFilePath)
	cmd := sh.Cmd(fmt.Sprintf("go run cmd/gogo/main.go gadget CompileGo --binaryName=gogo --inputFolderPath=./cmd/gogo --outputFolderPath=%s --tags=gogo --versionPath=\"github.com/2bit-software/gogo/cmd/gogo/cmds\"\n", dir))
	buildResult, err := cmd.String()
	require.NoErrorf(t, err, "failed to build gogo binary: %s", buildResult)

	// make sure the binary exists and runs
	_, err = os.Stat(gogoFilePath)
	require.NoError(t, err)
	res, err := sh.Cmd(fmt.Sprintf("%s --version", gogoFilePath)).String()
	require.NoError(t, err)
	require.Truef(t, strings.HasPrefix(res, "v"), "expected version string to start with 'v', got %s", res)

	root, err := mod.FindModuleRoot()
	require.NoError(t, err)

	// reset the working directory after the test
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(cwd)
	}()

	scenarioDir := path.Join(root, "scenarios", testFolder)
	t.Logf("scenarioDir: %s", scenarioDir)

	gadgetsPath := path.Join(dir, "gadgets")

	// see if we can build the scenario
	cmd = sh.Cmd(fmt.Sprintf("%s build --verbose -o %s", gogoFilePath, gadgetsPath)).
		AddEnv([]string{"GOGO_DISABLE_CACHE=true"}).Dir(scenarioDir)
	gadgetBuildResult, err := cmd.String()
	require.NoErrorf(t, err, "error building scenario: %s", gadgetBuildResult)
	require.NotEmpty(t, gadgetBuildResult)
	return gadgetsPath
}

func TestSetupBinaries(t *testing.T) {
	testFolder := "standard/.gogo"
	gadgetsBinaryPath := setupBinaries(t, testFolder)
	assert.FileExists(t, gadgetsBinaryPath)
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
			command:  "ContextWithNoUsage",
			args:     []string{},
			expected: "ContextWithNoUsage",
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
			// this test shows that if you fill a positional arg with an empty string (two quotes),
			// then we strip it out.
			command:  "SingleArgument",
			args:     []string{`""`},
			expected: "SingleArgument with arg1:",
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

	gadgetsBinaryPath := setupBinaries(t, testFolder)

	for _, test := range tests {
		t.Run(test.command, func(t *testing.T) {
			// run the scenario
			scenarioCmd := sh.Cmd(fmt.Sprintf("%s %s", gadgetsBinaryPath, test.command)).SetEnv([]string{})
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

	gadgetsBinaryPath := setupBinaries(t, testFolder)

	for _, test := range tests {
		t.Run(test.command, func(t *testing.T) {
			// run the scenario
			scenarioCmd := sh.Cmd(fmt.Sprintf("%s %s", gadgetsBinaryPath, test.command)).SetEnv([]string{})
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

	gadgetsBinaryPath := setupBinaries(t, testFolder)

	for _, test := range tests {
		t.Run(test.command, func(t *testing.T) {
			// run the scenario
			scenarioCmd := sh.Cmd(fmt.Sprintf("%s %s", gadgetsBinaryPath, test.command)).SetEnv([]string{})
			if len(test.args) > 0 {
				scenarioCmd.SetArgs(test.args...)
			}
			result, err := scenarioCmd.String()
			require.NoErrorf(t, err, "failed to run scenario %s: %s", test.command, result)
			require.Equal(t, test.expected, strings.TrimSpace(result))
		})
	}
}

// Here we're testing pure flag parsing, without any positional arguments
func TestFlagArgParsing(t *testing.T) {
	testFolder := "standard"

	tests := []struct {
		command  string
		args     []string
		expected string
	}{
		{
			command:  "SingleArgument",
			args:     []string{"--arg1=passedArg1"},
			expected: "SingleArgument with arg1: passedArg1",
		},
		{
			// tests that both arguments passed as flags is enough
			command:  "TwoDifferentArguments",
			args:     []string{"--arg1=passedArg1", "--arg2"},
			expected: "TwoDifferentArguments with arg1: passedArg1, arg2: true",
		},
		{
			// tests that passing only the first argument as a flag is enough when
			// two flags exist
			command:  "TwoDifferentArguments",
			args:     []string{"--arg1=passedArg1"},
			expected: "TwoDifferentArguments with arg1: passedArg1, arg2: false",
		},
		{
			// tests that passing only the second argument as a flag is enough when
			// two flags exist
			command:  "TwoDifferentArguments",
			args:     []string{"--arg2"},
			expected: "TwoDifferentArguments with arg1: , arg2: true",
		},
		{
			// flags in different order
			command:  "TwoDifferentArguments",
			args:     []string{"--arg2", "--arg1=passedArg1"},
			expected: "TwoDifferentArguments with arg1: passedArg1, arg2: true",
		},
	}

	gadgetsBinaryPath := setupBinaries(t, testFolder)

	for _, test := range tests {
		t.Run(test.command, func(t *testing.T) {
			// run the scenario
			scenarioCmd := sh.Cmd(fmt.Sprintf("%s %s", gadgetsBinaryPath, test.command)).SetEnv([]string{})
			if len(test.args) > 0 {
				scenarioCmd.SetArgs(test.args...)
			}
			result, err := scenarioCmd.String()
			require.NoErrorf(t, err, "failed to run scenario %s: %s", test.command, result)
			require.Equal(t, test.expected, strings.TrimSpace(result))
		})
	}
}

func TestFlagAndPositionalArgs(t *testing.T) {
	testFolder := "standard"

	tests := []struct {
		name     string
		command  string
		args     []string
		expected string
	}{
		{
			name:     "three arg func with all positional args",
			command:  "ThreeArgFuncWithContext",
			args:     []string{"passedArg1", "true", "3"},
			expected: "ThreeArgFuncWithContext with name: passedArg1, include: true, value: 3",
		},
		{
			name:     "three arg func with all flag args",
			command:  "ThreeArgFuncWithContext",
			args:     []string{"--name=passedArg1", "--include", "--value=3"},
			expected: "ThreeArgFuncWithContext with name: passedArg1, include: true, value: 3",
		},
		{
			name:     "three arg func with mixed positional and flag args",
			command:  "ThreeArgFuncWithContext",
			args:     []string{"passedArg1", "--include", "--value=3"},
			expected: "ThreeArgFuncWithContext with name: passedArg1, include: true, value: 3",
		},
		{
			name:     "three arg func with mixed positional and flag args",
			command:  "ThreeArgFuncWithContext",
			args:     []string{"passedArg1", "true", "--value=3"},
			expected: "ThreeArgFuncWithContext with name: passedArg1, include: true, value: 3",
		},
		{
			name:     "three arg func with mixed positional and flag args",
			command:  "ThreeArgFuncWithContext",
			args:     []string{"--value=3", "passedArg1", "true"},
			expected: "ThreeArgFuncWithContext with name: passedArg1, include: true, value: 3",
		},
	}

	gadgetsBinaryPath := setupBinaries(t, testFolder)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// run the scenario
			scenarioCmd := sh.Cmd(fmt.Sprintf("%s %s", gadgetsBinaryPath, test.command)).SetEnv([]string{})
			if len(test.args) > 0 {
				scenarioCmd.SetArgs(test.args...)
			}
			result, err := scenarioCmd.String()
			require.NoErrorf(t, err, "failed to run scenario %s: %s", test.command, result)
			require.Equal(t, test.expected, strings.TrimSpace(result))
		})
	}
}
