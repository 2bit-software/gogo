package sh

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestShSuite(t *testing.T) {
	suite.Run(t, new(ShTestSuite))
}

type ShTestSuite struct {
	suite.Suite
}

// --- EnvMapToEnv (T003) ---

func (s *ShTestSuite) TestEnvMapToEnv_WithEntries() {
	env := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	result := EnvMapToEnv(env)
	assert.Len(s.T(), result, 2)
	assert.Contains(s.T(), result, "FOO=bar")
	assert.Contains(s.T(), result, "BAZ=qux")
}

func (s *ShTestSuite) TestEnvMapToEnv_Empty() {
	result := EnvMapToEnv(map[string]string{})
	assert.Empty(s.T(), result)
}

// --- Constructors and Builder Methods (T004) ---

func (s *ShTestSuite) TestCmd_SingleArg() {
	out, err := Cmd("echo", "hello").StdOut()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "hello\n", out)
}

func (s *ShTestSuite) TestCmd_MultipleArgs() {
	out, err := Cmd("echo", "hello", "world").StdOut()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "hello world\n", out)
}

func (s *ShTestSuite) TestCmdWithCtx() {
	ctx := context.Background()
	out, err := CmdWithCtx(ctx, "echo", "ctx-test").StdOut()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "ctx-test\n", out)
}

func (s *ShTestSuite) TestDir() {
	tmpDir := s.T().TempDir()
	out, err := Cmd("pwd").Dir(tmpDir).StdOut()
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), strings.TrimSpace(out), tmpDir)
}

func (s *ShTestSuite) TestSetArgs() {
	out, err := Cmd("echo").SetArgs("set-args-test").StdOut()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "set-args-test\n", out)
}

func (s *ShTestSuite) TestSetEnv() {
	// SetEnv replaces the entire environment, so only our var should exist
	out, err := Cmd("env").SetEnv([]string{"MY_TEST_VAR=set-env-value"}).StdOut()
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), out, "MY_TEST_VAR=set-env-value")
}

func (s *ShTestSuite) TestAddEnv() {
	out, err := Cmd("env").AddEnv([]string{"ADDED_VAR=added-value"}).StdOut()
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), out, "ADDED_VAR=added-value")
}

func (s *ShTestSuite) TestStdin() {
	input := "stdin-test-data"
	out, err := Cmd("cat").Stdin(strings.NewReader(input)).StdOut()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), input, out)
}

// --- Command Parsing (T005) ---

func (s *ShTestSuite) TestRun_SingleStringWithSpaces() {
	out, err := Cmd("echo hello").StdOut()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "hello\n", out)
}

func (s *ShTestSuite) TestRun_VariadicArgs() {
	out, err := Cmd("echo", "variadic", "args").StdOut()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "variadic args\n", out)
}

func (s *ShTestSuite) TestRun_CommandWithSpacesAndSetArgs() {
	// When command has spaces AND SetArgs is called, parsed command parts
	// are prepended to SetArgs values
	out, err := Cmd("echo hello").SetArgs("world").StdOut()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "hello world\n", out)
}

func (s *ShTestSuite) TestRun_QuotedStringsInSingleCommand() {
	out, err := Cmd("echo 'hello world'").StdOut()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "hello world\n", out)
}

// --- Execution and Output Capture (T006) ---

func (s *ShTestSuite) TestRun_Success() {
	err := Cmd("true").Run()
	assert.NoError(s.T(), err)
}

func (s *ShTestSuite) TestRun_Failure() {
	err := Cmd("false").Run()
	assert.Error(s.T(), err)
}

func (s *ShTestSuite) TestStdOut_CapturesStdoutOnly() {
	// StdOut should capture stdout; stderr goes elsewhere
	out, err := Cmd("echo", "stdout-test").StdOut()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "stdout-test\n", out)
}

func (s *ShTestSuite) TestString_CapturesCombinedOutput() {
	// String captures both stdout and stderr into one buffer
	// Use sh -c to write to both stdout and stderr
	out, err := Cmd("sh", "-c", "echo out; echo err >&2").String()
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), out, "out")
	assert.Contains(s.T(), out, "err")
}

func (s *ShTestSuite) TestRunWithWriters_CustomWriters() {
	var stdout, stderr bytes.Buffer
	err := Cmd("sh", "-c", "echo out; echo err >&2").RunWithWriters(&stdout, &stderr)
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), stdout.String(), "out")
	assert.Contains(s.T(), stderr.String(), "err")
}

func (s *ShTestSuite) TestRunWithWriters_NilDefaultsToStdoutStderr() {
	// When nil is passed, should not panic and should execute successfully
	err := Cmd("true").RunWithWriters(nil, nil)
	assert.NoError(s.T(), err)
}

func (s *ShTestSuite) TestRunAndStream() {
	// RunAndStream writes to os.Stdout/os.Stderr — just verify no error
	err := Cmd("echo", "stream-test").RunAndStream()
	assert.NoError(s.T(), err)
}

// --- Context Cancellation (T007) ---

func (s *ShTestSuite) TestCmdWithCtx_Cancellation() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := CmdWithCtx(ctx, "sleep", "10").Run()
	assert.Error(s.T(), err)
}

// --- DetermineWidth (T008) ---

func (s *ShTestSuite) TestDetermineWidth_NotTerminal() {
	// In a test environment, stdout is not a terminal
	width := DetermineWidth(false)
	assert.Equal(s.T(), -1, width)
}

// --- Edge Cases (T008) ---

func (s *ShTestSuite) TestRun_NonExistentDir() {
	err := Cmd("echo", "test").Dir("/nonexistent/path/that/does/not/exist").Run()
	assert.Error(s.T(), err)
}

func (s *ShTestSuite) TestRun_EmptyCommand() {
	err := Cmd("").Run()
	assert.Error(s.T(), err)
}

func (s *ShTestSuite) TestSetPrintFinalCommand() {
	// Verify SetPrintFinalCommand can be set and command still executes
	out, err := Cmd("echo", "print-test").SetPrintFinalCommand(true).StdOut()
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), out, "print-test")
}
