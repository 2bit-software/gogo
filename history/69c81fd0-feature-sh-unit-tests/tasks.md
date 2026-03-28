# Tasks: pkg/sh Unit Tests

**Complexity**: Simple (2 files, ~350 LOC)
**Critical Path**: T001 → T002 → T003-T008 (parallel) → T009

## Task List

- [ ] T001 [US-] Fix RunWithWriters bug in `pkg/sh/sh.go` (lines 120-128)
  - Assign stdOut param to e.stdOut
  - Fix errOut nil-check to default to os.Stderr
  - Assign errOut param to e.stdErr
  - **Acceptance**: RunWithWriters correctly writes to provided writers

- [ ] T002 [US-] Create test suite scaffold in `pkg/sh/sh_test.go`
  - ShTestSuite struct with testify/suite
  - Suite runner function
  - **Acceptance**: `go test ./pkg/sh/...` compiles and runs (0 tests)

- [ ] T003 [P] [US5] Add EnvMapToEnv tests in `pkg/sh/sh_test.go`
  - Test with entries → KEY=VALUE format
  - Test empty map → nil/empty
  - **Acceptance**: FR-007 covered

- [ ] T004 [P] [US2,US3] Add constructor and builder method tests in `pkg/sh/sh_test.go`
  - Cmd/CmdWithCtx constructors
  - Dir, SetArgs, SetEnv, AddEnv, Stdin
  - **Acceptance**: FR-001, FR-003 covered

- [ ] T005 [P] [US2] Add command parsing tests in `pkg/sh/sh_test.go`
  - Single string with spaces
  - Variadic args
  - Command with spaces + SetArgs
  - **Acceptance**: FR-006, FR-009 covered

- [ ] T006 [P] [US1] Add execution and output capture tests in `pkg/sh/sh_test.go`
  - Run success/failure
  - StdOut, String capture
  - RunWithWriters with custom writers and nil defaults
  - RunAndStream
  - **Acceptance**: FR-002, FR-004, FR-010 covered

- [ ] T007 [P] [US4] Add context cancellation test in `pkg/sh/sh_test.go`
  - Cancelled context terminates command
  - **Acceptance**: FR-005 covered

- [ ] T008 [P] [US6] Add DetermineWidth and edge case tests in `pkg/sh/sh_test.go`
  - DetermineWidth returns -1 in non-TTY
  - Non-existent dir, empty command
  - **Acceptance**: FR-008 covered

- [ ] T009 [US-] Run tests and verify all pass
  - `go test -v -count=1 ./pkg/sh/...`
  - **Acceptance**: All tests pass, no flakiness

## Dependency Graph

```
T001 (bug fix)
  └─→ T002 (scaffold)
        └─→ T003, T004, T005, T006, T007, T008 (parallel)
              └─→ T009 (verify)
```

## FR Coverage

| FR | Task(s) |
|----|---------|
| FR-001 | T004 |
| FR-002 | T006 |
| FR-003 | T004 |
| FR-004 | T006 |
| FR-005 | T007 |
| FR-006 | T005 |
| FR-007 | T003 |
| FR-008 | T008 |
| FR-009 | T005 |
| FR-010 | T006 |
