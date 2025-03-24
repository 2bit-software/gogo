package gadgets

import (
	"errors"
	"fmt"
	"github.com/2bit-software/gogo/pkg/sh"
	"github.com/2bit-software/gogo/pkg/version"
	"os"
	"path"
	"text/template"
)

const (
	REQUIRED_VERSION = "v1.24.0"
)

var (
	HELLO_TEMPLATE          = template.New("hello.go.tmpl")
	HELLO_PKG_TEMPLATE      = template.New("hello.pkg.go.tmpl")
	HELLO_PKG_TEST_TEMPLATE = template.New("hello_test.pkg.go.tmpl")
)

func init() {
	var err error
	// load the main template file
	HELLO_TEMPLATE, err = HELLO_TEMPLATE.ParseFS(templates, "templates/hello.go.tmpl")
	if err != nil {
		panic(err)
	}
	// load the pkg template file
	HELLO_PKG_TEMPLATE, err = HELLO_PKG_TEMPLATE.ParseFS(templates, "templates/hello.pkg.go.tmpl")
	if err != nil {
		panic(err)
	}
	// load the pkg test template file
	HELLO_PKG_TEST_TEMPLATE, err = HELLO_PKG_TEST_TEMPLATE.ParseFS(templates, "templates/hello_test.pkg.go.tmpl")
	if err != nil {
		panic(err)
	}
}

// Init a new GoGo project. It's assumed the folder is the full path to what should be made
// and where the files should be places. The gogoModPath is the go.mod module path to use.
func Init(folder, goModPath string) error {
	fmt.Printf("Initializing GoGo workspace in %s with goModPath %s\n", folder, goModPath)
	// make sure the folder exists, and if not, create it
	if err := os.MkdirAll(folder, 0755); err != nil && !os.IsExist(err) {
		return err
	}

	err := ensureGoMod(folder, goModPath)
	if err != nil {
		return err
	}

	// handle example "hello" function
	err = renderExamples(folder, goModPath+"/pkg", GOGOIMPORTPATH)
	if err != nil {
		return err
	}

	// ensure the dependencies are installed
	err = ensureDeps(folder)
	return err
}

// ensureGoMod makes sure a go.mod file exists in the folder, and if not, creates one.
// It uses the goModPath as the mod package namespace, if it doesn't exist.
func ensureGoMod(folder, goModPath string) error {
	// make sure the folder exists, and if not, error out
	if _, err := os.Stat(folder); err != nil {
		return err
	}

	// detect if a go.mod file exists
	goModFilePath := path.Join(folder, "go.mod")
	_, err := os.Stat(goModFilePath)
	// return errors that are not "file does not exist"
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	exists := err == nil
	if exists {
		return nil
	}
	// create a go.mod file
	fmt.Printf("Creating go.mod file with namespace %s\n", goModPath)
	return sh.Cmd("go mod init " + goModPath).Dir(folder).RunAndStream()
}

func ensureDeps(folder string) error {
	// go get the gogo context
	// TODO: should this version (@main) correlate to something on the local machine?
	getPath := fmt.Sprintf("%s@main", GOGOIMPORTPATH)
	err := sh.Cmd(fmt.Sprintf("go get %s", getPath)).Dir(folder).RunAndStream()
	if err != nil {
		return err
	}

	// try adding the tool deps, if we're in go1.24+. This is a better way to pin the
	// gogo context package, without requiring us to import it for side-effects somewhere
	err = tryToAddToolDep(folder)
	if err != nil {
		return err
	}

	// now go mod tidy
	err = sh.Cmd("go mod tidy").Dir(folder).RunAndStream()
	return err
}

// ensureExample adds all example files. It assumes a base folder that contains a
// folder called "pkg". If the pkg folder does not exist, it will create it.
func renderExamples(folder, pkgImportPath, gogoImportPath string) error {
	// make sure the <folder>/pkg folder exists
	err := ensureFolder(path.Join(folder, "pkg"))
	if err != nil {
		return err
	}
	// render main file
	mainFilePath := path.Join(folder, "hello.go")
	mainFile, err := os.Create(mainFilePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = mainFile.Close()
	}()

	err = HELLO_TEMPLATE.Execute(mainFile, map[string]string{
		"GoGoImportPath": pkgImportPath,
		"PkgImportPath":  gogoImportPath,
	})
	if err != nil {
		return err
	}
	// render pkg file
	pkgFilePath := path.Join(folder, "pkg", "hello.go")
	pkgFile, err := os.Create(pkgFilePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = pkgFile.Close()
	}()

	err = HELLO_PKG_TEMPLATE.Execute(pkgFile, nil)
	if err != nil {
		return err
	}

	// render pkg test file
	testPkgFilePath := path.Join(folder, "pkg", "hello_test.go")
	testPkgFile, err := os.Create(testPkgFilePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = testPkgFile.Close()
	}()

	err = HELLO_PKG_TEST_TEMPLATE.Execute(testPkgFile, nil)
	return err
}

func ensureFolder(folder string) error {
	// make sure the folder exists, and if not, create it
	if err := os.MkdirAll(folder, 0755); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

// tryToAddToolDep attempts to add the gogo context as a tool dependency, if the go version is >= 1.24.
func tryToAddToolDep(folder string) error {
	acceptable, err := version.MeetsGoVersion(REQUIRED_VERSION)
	if err != nil {
		return err
	}
	if !acceptable {
		return nil
	}

	getPath := ""
	// if we're using go 1.24, also pin the gogo ctx as a tool, so it's stuck in deps
	err = sh.Cmd(fmt.Sprintf("go get -tool %s", getPath)).Dir(folder).RunAndStream()
	return err
}
