# Implementation Plan: pkg/sh Unit Tests

## Overview

Two deliverables: (1) fix the `RunWithWriters` bug, (2) add comprehensive unit tests for all exported functions in `pkg/sh`.

## Step 1: Fix RunWithWriters Bug

**File**: `pkg/sh/sh.go` (lines 120-128)

Current (broken):
```go
func (e *Executor) RunWithWriters(stdOut, errOut io.Writer) error {
    if stdOut == nil {
        stdOut = os.Stdout
    }
    if errOut == nil {
        e.stdErr = errOut  // assigns nil to stdErr
    }
    return e.Run()
}
```

Fixed:
```go
func (e *Executor) RunWithWriters(stdOut, errOut io.Writer) error {
    if stdOut == nil {
        stdOut = os.Stdout
    }
    if errOut == nil {
        errOut = os.Stderr
    }
    e.stdOut = stdOut
    e.stdErr = errOut
    return e.Run()
}
```

**Changes**:
- Assign `stdOut` param to `e.stdOut` (was never assigned)
- Fix `errOut` nil-check: default to `os.Stderr`, not assign nil
- Assign `errOut` param to `e.stdErr`

**Dependency**: None. Do this first so tests can verify correct behavior.

## Step 2: Create Test File

**File**: `pkg/sh/sh_test.go`

**Structure**: Single `ShTestSuite` using `testify/suite`.

### Test Groups (in order of implementation):

#### Group A: EnvMapToEnv (pure function, no side effects)
- `TestEnvMapToEnv_WithEntries` — verify KEY=VALUE formatting
- `TestEnvMapToEnv_Empty` — verify empty map returns nil

#### Group B: Constructor Functions
- `TestCmd_NoArgs` — `Cmd()` creates executor with empty command
- `TestCmd_SingleArg` — `Cmd("echo")` sets cmd field
- `TestCmd_MultipleArgs` — `Cmd("echo", "hello")` sets cmd and args
- `TestCmdWithCtx` — verify context is propagated

#### Group C: Builder Methods
- `TestDir` — set working dir, run `pwd`, verify output
- `TestSetArgs` — set args on existing executor
- `TestSetEnv` — replace environment entirely
- `TestAddEnv` — append to environment
- `TestStdin` — pipe input via `strings.NewReader`, read with `cat`

#### Group D: Command Parsing
- `TestRun_SingleStringWithSpaces` — `Cmd("echo hello")` parses correctly
- `TestRun_VariadicArgs` — `Cmd("echo", "hello")` uses args directly
- `TestRun_CommandWithSpacesAndSetArgs` — parsed cmd parts + SetArgs combine

#### Group E: Execution & Output Capture
- `TestRun_Success` — `Cmd("true").Run()` returns nil
- `TestRun_Failure` — `Cmd("false").Run()` returns error
- `TestStdOut` — captures stdout only
- `TestString` — captures combined stdout+stderr
- `TestRunWithWriters` — custom writers receive output
- `TestRunWithWriters_NilDefaults` — nil writers default to os.Stdout/Stderr
- `TestRunAndStream` — runs without error (output goes to os.Stdout/Stderr)

#### Group F: Context Cancellation
- `TestCmdWithCtx_Cancellation` — cancelled context terminates command

#### Group G: DetermineWidth
- `TestDetermineWidth_NotTerminal` — returns -1 in test environment

#### Group H: Edge Cases
- `TestRun_NonExistentDir` — Dir set to invalid path returns error
- `TestRun_EmptyCommand` — empty Cmd() returns error
- `TestSetPrintFinalCommand` — verify flag can be set (execution test)

## Step 3: Run Tests

```
go test -v -count=1 ./pkg/sh/...
```

Verify all pass, no flaky tests from timing-dependent operations.

## Dependencies

```
Step 1 (bug fix) → Step 2 (tests) → Step 3 (verify)
```

Step 2 Group A-C have no interdependencies and could be written in any order. Groups D-H depend on basic execution working (Group E).

## Risk Assessment

- **Low risk**: All tests use deterministic commands (`echo`, `true`, `false`, `cat`, `pwd`)
- **Context cancellation test**: Uses `sleep` with short timeout — keep timeout generous enough to avoid flakiness (100ms+ for cancellation)
- **Platform**: Tests use Unix commands — should work on macOS and Linux (CI)
