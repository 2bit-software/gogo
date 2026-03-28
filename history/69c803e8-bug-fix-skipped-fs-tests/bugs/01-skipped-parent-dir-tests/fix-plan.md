# Fix Plan

## Changes Required

### 1. Remove `os.Stat` dependency (`parent.go`)

- Remove `"os"` from imports
- Delete lines 49-52 (the `os.Stat` block)
- The algorithm works correctly without it — filenames become extra components that diverge naturally

### 2. Fix common parent path reconstruction (`parent.go`)

Replace:
```go
commonParent := filepath.Join(commonComponents...)
if !strings.HasPrefix(commonParent, filepath.VolumeName(commonParent)) {
    commonParent = string(filepath.Separator) + commonParent
}
```

With:
```go
commonParent := strings.Join(commonComponents, string(filepath.Separator))
if commonParent == "" {
    commonParent = string(filepath.Separator)
}
```

`strings.Join` preserves the empty first component from splitting absolute paths, correctly producing `/home/user/docs`. The empty check handles the root directory edge case.

### 3. Unskip tests (`parent_test.go`)

- Remove `t.Skip("these do not work yet")` from `TestParentDirWithRelativesUnix` (line 11)
- Remove `t.Skip("these do not work yet")` from `TestParentDirWithRelativesWindows` (line 132) — keep the platform skip on line 133
- Remove `t.Skip("these do not work yet")` from `TestParentDirWithRelativesWithInvalidPaths` (line 200)
- Adjust invalid path test: null byte path (`\x00`) may not cause an error with pure path computation — verify and update expectation if needed

### 4. Remove WIP comment (`parent.go`)

- Remove line 10: `// WIP: This is in progress and is an optimization`

## Test Verification

Run `go test ./pkg/fs/ -v -run TestParentDir` and confirm all non-platform-skipped tests pass.
