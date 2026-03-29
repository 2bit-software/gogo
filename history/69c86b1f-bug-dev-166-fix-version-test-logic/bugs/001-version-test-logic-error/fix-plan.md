# Fix Plan

## Changes Required

### 1. Normalize inputs in `meetsGoVersionHelper` (`pkg/version/version.go`)

Add "v" prefix normalization so callers don't need to know about `semver.Compare`'s requirements:

```go
func meetsGoVersionHelper(required, current string) (bool, error) {
    if !strings.HasPrefix(required, "v") {
        required = "v" + required
    }
    if !strings.HasPrefix(current, "v") {
        current = "v" + current
    }
    if !semver.IsValid(required) {
        return false, fmt.Errorf("invalid required version: %s", required)
    }
    if !semver.IsValid(current) {
        return false, fmt.Errorf("invalid current version: %s", current)
    }
    compareResult := semver.Compare(required, current)
    return compareResult <= 0, nil
}
```

### 2. Fix test assertion (`pkg/version/version_test.go:51`)

Change `assert.True` to `assert.False` and add `assert.NoError` for consistency:

```go
func (s *VersionTestSuite) TestGoVersionLowerThan124() {
    str := "go version go1.23.2 linux/amd64"
    version, err := getGoVersionString(str)
    assert.NoError(s.T(), err)
    result, err := meetsGoVersionHelper("1.24.0", version)
    assert.NoError(s.T(), err)
    assert.False(s.T(), result)
}
```

### 3. Add edge case tests (`pkg/version/version_test.go`)

- Same version (e.g., 1.24.0 vs 1.24.0) → `true`
- Patch version difference (1.24.1 vs 1.24.0 required) → `true`
- Required with "v" prefix → still works
- Invalid version string → returns error
- Pre-release versions if applicable

### 4. Add `MeetsGoVersion` test (if feasible)

`MeetsGoVersion` depends on `sh.Cmd("go version")`. Testing options:
- **Option A**: Test against the actual Go runtime (integration-style) — simplest, verifies real behavior
- **Option B**: Refactor to inject the command executor — more testable but scope creep

Recommend **Option A**: call `MeetsGoVersion("v1.0.0")` (should return `true` since any modern Go meets 1.0) and `MeetsGoVersion("v99.0.0")` (should return `false`).

## Verification

After applying fixes:
1. `TestGoVersionLowerThan124` should pass with `assert.False`
2. All existing tests should still pass
3. New edge case tests should pass
4. `go test -v -run "TestVersionSuite" -count=1 ./pkg/version/`
