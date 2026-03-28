---
status: complete
updated: 2026-03-28
---

# Research: pkg/sh Unit Tests

## Executive Summary

`pkg/sh` is a 2-file, ~215 LOC shell execution package with 15 exported symbols and zero test coverage. It's used throughout the codebase as the primary execution layer. The project uses testify/suite v1.9.0 for testing, and the package's builder pattern with pure functions makes it highly testable.

## Findings

### Codebase Context

- **Package structure**: Two files — `sh.go` (Executor builder + execution) and `info.go` (terminal width detection)
- **Builder pattern**: All builder methods return `*Executor` for chaining. Terminal methods are `Run()`, `StdOut()`, `String()`, `RunAndStream()`, `RunWithWriters()`
- **Command parsing**: When no args are provided and the command string contains spaces, `mvdan/sh/shell.Fields()` is used to parse the command. When args are provided and the command has spaces, it parses the command and prepends those parts to the existing args.
- **Environment handling**: `CmdWithCtx` initializes env from `os.Environ()`. `SetEnv` replaces entirely, `AddEnv` appends.
- **Consumers**: `pkg/gadgets`, `pkg/version`, `cmd/gogo/cmds/config.go`, `.gogo/pkg/git.go`, `.gogo/pkg/compile.go`, `gogo_cli_test.go`
- **Usage patterns observed**: Single-string commands (`Cmd("go mod tidy")`), variadic args (`Cmd("go", "mod", "tidy")`), chained Dir+Run, StdOut capture, String capture, AddEnv with String
- **Existing test style**: `pkg/version` uses testify/suite; `pkg/fs` uses table-driven tests with stdlib. Project rules mandate testify/suite + testify/assert.

### Domain Knowledge

- **Testing shell execution**: Use simple, deterministic commands (`echo`, `true`, `false`, `cat`) for unit tests
- **Context cancellation**: Use `context.WithCancel` or `context.WithTimeout` plus `sleep` to verify cancellation
- **Stdin piping**: Use `strings.NewReader` to inject stdin, verify with `cat`
- **Environment testing**: Set unique env vars and echo them to verify propagation
- **DetermineWidth**: Depends on real terminal via `golang.org/x/term` — in CI/test, stdout is typically not a TTY, so it returns -1. Testing the non-terminal path is reliable; testing the terminal path requires mocking or is best left to integration tests.

### Bug Spotted During Research

`RunWithWriters` has a likely bug on line 125-126 of sh.go:
```go
if errOut == nil {
    e.stdErr = errOut  // sets stdErr to nil when errOut is nil
}
```
This should probably be:
```go
if errOut != nil {
    e.stdErr = errOut
}
```
And the stdOut writer is never assigned to `e.stdOut`. The function sets local `stdOut` to `os.Stdout` if nil, but never passes it to the executor. Tests should document this behavior.

## Decision Points

- [x] **D1**: Use testify/suite (per project conventions)
- [x] **D2**: Use real commands (`echo`, `true`, `false`, `cat`) rather than mock exec — the package wraps os/exec and tests should verify actual execution
- [x] **D3**: Skip terminal-dependent DetermineWidth tests — CI has no TTY, test non-terminal path only
- [x] **D4**: Fix the RunWithWriters bug in this ticket

## Recommendations

1. **Test file structure**: Single file `pkg/sh/sh_test.go` with a `ShTestSuite` using testify/suite
2. **Use real commands**: `echo`, `true`, `false`, `cat`, `env` for deterministic, cross-platform-ish tests
3. **Test DetermineWidth minimally**: Verify it returns -1 when not in a terminal (the CI case)
4. **Document RunWithWriters bug**: Write a test that captures the current (broken) behavior, so the fix can be done in a follow-up

## Sources

- `pkg/sh/sh.go` — main executor implementation
- `pkg/sh/info.go` — DetermineWidth
- `pkg/version/version_test.go` — testify/suite pattern reference
- `pkg/gadgets/builder_test.go` — test helper patterns
- Linear DEV-162 — ticket requirements
