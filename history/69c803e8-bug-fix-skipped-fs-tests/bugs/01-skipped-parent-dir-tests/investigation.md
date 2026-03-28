# Investigation

## Root Causes

Two bugs in `ParentDirWithRelatives` (`pkg/fs/parent.go`):

### Bug 1: `os.Stat` on non-existent paths (line 49-52)

```go
fileInfo, err := os.Stat(p)
if err == nil && !fileInfo.IsDir() {
    p = filepath.Dir(p)
}
```

The function calls `os.Stat(p)` to check if a path is a file vs directory, so it can strip the filename. Tests use paths like `/home/user/file.txt` that don't exist on the test machine. When stat fails, the code path is skipped — but this is harmless for the algorithm because the filename becomes an extra component that naturally diverges during common prefix computation.

The real problem: this makes the function depend on filesystem state when it should be pure path computation. It also uses `p` after stripping the volume name, which on Unix is the same path but conceptually wrong.

**Fix**: Remove the `os.Stat` block entirely. The algorithm works correctly without it.

### Bug 2: `filepath.Join` drops leading `/` (lines 76-79)

```go
commonParent := filepath.Join(commonComponents...)
if !strings.HasPrefix(commonParent, filepath.VolumeName(commonParent)) {
    commonParent = string(filepath.Separator) + commonParent
}
```

Splitting `/home/user/docs` by `/` gives `["", "home", "user", "docs"]`. Then `filepath.Join("", "home", "user", "docs")` produces `"home/user/docs"` — no leading slash.

The guard on line 77 tries to fix this, but `filepath.VolumeName("home/user/docs")` returns `""` on Unix, and `strings.HasPrefix(anything, "")` is always `true`. So the condition `!true` = `false`, and the leading `/` is **never** prepended.

**Fix**: Replace `filepath.Join` with `strings.Join` which preserves the empty first component, producing `/home/user/docs`. Handle the edge case where common parent is root (`[""]` joins to `""`, should be `/`).

## Callers

`ParentDirWithRelatives` has zero callers in production code — only used in tests. No regression risk from changing behavior.
