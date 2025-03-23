package gadgets

import (
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

// TestInitLocal tests the Init function, with no parameters
// So it should make a local .gogo folder and a default go.mod path
func TestInitLocal(t *testing.T) {
	// revert to the original directory after the test
	cwd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		err = os.Chdir(cwd)
		require.NoError(t, err)
	}()

	// get a temporary directory to work in
	tmpDir := t.TempDir()
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	folder := path.Join(".", ".gogo")
	err = Init(folder, GOGOIMPORTPATH)
	require.NoError(t, err)

	// assert the .gogo folder exists
	_, err = os.Stat(folder)
	require.NoError(t, err)

	// now assert the go.mod file exists
	goModPath := path.Join(folder, "go.mod")
	_, err = os.Stat(goModPath)
	require.NoError(t, err)
}

func TestRenderExamples(t *testing.T) {
	// revert to the original directory after the test
	cwd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		err = os.Chdir(cwd)
		require.NoError(t, err)
	}()

	// get a temporary directory to work in
	tmpDir := t.TempDir()
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// render examples
	err = renderExamples(tmpDir, "github.com/2bit-software/example/pkg", GOGOIMPORTPATH)
	require.NoError(t, err)

	// check if all the files exist
	files := []string{
		"hello.go",
		"pkg/hello.go",
		"pkg/hello_test.go",
	}
	for _, file := range files {
		filePath := path.Join(tmpDir, file)
		_, err = os.Stat(filePath)
		require.NoError(t, err)
	}
}
