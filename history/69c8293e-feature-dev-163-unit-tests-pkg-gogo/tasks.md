# Tasks: Unit Tests for pkg/gogo

**Complexity**: Simple (3 files affected)
**Total Tasks**: 6
**Critical Path**: T001 → T002 → T003/T004/T005 (parallel) → T006

## Tasks

- [ ] T001 [US1,US2,US3] Add testify dependency to `pkg/gogo/go.mod`
  - **File**: `pkg/gogo/go.mod`, `pkg/gogo/go.sum`
  - **Action**: Run `cd pkg/gogo && go get github.com/stretchr/testify`
  - **Acceptance**: `github.com/stretchr/testify` appears in `pkg/gogo/go.mod` require block
  - **Traces to**: Plan Step 1, FR-007

- [ ] T002 Create test file with imports and `fieldOf[T]` helper in `pkg/gogo/args_test.go`
  - **File**: `pkg/gogo/args_test.go`
  - **Action**: Create file with `package gogo`, imports (`reflect`, `testing`, `math`, `testify/assert`), and generic helper:
    ```go
    func fieldOf[T any]() reflect.Value {
        var v T
        return reflect.ValueOf(&v).Elem()
    }
    ```
  - **Acceptance**: File compiles with `go build ./pkg/gogo/...`
  - **Traces to**: Plan Step 2, FR-007
  - **Depends on**: T001

- [ ] T003 [P] [US1] Implement `TestSetFieldFromString` with table-driven tests
  - **File**: `pkg/gogo/args_test.go`
  - **Action**: Add table-driven test covering all type branches. Test cases:
    - String: "hello", "" (empty) → field set (US1-1, US1-2)
    - Bool true: "true", "1", "t", "TRUE", "True", "T" → true (US1-3)
    - Bool false: "false", "0", "f", "FALSE", "False", "F" → false (US1-4)
    - Bool invalid: "yes", "no", "maybe", "" → error contains "invalid boolean value for" (US1-5)
    - Int: "42" → 42 (US1-6)
    - Int invalid: "abc" → error contains "invalid integer value for" (US1-7)
    - Int8 boundary: "127" ok, "128" overflow, "-128" ok, "-129" overflow (US1-8 to US1-11)
    - Uint8 boundary: "255" ok, "256" overflow, "-1" → "invalid unsigned integer value for" (US1-12 to US1-14)
    - Float64: "3.14" ok, "abc" → "invalid float value for" (US1-15, US1-16)
    - Float32: "3.14" ok, "3.5e38" → "overflows" (US1-17, US1-18)
    - Unsupported: `[]string` slice → "unsupported field type for" (US1-19)
  - **Acceptance**: `go test ./pkg/gogo/... -run TestSetFieldFromString -count=1` passes
  - **Traces to**: Plan Step 3, FR-001, FR-002, SC-003
  - **Depends on**: T002

- [ ] T004 [P] [US2] Implement `TestHydrateFromPositional` with grouped subtests
  - **File**: `pkg/gogo/args_test.go`
  - **Action**: Add test function with `t.Run` subgroups. Test cases:
    - **Basic assignment**: 3 ordered fields + 3 args → all filled (US2-1)
    - **Pre-set fields**: A(order:0)="existing", args ["x","y"] → A stays "existing", B="y" (US2-2)
    - **Empty quotes first pass**: 2-byte `""` (0x22,0x22) consumed but field not set (US2-3)
    - **Escaped quotes first pass**: 4-byte `\"\"` (0x5C,0x22,0x5C,0x22) consumed but field not set (US2-4)
    - **Fewer args**: fewer args than fields → extra fields zero, no error (US2-5)
    - **Non-pointer**: raw struct → error "expected pointer to struct" (US2-6)
    - **Pointer to non-struct**: *string → error "expected pointer to struct" (US2-7)
    - **Nil input**: nil → error (US2-8)
    - **Invalid order tag**: order:"abc" → error "invalid order tag for field" (US2-9)
    - **Unexported fields**: skipped without error (US2-10)
    - **More args than fields**: extra args unused, no error (US2-11)
    - **Empty args**: no-op, no error (US2-12)
    - **No order tags**: no-op, no error (US2-13)
    - **Gapped order tags**: order:"0" + order:"3", 4 args → maps to positional indices (US2-14)
    - **Mixed types**: string + int fields with conversion (US2-15)
  - **Acceptance**: `go test ./pkg/gogo/... -run TestHydrateFromPositional -count=1` passes
  - **Traces to**: Plan Step 4, FR-003, FR-004, FR-005, SC-004
  - **Depends on**: T002

- [ ] T005 [P] [US3] Implement `TestParseArgs` smoke tests
  - **File**: `pkg/gogo/args_test.go`
  - **Action**: Add table-driven smoke tests. Test cases:
    - `["--name", "foo"]` → field="foo", no positional (US3-1)
    - `["--name", "foo", "bar"]` → field="foo", positional=["bar"] (US3-2)
    - `["--unknown"]` → error (US3-3)
  - **Acceptance**: `go test ./pkg/gogo/... -run TestParseArgs -count=1` passes
  - **Traces to**: Plan Step 5, FR-006, SC-002
  - **Depends on**: T002

- [ ] T006 Run full test suite and verify all pass
  - **Action**: Run `go test ./pkg/gogo/... -count=1 -v`
  - **Acceptance**: All tests pass, zero failures
  - **Traces to**: Plan Step 6, SC-001
  - **Depends on**: T003, T004, T005

## Dependency Graph

```
T001 → T002 → T003 [P]
              → T004 [P] → T006
              → T005 [P]
```

## Traceability Matrix

| FR | Tasks |
|----|-------|
| FR-001 | T003 |
| FR-002 | T003 |
| FR-003 | T004 |
| FR-004 | T004 |
| FR-005 | T004 |
| FR-006 | T005 |
| FR-007 | T001, T002 |

| SC | Tasks |
|----|-------|
| SC-001 | T006 |
| SC-002 | T003, T004, T005 |
| SC-003 | T003 |
| SC-004 | T004 |
