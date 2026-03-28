# Reproduction

## Steps

1. Run `go test ./pkg/fs/ -v -run TestParentDir`
2. Observe all tests are skipped

## Expected

Tests run and pass.

## Actual

```
--- SKIP: TestParentDirWithRelativesUnix (0.00s)
    parent_test.go:11: these do not work yet
--- SKIP: TestParentDirWithRelativesWindows (0.00s)
    parent_test.go:132: these do not work yet
--- SKIP: TestParentDirWithRelativesWithInvalidPaths (0.00s)
    parent_test.go:200: these do not work yet
```

## To verify bugs without skip

Remove `t.Skip(...)` lines and run tests. Two failures manifest:
1. Common parent paths are missing the leading `/` on Unix
2. Paths that don't exist on disk cause unexpected behavior via `os.Stat`
