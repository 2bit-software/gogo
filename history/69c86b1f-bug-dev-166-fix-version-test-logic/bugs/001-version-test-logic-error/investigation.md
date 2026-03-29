# Investigation

## Execution Flow

### `meetsGoVersionHelper(required, current string) (bool, error)`

1. Calls `semver.Compare(required, current)`
2. Returns `compareResult <= 0` (i.e., required <= current means version is sufficient)

### `getGoVersionString(version string) (string, error)`

1. Parses `"go version go1.XX.Y os/arch"` format
2. Strips `"go"` prefix, prepends `"v"` → returns `"v1.XX.Y"`

### `MeetsGoVersion(required string) (bool, error)`

1. Runs `go version` via `sh.Cmd`
2. Parses output with `getGoVersionString` → gets `"v1.XX.Y"`
3. Passes `required` (as-is) and parsed current to `meetsGoVersionHelper`

## The Bug

`semver.Compare` requires both arguments to have a `"v"` prefix. Strings without `"v"` are invalid and always compare less-than valid strings.

### In Tests

Tests pass `"1.24.0"` (no "v") as `required`. `getGoVersionString` returns `"v1.XX.Y"` (with "v"). Since `"1.24.0"` is invalid:

| Test | required | current | Compare result | Function returns | Correct answer |
|------|----------|---------|---------------|-----------------|---------------|
| HigherThan124 | "1.24.0" | "v1.25.0" | -1 | true | true |
| Exactly124 | "1.24.0" | "v1.24.0" | -1 | true | true |
| LowerThan124 | "1.24.0" | "v1.23.2" | -1 | **true** | **false** |

All three return `true` because the invalid string always compares less-than. Tests 1 and 2 pass by accident.

### In Production

The sole caller (`pkg/gadgets/init.go:176`) passes `REQUIRED_VERSION = "v1.24.0"` — with the "v" prefix. So production comparisons work correctly. The bug is **latent in production** but **active in tests**.

## Impact Assessment

- **Production risk**: Low — the only caller uses the correct "v" prefix format
- **Test reliability**: Broken — tests pass for the wrong reason, providing false confidence
- **API fragility**: High — any new caller that omits "v" will get silently wrong results
