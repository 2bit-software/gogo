# Task Execution - Business Specification

## Overview

Task Execution is where GoGo delivers its core value — running developer-authored tasks instantly from the command line. GoGo compiles tasks into fast binaries, caches them intelligently to avoid unnecessary rebuilds, and provides options for pre-building and optimizing binaries. The developer experience is: name a task, pass arguments, get results.

## Capabilities Summary

| Capability | Description | Priority |
|------------|-------------|----------|
| Run a task with arguments | Execute any task by name with positional or named arguments | Core |
| Smart caching | Reuse compiled binaries when source files haven't changed | Core |
| Force rebuild | Bypass cache and recompile from scratch | Supporting |
| Pre-build binaries | Compile tasks ahead of time for instant execution | Supporting |
| Optimized builds | Create smaller binaries for distribution | Supporting |

---

## 1. Business Purpose

### Problem Solved
Developers need to run automation tasks quickly and reliably. Interpreted task runners add startup overhead; manually managing build artifacts is tedious. GoGo compiles tasks into native binaries for fast execution while handling all build management transparently.

### Who Benefits
- **Go Developers**: Run tasks as fast as native binaries with zero build management overhead. Repeated runs are instant thanks to caching.

### Business Value
Combines the convenience of scripting with the performance of compiled binaries. Smart caching means developers rarely wait for builds, and when they do, it's only because something actually changed.

---

## 2. User Journey

### Typical Workflow

```
1. Developer specifies a task name and any required arguments
2. GoGo locates the task workspace and parses task definitions
3. GoGo checks if a cached binary exists and is up to date
4. If cache is fresh, GoGo runs the cached binary immediately
5. If cache is stale or missing, GoGo compiles a new binary
6. The task runs with streaming output visible to the developer
7. GoGo reports success or displays the error message from the task
```

### Common Scenarios

| Scenario | Description | Outcome |
|----------|-------------|---------|
| First run | No cached binary exists yet | GoGo compiles and caches a binary, then runs it |
| Repeated run, no changes | Source files unchanged since last build | Cached binary runs instantly — no compilation |
| Run after editing a task | Source files are newer than the cached binary | GoGo recompiles automatically, then runs |
| Forced rebuild | Developer wants a clean build regardless of cache state | Cache is bypassed and a fresh binary is compiled |
| Pre-build for speed | Developer compiles ahead of time before a demo or CI run | Binary is ready; execution is instant |
| Optimized build | Developer wants a smaller binary to share or deploy | Symbols and debug info are stripped for a compact binary |
| Pass arguments by position | Developer provides values in order after the task name | Values are matched to parameters left-to-right |
| Pass arguments by name | Developer uses named flags | Values are matched to the corresponding parameter |

---

## 3. What Users Provide

### For Running a Task
- **Task name**: The name of the task to run (matches the function name)
- **Arguments** (if the task requires them): Values provided either positionally (in order) or as named flags

### For Build Control
- **Disable cache flag** (optional): Forces a fresh build regardless of cache state
- **Keep artifacts flag** (optional): Preserves intermediate generated files for debugging
- **Source directory** (optional): Explicitly specify where to find the task workspace
- **Output path** (optional): Specify where to save a pre-built binary
- **Optimize flag** (optional): Strip debug information for smaller binaries
- **Verbose flag** (optional): Show detailed build and execution output

---

## 4. What Users Receive

### Task Output
- Real-time streaming output from the running task
- The task's own printed output appears directly in the developer's terminal

### Confirmations
- Successful tasks exit cleanly with no extra messaging
- Failed tasks display the error message returned by the task function

### Build Artifacts (when requested)
- Pre-built binary at the specified output path
- Intermediate generated files preserved when keep-artifacts is enabled

---

## 5. Business Rules

| Rule | Description | Why It Exists |
|------|-------------|---------------|
| BR-001 | Cached binaries are reused when all source files are older than the binary | Avoids unnecessary compilation for fast repeated execution |
| BR-002 | Source file changes are detected by comparing file modification timestamps | Simple, reliable change detection without complex dependency tracking |
| BR-003 | Both task source files and dependency files are checked for changes | Ensures the binary reflects all relevant source changes |
| BR-004 | Cache is stored in a temporary directory keyed to the workspace | Each project gets its own cached binary without conflicts |
| BR-005 | Arguments can be passed positionally or by name | Flexibility — quick invocations use positional; complex ones use named |
| BR-006 | Type validation is applied to arguments before execution | Prevents runtime errors from mistyped arguments |

### Validation Rules
- Named task must exist in the workspace
- Argument types must match what the task expects (text, number, decimal, true/false)
- Required arguments must be provided (unless they have defaults)

### Timing Rules
- Compilation happens on-demand — only when the binary is missing or stale
- Pre-building creates the binary ahead of time, so execution is immediate later

---

## 6. Error Scenarios

| Situation | User Experience | Resolution |
|-----------|-----------------|------------|
| Task not found | Developer is told the task name doesn't match any known task | Check task name spelling and ensure the task is properly authored |
| Wrong argument type | Developer sees a validation error (e.g., "expected a number") | Provide the argument in the correct type |
| Missing required argument | Developer is told which argument is missing | Provide the missing argument |
| Compilation failure | Build error is displayed with details | Fix the error in the task source code |
| Task returns an error | Error message from the task is displayed | Address the issue described in the error message |

### Common Issues
- **Stale cache causing unexpected behavior**: Use the force-rebuild flag to ensure a fresh compilation.
- **Slow first run**: Expected — the first run compiles the binary. Subsequent runs use the cache and are much faster.

---

## 7. Related Domains

### Depends On
- **Workspace Setup**: Execution needs a workspace to find task source files
- **Task Authoring**: Execution compiles and runs authored tasks
- **Task Discovery**: Developers typically discover tasks before executing them

### Used By
- None — Task Execution is the terminal domain in the workflow

---

## 8. Access Control

Task execution is entirely local. Tasks run with the developer's own permissions and can do anything the developer could do manually (file operations, network calls, etc.). GoGo does not sandbox task execution.

---

## 9. Lifecycle

```
Task invoked → Cache checked → [Stale/Missing: Compile → Cache] → Execute → Report result
```

### State Definitions
- **Invoked**: Developer has specified a task to run
- **Cache checked**: GoGo compares source timestamps against the cached binary
- **Compile**: Source is fresh or cache is bypassed — GoGo generates and compiles a new binary
- **Cache**: New binary is stored for future reuse
- **Execute**: The binary runs with the provided arguments
- **Report result**: Success (clean exit) or failure (error message displayed)

### Transitions
- Invoked → Cache checked: Automatic on every run
- Cache checked → Execute: Cache is fresh — skip compilation
- Cache checked → Compile: Cache is stale or missing
- Compile → Cache: New binary saved
- Cache → Execute: Fresh binary runs
- Execute → Report result: Task completes or fails

---

## 10. Success Metrics

| Metric | Target | Why It Matters |
|--------|--------|----------------|
| Cached execution time | Near-instant (under 100ms overhead) | Developers shouldn't feel the tool's overhead |
| Cache hit rate on unchanged source | 100% | Unnecessary rebuilds waste developer time |
| Correct stale detection | 100% — always rebuild when source changes | Stale binaries cause subtle, confusing bugs |
| Argument validation before execution | All type mismatches caught pre-run | Fail fast with a clear message, not a cryptic runtime error |
