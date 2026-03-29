# Bug Report: version_test.go logic error and MeetsGoVersion always returns true

## Source

**Linear Ticket**: [DEV-166](https://linear.app/heinsight/issue/DEV-166/fix-version-testgo-logic-error-and-improve-coverage)

## Symptoms

1. **Test assertion error (line 51)**: `TestGoVersionLowerThan124` tests Go 1.23.2 against a minimum requirement of 1.24.0 and asserts `true`. Go 1.23.2 does NOT meet 1.24.0, so this should assert `false`.

2. **`MeetsGoVersion()` always returns `true`**: The public function passes `required` without a "v" prefix to `meetsGoVersionHelper`, but `getGoVersionString` returns `current` with a "v" prefix. The `semver.Compare` function treats strings without "v" as invalid, and invalid strings always sort less-than valid ones. This means `semver.Compare(required, current)` always returns `-1`, and `-1 <= 0` is always `true`.

3. **No tests for `MeetsGoVersion()`**: The public function is untested — only internal helpers have coverage.

## Affected Files

- `pkg/version/version.go` — `meetsGoVersionHelper()` and `MeetsGoVersion()`
- `pkg/version/version_test.go` — incorrect assertion and missing coverage

## Impact

Any code relying on `MeetsGoVersion()` to gate behavior on Go version will always pass, even if the installed Go version is too old.
