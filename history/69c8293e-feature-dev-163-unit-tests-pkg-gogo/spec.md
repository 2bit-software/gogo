# Feature Specification: Unit Tests for pkg/gogo

**Feature Branch**: `morganhein/dev-163-add-unit-tests-for-pkggogo-args-cli-context`
**Created**: 2026-03-28
**Status**: Draft
**Input**: User description: "let's work on dev-163"

## Prerequisites

- Add `github.com/stretchr/testify` to `pkg/gogo/go.mod` (run `go get github.com/stretchr/testify` inside `pkg/gogo/`)
- Test file: `pkg/gogo/args_test.go` using `package gogo` (same-package access required for unexported `setFieldFromString`)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - setFieldFromString covers all type conversions (Priority: P1)

A developer refactoring `setFieldFromString` can run unit tests to verify all supported type conversions work correctly, including boundary values and invalid inputs.

**Why this priority**: This is the core type conversion logic. Bugs here silently corrupt argument values across all gadgets.

**Independent Test**: Can be tested in isolation by constructing `reflect.Value` fields of the target type via a temporary struct and calling `setFieldFromString` directly.

**Note**: `setFieldFromString` is unexported but accessible from tests in `package gogo`.

**Acceptance Scenarios**:

1. **Given** a `string` field, **When** called with `"hello"`, **Then** the field is set to `"hello"`
2. **Given** a `string` field, **When** called with `""` (empty Go string), **Then** the field is set to `""` (empty)
3. **Given** a `bool` field, **When** called with any of `"true"`, `"1"`, `"t"`, `"TRUE"`, `"True"`, `"T"`, **Then** the field is set to `true` (these are the values `strconv.ParseBool` accepts)
4. **Given** a `bool` field, **When** called with any of `"false"`, `"0"`, `"f"`, `"FALSE"`, `"False"`, `"F"`, **Then** the field is set to `false`
5. **Given** a `bool` field, **When** called with `"yes"`, `"no"`, `"maybe"`, or `""` (empty Go string), **Then** an error containing `"invalid boolean value for"` is returned
6. **Given** an `int` field, **When** called with `"42"`, **Then** the field is set to `42`
7. **Given** an `int` field, **When** called with `"abc"`, **Then** an error containing `"invalid integer value for"` is returned
8. **Given** an `int8` field, **When** called with `"127"`, **Then** the field is set to `127`
9. **Given** an `int8` field, **When** called with `"128"`, **Then** an error containing `"overflows"` is returned
10. **Given** an `int8` field, **When** called with `"-128"`, **Then** the field is set to `-128`
11. **Given** an `int8` field, **When** called with `"-129"`, **Then** an error containing `"overflows"` is returned
12. **Given** a `uint8` field, **When** called with `"255"`, **Then** the field is set to `255`
13. **Given** a `uint8` field, **When** called with `"256"`, **Then** an error containing `"overflows"` is returned
14. **Given** a `uint8` field, **When** called with `"-1"`, **Then** an error containing `"invalid unsigned integer value for"` is returned (strconv.ParseUint rejects negative values)
15. **Given** a `float64` field, **When** called with `"3.14"`, **Then** the field is set to `3.14`
16. **Given** a `float64` field, **When** called with `"abc"`, **Then** an error containing `"invalid float value for"` is returned
17. **Given** a `float32` field, **When** called with `"3.14"`, **Then** the field is set (within float32 precision)
18. **Given** a `float32` field, **When** called with `"3.5e38"` (exceeds `math.MaxFloat32` ~3.4e38), **Then** an error containing `"overflows"` is returned
19. **Given** a field of unsupported type (e.g., `[]string` slice), **When** called, **Then** an error containing `"unsupported field type for"` is returned

**Type coverage note**: Test at least one representative from each kind group: `string`, `bool`, one `int` variant, one `uint` variant, one `float` variant. Boundary/overflow tests for `int8` and `uint8` are sufficient to prove the overflow check works — the same `OverflowInt`/`OverflowUint` mechanism applies to all sizes. Add `int16`/`uint16` only if you want extra confidence.

---

### User Story 2 - HydrateFromPositional fills struct fields correctly (Priority: P1)

A developer modifying the two-pass assignment algorithm can verify that positional args are correctly assigned to struct fields based on `order` tags.

**Why this priority**: This is the core positional argument routing logic. The two-pass algorithm is complex and non-obvious.

**Independent Test**: Construct structs with `order` tags and call `HydrateFromPositional` with various positional arg slices.

**Two-pass algorithm summary**:
- **First pass**: Iterates fields sorted by `order` tag value (ascending). For each unset field, tries to assign `positional[field.order]`. Skips fields that are already non-zero. Strips empty-quote values (`""` and `\"\"`).
- **Second pass**: Iterates the same sorted fields again. For each still-unset field, assigns the next unused positional arg sequentially. Does NOT perform empty-quote stripping in the second pass.

**Acceptance Scenarios**:

1. **Given** a struct `{A string order:"0"; B string order:"1"; C string order:"2"}` and positional args `["x", "y", "z"]`, **When** called, **Then** A="x", B="y", C="z"
2. **Given** a struct where field A (order:"0") is pre-set to `"existing"` and positional args `["x", "y"]`, **When** called, **Then** A remains `"existing"`, B="y", and the arg "x" is available — if there's an unfilled field, it may receive "x" in the second pass
3. **Given** a positional arg that is a 2-byte Go string containing two literal double-quote characters (byte content: `0x22 0x22`), **When** processed in the **first pass**, **Then** the arg is consumed (marked as used) but the target field is NOT set
4. **Given** a positional arg that is a 4-byte Go string containing `\`, `"`, `\`, `"` (byte content: `0x5C 0x22 0x5C 0x22`), **When** processed in the **first pass**, **Then** same behavior as scenario 3 — consumed but not set
5. **Given** fewer positional args than ordered fields, **When** called, **Then** extra fields remain at zero value, no error returned
6. **Given** a non-pointer input (e.g., a raw struct value), **When** called, **Then** an error containing `"expected pointer to struct, got"` is returned
7. **Given** a pointer to a non-struct (e.g., `*string`), **When** called, **Then** an error containing `"expected pointer to struct, got"` is returned
8. **Given** a nil input, **When** called, **Then** an error is returned (reflect.ValueOf(nil).Kind() is Invalid, not Ptr)
9. **Given** a struct with an invalid order tag (e.g., `order:"abc"`), **When** called, **Then** an error containing `"invalid order tag for field"` is returned
10. **Given** a struct with unexported fields that have order tags, **When** called, **Then** unexported fields are skipped without error
11. **Given** a struct with 2 ordered fields and 3 positional args, **When** called, **Then** the first two args fill the fields by position (first pass), and the third arg is unused (no extra fields to fill in second pass)
12. **Given** an empty positional args slice, **When** called, **Then** no fields are modified, no error
13. **Given** a struct with no `order` tags on any fields, **When** called, **Then** no fields are modified, no error
14. **Given** a struct with order tags that have gaps (e.g., `order:"0"` and `order:"3"`) and 4 positional args, **When** called, **Then** the first field gets `positional[0]` and the second field gets `positional[3]` (maps to the position index, not sequential)
15. **Given** a struct with mixed types `{Name string order:"0"; Count int order:"1"}` and args `["hello", "42"]`, **When** called, **Then** Name="hello", Count=42 (type conversion applies per field)

**Important note on duplicate order values**: `sort.Slice` is NOT stable in Go. If two fields share the same `order` tag value, their processing order is non-deterministic. Tests should avoid duplicate order values.

---

### User Story 3 - ParseArgs wraps go-flags correctly (Priority: P2)

A developer can verify that `ParseArgs` correctly delegates to `go-flags` and returns remaining positional args and errors.

**Why this priority**: This is a thin one-line wrapper. Smoke tests provide regression safety without over-testing the third-party library.

**Acceptance Scenarios**:

1. **Given** a struct with `long:"name"` string field and args `["--name", "foo"]`, **When** `ParseArgs` is called, **Then** the field is set to `"foo"` and the returned positional slice is empty
2. **Given** args with both flags and positional args `["--name", "foo", "bar"]`, **When** called, **Then** the field is set to `"foo"` and `["bar"]` is returned as positional args
3. **Given** an unknown flag `["--unknown"]`, **When** called, **Then** an error is returned

---

### Edge Cases

- Empty positional args slice: no-op, no error (scenario 2.12)
- Struct with no order tags: no-op, no error (scenario 2.13)
- Gaps in order tags: maps to positional index, not sequential (scenario 2.14)
- Duplicate order values: non-deterministic due to unstable sort — avoid in tests
- Second pass does NOT strip empty quotes — only first pass does (design note, not a test target unless flagged as a bug)
- Pre-set fields (non-zero) skipped in first pass via `IsZero()` check (scenario 2.2)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Tests MUST cover all type branches in `setFieldFromString`: string, bool, at least one int variant, at least one uint variant, at least one float variant, and the unsupported-type default branch
- **FR-002**: Tests MUST cover overflow detection for at least `int8` and `uint8` (boundary values), plus `float32` overflow
- **FR-003**: Tests MUST cover HydrateFromPositional's first-pass (position-indexed) and second-pass (sequential leftover) assignment
- **FR-004**: Tests MUST cover input validation: non-pointer, pointer-to-non-struct, nil, invalid order tags
- **FR-005**: Tests MUST cover empty-quote stripping in the first pass: both the 2-byte `""` (two double-quotes) and 4-byte `\"\"` (backslash-quote-backslash-quote) variants
- **FR-006**: Tests MUST cover ParseArgs flag parsing delegation and positional arg return (smoke tests)
- **FR-007**: Tests MUST use `testify/assert` with table-driven test patterns (add testify to `pkg/gogo/go.mod`)

### Key Entities

- **Test structs**: Structs with various field types and `order`/`long`/`short` tags, defined in the test file
- **Positional args**: `[]string` slices representing command-line positional arguments

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All tests pass with `go test ./pkg/gogo/... -count=1`
- **SC-002**: `ParseArgs` and `HydrateFromPositional` (exported) and `setFieldFromString` (unexported) each have dedicated test functions
- **SC-003**: Every type branch in `setFieldFromString` has at least one success and one failure test case
- **SC-004**: Every error return path in `HydrateFromPositional` has a test case

## Testing Requirements *(mandatory)*

### Test Strategy

- **File**: `pkg/gogo/args_test.go` with `package gogo`
- **Framework**: `testify/assert` (must be added to `pkg/gogo/go.mod`)
- **Pattern**: Table-driven tests with `t.Run` subtests
- **Scope**: Pure unit tests — no binary compilation, no external dependencies beyond testify

### FR to Test Mapping

| FR | Test Type | Description |
|----|-----------|-------------|
| FR-001 | Unit | Table-driven tests for each type in setFieldFromString |
| FR-002 | Unit | Boundary value tests for int8, uint8, float32 overflow |
| FR-003 | Unit | Struct configurations testing first-pass and second-pass assignment |
| FR-004 | Unit | Error path tests for invalid inputs to HydrateFromPositional |
| FR-005 | Unit | Tests with exact byte content for empty-quote variants |
| FR-006 | Unit | Smoke tests for ParseArgs delegation |
| FR-007 | N/A | Convention requirement — verified by code review |

### Edge Case Coverage

- Empty positional args slice -> Verified no-op behavior
- Struct with no order tags -> Verified no-op behavior
- Unexported fields -> Verified they are skipped
- Pre-set fields (non-zero) -> Verified they are skipped in first pass
- Overflow boundaries for int8, uint8, float32 -> Verified error returned
- Empty-quote stripping in first pass only -> Verified both variants consumed but field not set
