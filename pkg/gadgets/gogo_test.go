package gadgets

import (
	"github.com/2bit-software/gogo/pkg/mod"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"path"
	"strings"
	"testing"
)

// Test the description of the function. The description is used
// in listing the summary of the function.
func TestDescription(t *testing.T) {
	// revert to the original directory after the test
	cwd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		err = os.Chdir(cwd)
		require.NoError(t, err)
	}()

	// change to the directory where the test functions are located, buildFuncList requires this
	root, err := mod.FindModuleRoot()
	require.NoError(t, err)
	scenarioPath := path.Join(root, "scenarios", "standard")
	err = os.Chdir(scenarioPath)
	require.NoError(t, err)

	opts := RunOpts{
		Verbose: false,
	}
	wd, err := os.Getwd()
	require.NoError(t, err)
	funcList, err := BuildFuncList(opts, wd)
	require.NoError(t, err)
	require.NoError(t, err)
	output := generateFuncListOutput(funcList, 300)
	// assert that the output contains the expected description
	description := `This is the description for the function. Without any other arguments to the ctx, this will show up in the list view and the --help output.`
	found := false
	for _, out := range output {
		if strings.Contains(out, description) {
			found = true
		}
	}
	assert.True(t, found)
}

// Test the ShortDescription of the function.
func TestShortDescription(t *testing.T) {
	// revert to the original directory after the test
	cwd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		err = os.Chdir(cwd)
		require.NoError(t, err)
	}()

	// change to the directory where the test functions are located, buildFuncList requires this
	root, err := mod.FindModuleRoot()
	require.NoError(t, err)
	scenarioPath := path.Join(root, "scenarios", "standard")
	err = os.Chdir(scenarioPath)
	require.NoError(t, err)

	opts := RunOpts{
		Verbose: false,
	}
	wd, err := os.Getwd()
	require.NoError(t, err)
	funcList, err := BuildFuncList(opts, wd)
	require.NoError(t, err)
	require.NoError(t, err)
	output := generateFuncListOutput(funcList, 300)
	// assert that the output contains the expected short description
	description := `this is a short description set specifically for the BasicShortDescription function`
	found := false
	for _, out := range output {
		if strings.Contains(out, description) {
			found = true
		}
	}
	assert.True(t, found)
}

func TestBuild(t *testing.T) {
	l := log.New(os.Stdout, "", log.LstdFlags)
	root, err := mod.FindModuleRoot()
	require.NoError(t, err)
	// make a temp dir
	tmpDir, err := os.MkdirTemp("", "gogo-test")
	require.NoError(t, err)

	opts := BuildOpts{
		KeepArtifacts:  false,
		DisableCache:   true,
		Optimize:       false,
		SourceDir:      path.Join(root, "scenarios", "standard", ".gogo"),
		OutputDir:      tmpDir,
		BinaryFilepath: path.Join(tmpDir, "gadgets"),
	}

	// build the function
	err = Build(l, opts)
	require.NoError(t, err)
}
