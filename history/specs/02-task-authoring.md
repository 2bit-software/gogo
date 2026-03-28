# Task Authoring - Business Specification

## Overview

Task Authoring is the core creative activity in GoGo. Developers write standard Go functions and GoGo automatically converts them into fully-featured CLI commands. Function names become command names, parameters become typed arguments, and comments become help text — no configuration files, decorators, or registration code needed.

## Capabilities Summary

| Capability | Description | Priority |
|------------|-------------|----------|
| Define a task as a function | Write an exported function that becomes a runnable command | Core |
| Add typed arguments | Define parameters that become validated CLI arguments | Core |
| Document tasks with comments | Add descriptions and help text through code comments | Core |
| Signal success or failure | Indicate task outcomes through return values | Core |
| Access task context | Use enriched context for argument metadata and constraints | Supporting |
| Organize tasks across files | Split tasks into multiple files within the workspace | Supporting |

---

## 1. Business Purpose

### Problem Solved
Building CLI tools for project automation typically requires significant boilerplate: argument parsing, type validation, help text generation, and command registration. GoGo eliminates all of this by deriving the CLI interface directly from function signatures and comments.

### Who Benefits
- **Go Developers**: Write tasks as simple functions without learning a framework or writing boilerplate. The mental model is "write a function, run it as a command."

### Business Value
Dramatically lowers the cost of creating project automation. Tasks that would require dozens of lines of CLI setup code become single functions. This encourages teams to automate more and maintain their automation because the overhead is nearly zero.

---

## 2. User Journey

### Typical Workflow

```
1. Developer creates a new file in the task workspace folder
2. Developer writes an exported function with a descriptive name
3. Developer adds parameters for any inputs the task needs
4. Developer adds a comment above the function describing what it does
5. The task is immediately available to run — no registration or configuration needed
```

### Common Scenarios

| Scenario | Description | Outcome |
|----------|-------------|---------|
| Simple task, no arguments | Developer writes a function with no parameters | Task runs with just its name, no arguments needed |
| Task with typed inputs | Developer adds text, number, or boolean parameters | Each parameter becomes a named, typed CLI argument |
| Task that can fail | Developer returns an error when something goes wrong | GoGo reports the failure to the user with the error message |
| Self-documenting task | Developer adds a comment above the function | Comment appears as description when browsing tasks |
| Constrained arguments | Developer uses task context to set allowed values for an argument | GoGo validates input against the allowed values |

---

## 3. What Users Provide

### For Defining a Task
- **Function name**: Must start with an uppercase letter (Go convention for exported functions). The name becomes the command name.
- **Parameters** (optional): Each parameter becomes a CLI argument. Supported types are text, whole numbers, decimal numbers, and true/false values.
- **Comment** (optional): A comment above the function. The first line becomes the short description; additional lines become extended help text.

### For Using Task Context
- **Task context parameter**: Must be the first parameter if used. Provides access to argument metadata and configuration.
- **Argument constraints** (optional): Developers can specify allowed values, default values, short flags, and help text for each argument through the context.

---

## 4. What Users Receive

### Confirmations
- When a task is authored correctly, it appears in the task listing immediately
- When a task runs successfully, it exits cleanly
- When a task fails, the error message is displayed clearly

### Information
- Developers see their function's comment as the task description
- Developers see their parameters as documented, typed arguments in help output

---

## 5. Business Rules

| Rule | Description | Why It Exists |
|------|-------------|---------------|
| BR-001 | Only exported functions become tasks | Convention separates public tasks from private helper functions |
| BR-002 | Arguments must be simple types (text, numbers, true/false) | CLI arguments are inherently scalar — complex types cannot be passed on a command line |
| BR-003 | Functions can only return nothing or an error | Tasks are actions, not data producers — they succeed or fail |
| BR-004 | Task context must be the first parameter if used | Consistent convention makes parsing reliable |
| BR-005 | Test files are ignored | Tasks in test files are for testing GoGo itself, not meant to be run as commands |
| BR-006 | The first comment line is the short description | Provides concise listing text while allowing extended documentation |

### Validation Rules
- Function names must start with an uppercase letter
- Parameters cannot be pointer types or complex structures
- Only one error return value is allowed (no multiple returns)
- Task context, if used, must be the first parameter

---

## 6. Error Scenarios

| Situation | User Experience | Resolution |
|-----------|-----------------|------------|
| Function uses unsupported parameter types | Function is silently skipped — it won't appear as a task | Change parameters to supported simple types |
| Function returns multiple values | Function is skipped | Change to return only an error or nothing |
| Syntax error in task file | Build fails with a clear error message | Fix the syntax error in the source file |
| Duplicate function names across files | Build fails due to naming conflict | Rename one of the conflicting functions |

### Common Issues
- **Task not appearing in listing**: The function is likely unexported (starts with lowercase) or uses unsupported parameter types. Check that the function name starts with an uppercase letter and all parameters are simple types.
- **Unexpected argument behavior**: When using task context, ensure it is the first parameter and that argument names match the function parameters.

---

## 7. Related Domains

### Depends On
- **Workspace Setup**: Tasks are written in the workspace folder that Setup creates

### Used By
- **Task Discovery**: Discovery reads authored tasks to present listings and help text
- **Task Execution**: Execution compiles and runs the authored tasks

---

## 8. Access Control

All task authoring happens locally on the developer's machine. Any developer with file system access to the project can create, modify, or delete task files.

---

## 9. Lifecycle

```
Function written → Detected by GoGo → Available as command → Modified → Re-detected on next run
```

### State Definitions
- **Written**: Developer has saved a function in the workspace
- **Detected**: GoGo has parsed the function and validated its signature
- **Available**: The task appears in listings and can be run
- **Modified**: Developer has changed the function; cached binary is now stale

### Transitions
- Written → Detected: Happens automatically when GoGo parses the workspace
- Detected → Available: Immediate — valid functions become available tasks
- Available → Modified: Developer edits the source file
- Modified → Detected: GoGo re-parses on next invocation

---

## 10. Success Metrics

| Metric | Target | Why It Matters |
|--------|--------|----------------|
| Lines of code per task vs. traditional CLI setup | 80%+ reduction | Core value proposition — simplicity |
| Time from writing a function to running it | Under 5 seconds (first run) | Instant feedback loop encourages authoring |
| Functions rejected due to convention violations | Clear, actionable feedback | Developers should understand why a function isn't recognized |
