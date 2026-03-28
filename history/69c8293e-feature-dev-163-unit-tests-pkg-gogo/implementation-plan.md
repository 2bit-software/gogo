# Implementation Plan: Unit Tests for pkg/gogo

**Created**: 2026-03-28
**Status**: Draft
**Spec**: [spec.md](spec.md)
**Research**: [research.md](research.md)

## Overview

Create `pkg/gogo/args_test.go` with comprehensive unit tests for `setFieldFromString`, `HydrateFromPositional`, and `ParseArgs`. Single file, table-driven tests, testify assertions.

## Implementation Steps

### Step 1: Add testify dependency

**Action**: Run `go get github.com/stretchr/testify` in `pkg/gogo/`
**Why**: testify is not in `pkg/gogo/go.mod` yet; needed for `assert` package
**Traces to**: FR-007

### Step 2: Create test file with test helper

**Action**: Create `pkg/gogo/args_test.go` with `package gogo`
**Contents**:
- Import `reflect`, `testing`, `math`, `github.com/stretchr/testify/assert`
- Helper function to create a `reflect.Value` from a temporary struct field for `setFieldFromString` tests:

```go
// fieldOf creates a settable reflect.Value of the given type for testing.
// Usage: field := fieldOf[int8]()
func fieldOf[T any]() reflect.Value {
    var v T
    return reflect.ValueOf(&v).Elem()
}
```

**Traces to**: FR-007, prerequisites

### Step 3: Implement TestSetFieldFromString

**Action**: Table-driven test function covering all type branches
**Structure**:

```go
func TestSetFieldFromString(t *testing.T) {
    tests := []struct {
        name      string
        field     reflect.Value
        value     string
        fieldName string
        wantErr   string // substring match; empty = no error
        check     func(t *testing.T, field reflect.Value) // optional value assertion
    }{...}
}
```

**Test cases** (grouped by type):

| Group | Cases | Spec Scenarios |
|-------|-------|---------------|
| String | set "hello", set "" (empty) | US1-1, US1-2 |
| Bool true | "true", "1", "t", "TRUE", "True", "T" | US1-3 |
| Bool false | "false", "0", "f", "FALSE", "False", "F" | US1-4 |
| Bool invalid | "yes", "no", "maybe", "" → error "invalid boolean value for" | US1-5 |
| Int | "42" success | US1-6 |
| Int invalid | "abc" → error "invalid integer value for" | US1-7 |
| Int8 boundary | "127" ok, "128" overflow, "-128" ok, "-129" overflow | US1-8 to US1-11 |
| Uint8 boundary | "255" ok, "256" overflow, "-1" → "invalid unsigned integer" | US1-12 to US1-14 |
| Float64 | "3.14" success, "abc" → error "invalid float value for" | US1-15, US1-16 |
| Float32 | "3.14" success, "3.5e38" → overflow | US1-17, US1-18 |
| Unsupported | []string slice → error "unsupported field type for" | US1-19 |

**Traces to**: FR-001, FR-002, SC-003

### Step 4: Implement TestHydrateFromPositional

**Action**: Table-driven test function covering the two-pass algorithm
**Structure**:

```go
func TestHydrateFromPositional(t *testing.T) {
    // Sub-tests grouped by concern
    t.Run("basic assignment", ...)
    t.Run("pre-set fields", ...)
    t.Run("empty quote stripping", ...)
    t.Run("input validation", ...)
    t.Run("edge cases", ...)
}
```

**Test cases by group**:

| Group | Cases | Spec Scenarios |
|-------|-------|---------------|
| Basic assignment | 3 ordered fields + 3 args → all filled | US2-1 |
| Pre-set fields | Struct {A order:"0" = "existing", B order:"1"}, args ["x","y"]. First pass: A skipped (non-zero), B="y". Second pass: no unfilled fields remain, "x" is unused. Assert: A="existing", B="y" | US2-2 |
| Empty quotes (first pass) | 2-byte `""` consumed but not set | US2-3 |
| Empty quotes (first pass) | 4-byte `\"\"` consumed but not set | US2-4 |
| Fewer args than fields | Extra fields stay zero, no error | US2-5 |
| Non-pointer input | Raw struct → error "expected pointer to struct" | US2-6 |
| Pointer to non-struct | *string → error "expected pointer to struct" | US2-7 |
| Nil input | nil → error | US2-8 |
| Invalid order tag | order:"abc" → error "invalid order tag" | US2-9 |
| Unexported fields | Skipped without error | US2-10 |
| More args than fields | Extra args unused | US2-11 |
| Empty args | No-op, no error | US2-12 |
| No order tags | No-op, no error | US2-13 |
| Gapped order tags | order:"0" and order:"3" map to positional indices | US2-14 |
| Mixed types | string + int fields with type conversion | US2-15 |

**Traces to**: FR-003, FR-004, FR-005, SC-004

### Step 5: Implement TestParseArgs

**Action**: Smoke tests for the go-flags wrapper
**Structure**: Table-driven with struct + args → check field value + positional return

| Cases | Spec Scenarios |
|-------|---------------|
| Flag sets field, no positional remaining | US3-1 |
| Flag + positional → field set, positional returned | US3-2 |
| Unknown flag → error | US3-3 |

**Traces to**: FR-006, SC-002

### Step 6: Run tests and verify

**Action**: Run `go test ./pkg/gogo/... -count=1 -v`
**Expected**: All tests pass
**Traces to**: SC-001

## Dependencies

```
Step 1 → Step 2 → Steps 3,4,5 (can be written in any order) → Step 6
```

Steps 3, 4, and 5 have no dependencies on each other — they test different functions. However, since they're all in one file, they'll be written sequentially.

## Technical Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Single test file | `args_test.go` | All three functions are in `args.go`; one test file per source file |
| Generic helper `fieldOf[T]()` | Generics-based reflect.Value creation | Cleaner than creating full structs for each setFieldFromString case |
| Subtests for HydrateFromPositional | `t.Run` groups | Organizes 15+ cases by concern rather than one flat table |
| No tests for `cli.go` or `context.go` | Skip | Type aliases and no-op stubs — zero testable logic |
| Error assertions | `assert.ErrorContains` | Checks error message substrings without brittle exact-match |

## Risks

- **Low**: `go get testify` may pull transitive deps that conflict. Mitigated: testify is already used in other modules in this repo.
- **Low**: `fieldOf[T]` helper requires Go 1.18+ generics. Confirmed: `go.mod` specifies Go 1.23.4.
