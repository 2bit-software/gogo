package version

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"os/exec"
	"testing"
)

func TestVersionSuite(t *testing.T) {
	suite.Run(t, new(VersionTestSuite))
}

type VersionTestSuite struct {
	suite.Suite
	originalExecCommand func(string, ...string) *exec.Cmd
}

// TestHelperProcess isn't a real test - it's used to mock exec.Command
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_TEST_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Println(os.Getenv("MOCK_OUTPUT"))
	os.Exit(0)
}

func (s *VersionTestSuite) TestGoVersionHigherThan124() {
	str := "go version go1.25.0 linux/amd64"
	version, err := getGoVersionString(str)
	assert.NoError(s.T(), err)
	result, err := meetsGoVersionHelper("1.24.0", version)
	assert.True(s.T(), result)
}

func (s *VersionTestSuite) TestGoVersion124() {
	str := "go version go1.24.0 linux/amd64"
	version, err := getGoVersionString(str)
	assert.NoError(s.T(), err)
	result, err := meetsGoVersionHelper("1.24.0", version)
	assert.True(s.T(), result)
}

func (s *VersionTestSuite) TestGoVersionLowerThan124() {
	str := "go version go1.23.2 linux/amd64"
	version, err := getGoVersionString(str)
	assert.NoError(s.T(), err)
	result, err := meetsGoVersionHelper("1.24.0", version)
	assert.True(s.T(), result)
}

func (s *VersionTestSuite) TestInvalidVersionFormat() {
	str := "go version invalid"
	_, err := getGoVersionString(str)
	assert.Error(s.T(), err)
}
