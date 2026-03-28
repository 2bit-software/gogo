# Feature Specification: Unit Tests for pkg/sh

**Feature Branch**: `morganhein/dev-162-add-unit-tests-for-pkgsh-shell-execution`
**Created**: 2026-03-28
**Status**: Draft
**Input**: Linear DEV-162 — Add unit tests for pkg/sh (shell execution)

## User Scenarios & Testing

### User Story 1 - Basic Command Execution and Output Capture (Priority: P1)

A developer runs commands via `sh.Cmd()` and captures output. Tests must verify that commands execute correctly and output is captured accurately via `StdOut()` and `String()`.

**Why this priority**: This is the core functionality — if execution and output capture don't work, nothing else matters.

**Independent Test**: Run `echo hello` via `Cmd("echo", "hello").StdOut()` and verify output.

**Acceptance Scenarios**:

1. **Given** a simple command, **When** executed via `StdOut()`, **Then** stdout is captured as a string
2. **Given** a simple command, **When** executed via `String()`, **Then** combined stdout+stderr is captured
3. **Given** a failing command, **When** executed via `Run()`, **Then** a non-nil error is returned
4. **Given** a successful command, **When** executed via `Run()`, **Then** nil error is returned

---

### User Story 2 - Command Parsing (Priority: P1)

Commands can be specified as a single string with spaces (`Cmd("go mod tidy")`) or as variadic arguments (`Cmd("go", "mod", "tidy")`). The parser must handle both correctly.

**Why this priority**: Two distinct code paths in `Run()` — both are heavily used throughout the codebase.

**Independent Test**: Compare output of `Cmd("echo hello")` vs `Cmd("echo", "hello")`.

**Acceptance Scenarios**:

1. **Given** a single string command with spaces, **When** executed, **Then** it is parsed and runs correctly
2. **Given** variadic args, **When** executed, **Then** args are used as-is
3. **Given** a command string with spaces AND SetArgs called, **When** executed, **Then** parsed command parts are prepended to SetArgs values

---

### User Story 3 - Builder Methods (Priority: P2)

Builder methods (`Dir`, `SetArgs`, `SetEnv`, `AddEnv`, `Stdin`) configure the executor. Each must correctly influence execution.

**Why this priority**: These are the most-used chainable methods in the codebase.

**Independent Test**: Each builder method can be tested independently by verifying its effect on command execution.

**Acceptance Scenarios**:

1. **Given** `Dir(path)` is set, **When** command runs `pwd`, **Then** output matches the specified directory
2. **Given** `SetArgs("hello")` is called on `Cmd("echo")`, **When** executed, **Then** output is "hello"
3. **Given** `SetEnv` with a custom env, **When** command prints env, **Then** only custom env vars exist
4. **Given** `AddEnv` with extra vars, **When** command prints a specific var, **Then** the added var is present
5. **Given** `Stdin` with a reader, **When** `cat` reads from stdin, **Then** output matches input

---

### User Story 4 - Context Cancellation (Priority: P2)

`CmdWithCtx` accepts a context for cancellation. A cancelled context must terminate the running command.

**Why this priority**: Critical for timeout and graceful shutdown scenarios.

**Independent Test**: Create a cancelled context, run a long-running command, verify it returns an error.

**Acceptance Scenarios**:

1. **Given** a context that is cancelled, **When** a long-running command is executed, **Then** Run() returns a context error
2. **Given** a context with timeout, **When** command exceeds timeout, **Then** Run() returns an error

---

### User Story 5 - EnvMapToEnv Utility (Priority: P3)

`EnvMapToEnv` converts `map[string]string` to `[]string` of `KEY=VALUE` pairs.

**Why this priority**: Pure utility function, simple to test, lower risk.

**Independent Test**: Pass a map, verify output slice contains expected `KEY=VALUE` entries.

**Acceptance Scenarios**:

1. **Given** a map with entries, **When** converted, **Then** output contains `KEY=VALUE` strings for each entry
2. **Given** an empty map, **When** converted, **Then** output is nil/empty

---

### User Story 6 - DetermineWidth (Priority: P3)

`DetermineWidth` returns the terminal width or -1 if not in a terminal.

**Why this priority**: Simple function with limited testability in CI (no TTY).

**Independent Test**: In a test environment (no TTY), verify it returns -1.

**Acceptance Scenarios**:

1. **Given** stdout is not a terminal (CI/test), **When** called, **Then** returns -1

---

### Edge Cases

- What happens when `Cmd()` is called with no arguments? (empty command)
- What happens when `Dir()` is set to a non-existent directory?
- What happens when `SetEnv` is called with an empty slice?
- How does command parsing handle quoted strings? (`Cmd("echo 'hello world'")`)

## Requirements

### Functional Requirements

- **FR-001**: Tests MUST cover `Cmd()` and `CmdWithCtx()` constructor functions
- **FR-002**: Tests MUST verify output capture via `StdOut()` and `String()`
- **FR-003**: Tests MUST verify builder methods: `Dir`, `SetArgs`, `SetEnv`, `AddEnv`, `Stdin`
- **FR-004**: Tests MUST verify error propagation from failed commands
- **FR-005**: Tests MUST verify context cancellation terminates execution
- **FR-006**: Tests MUST verify command string parsing (single string with spaces)
- **FR-007**: Tests MUST verify `EnvMapToEnv` key=value formatting
- **FR-008**: Tests MUST verify `DetermineWidth` returns -1 when not in a terminal (skip TTY-dependent tests)
- **FR-009**: Tests MUST verify combined args behavior (command with spaces + SetArgs)
- **FR-010**: `RunWithWriters` MUST be fixed — stdOut writer must be assigned to executor, and errOut nil-check logic must be corrected

## Success Criteria

### Measurable Outcomes

- **SC-001**: All exported functions in `pkg/sh` have at least one test
- **SC-002**: Tests pass in CI (no terminal dependency)
- **SC-003**: `go test ./pkg/sh/...` runs successfully with zero failures

## Testing Requirements

### Test Strategy

- Single test file: `pkg/sh/sh_test.go`
- Use `testify/suite` with a `ShTestSuite` struct (per project conventions)
- Use real OS commands (`echo`, `true`, `false`, `cat`) for execution tests
- Table-driven subtests where multiple inputs test the same behavior
- No mocking of os/exec — tests verify actual shell execution

### FR to Test Mapping

| FR | Test Type | Description |
|----|-----------|-------------|
| FR-001 | Unit | Verify Cmd/CmdWithCtx create valid executors |
| FR-002 | Unit | Verify StdOut/String capture command output |
| FR-003 | Unit | Verify each builder method affects execution |
| FR-004 | Unit | Verify failed commands return non-nil error |
| FR-005 | Unit | Verify cancelled context stops command |
| FR-006 | Unit | Verify space-separated string is parsed into cmd+args |
| FR-007 | Unit | Verify EnvMapToEnv produces KEY=VALUE pairs |
| FR-008 | Unit | Verify DetermineWidth returns -1 in non-TTY |
| FR-009 | Unit | Verify parsed command parts + SetArgs combine correctly |

### Edge Case Coverage

- Empty command string -> test for error or specific behavior
- Non-existent directory in Dir() -> test for error from Run()
- Quoted strings in single-command parsing -> verify mvdan/sh parsing
