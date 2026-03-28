# Initiative: sh-unit-tests

**Type**: feature
**Status**: completed
**Created**: 2026-03-28
**ID**: 69c81fd0-feature-sh-unit-tests

## Steps

| Step | Status | Updated |
|------|--------|--------|
| spec | completed | 2026-03-28 11:37 |
| plan | completed | 2026-03-28 11:42 |
| tasks | completed | 2026-03-28 11:45 |
| implement | completed | 2026-03-28 11:48 |

## Source

**Linear Ticket**: [DEV-162](https://linear.app/heinsight/issue/DEV-162/add-unit-tests-for-pkgsh-shell-execution)
**Title**: Add unit tests for pkg/sh (shell execution)

## Description

<!-- Add a description of this initiative -->

## Goals

<!-- Define the goals for this initiative -->

## Progress

<!-- Track progress here -->

## Completion

**Completed**: 2026-03-28
**Duration**: Same day

### Outcomes
- Bug fix: Fixed `RunWithWriters` — stdOut never assigned, errOut nil-check inverted
- Bug fix: Fixed `Run` — `c.Stderr` was set to `e.stdOut` instead of `e.stdErr`
- Feature: Added 26 unit tests covering all 15 exported symbols in pkg/sh
- All existing tests pass (no regressions)
