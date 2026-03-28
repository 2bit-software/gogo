# Initiative: fix-skipped-fs-tests

**Type**: bug
**Status**: completed
**Created**: 2026-03-28
**ID**: 69c803e8-bug-fix-skipped-fs-tests

## Steps

| Step | Status | Updated |
|------|--------|--------|
| investigate | completed | 2026-03-28 09:42 |
| fix | completed | 2026-03-28 09:47 |
| verify | completed | 2026-03-28 09:48 |

## Source

**Linear Ticket**: [DEV-136](https://linear.app/heinsight/issue/DEV-136/fix-skipped-filesystem-traversal-tests)
**Title**: Fix skipped filesystem traversal tests

## Description

Fix skipped filesystem traversal tests in `pkg/fs/parent_test.go`. All tests for `ParentDirWithRelatives` are skipped with `t.Skip("these do not work yet")`. Two implementation bugs prevent the tests from passing.

## Completion

**Completed**: 2026-03-28
**Duration**: Same day

### Outcomes
- Bug: Removed `os.Stat` dependency from `ParentDirWithRelatives` — function is now pure path computation
- Bug: Fixed `filepath.Join` dropping leading `/` on Unix — replaced with `strings.Join`
- Tests: Unskipped 3 test functions (8 test cases now passing)
- Tests: Fixed null byte test expectation to match pure path computation behavior
- Cleanup: Removed WIP comment and unused `os` import

### Files Changed
- `pkg/fs/parent.go` — implementation fix
- `pkg/fs/parent_test.go` — unskipped tests, fixed expectation
