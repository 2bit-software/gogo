# Bug Report: Skipped ParentDirWithRelatives Tests

## Symptoms

All 3 test functions in `pkg/fs/parent_test.go` are unconditionally skipped with `t.Skip("these do not work yet")`:

- `TestParentDirWithRelativesUnix` (6 test cases)
- `TestParentDirWithRelativesWindows` (3 test cases)
- `TestParentDirWithRelativesWithInvalidPaths` (2 test cases)

## Impact

`ParentDirWithRelatives` is marked "WIP" (`parent.go:10`) and has zero passing test coverage. This function computes common parent directories and relative paths — foundational logic for workspace discovery.

## Context

The function exists but is not yet called by any production code. It appears to be a planned optimization for directory traversal. The tests were written alongside the implementation but skipped when they failed.
