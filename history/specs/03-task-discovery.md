# Task Discovery - Business Specification

## Overview

Task Discovery enables developers to browse, understand, and learn about all available tasks in a project without reading source code. GoGo presents a formatted, color-coded listing of tasks with descriptions and argument details, making it easy to find and understand the right task for the job.

## Capabilities Summary

| Capability | Description | Priority |
|------------|-------------|----------|
| List all available tasks | View formatted task listing with descriptions | Core |
| View task arguments | See expected arguments with types and defaults | Core |
| Discover from subdirectories | Find tasks from anywhere in the project tree | Core |
| View version and build info | Check installed GoGo version and build details | Supporting |

---

## 1. Business Purpose

### Problem Solved
In projects with many automation tasks, developers need a quick way to find out what tasks exist, what they do, and how to invoke them — without digging through source files. Task Discovery provides this at-a-glance understanding.

### Who Benefits
- **Go Developers**: Quickly find the right task for a job, understand what arguments it expects, and see how to use it — all from the command line.

### Business Value
Reduces the time spent searching for and understanding automation tasks. New team members can immediately see what automation is available in a project, improving onboarding and reducing tribal knowledge dependencies.

---

## 2. User Journey

### Typical Workflow

```
1. Developer runs GoGo without specifying a task
2. GoGo searches for the nearest task workspace (current directory and parent directories)
3. GoGo parses all task files in the workspace
4. Developer sees a formatted list of tasks with short descriptions
5. Developer identifies the task they need and sees its argument requirements
```

### Common Scenarios

| Scenario | Description | Outcome |
|----------|-------------|---------|
| Browse all tasks | Developer runs GoGo with no arguments | Formatted list of all tasks with short descriptions |
| Find tasks from a subdirectory | Developer is deep in the project tree | GoGo walks up the directory tree to find the workspace |
| Check version | Developer wants to know which GoGo version is installed | Version number and build details displayed |

---

## 3. What Users Provide

### For Listing Tasks
- Nothing required — just run GoGo without arguments
- **Source directory** (optional): Explicitly specify where to look for the task workspace

---

## 4. What Users Receive

### Information
- A color-coded, terminal-width-aware listing of all tasks
- Each task shows its name and short description, aligned in columns
- Task arguments are displayed with their names, types, and default values
- Text wraps intelligently based on terminal width

### Version Details
- Version number and commit reference (compact format)
- Detailed build information including build time, builder identity, and platform details

---

## 5. Business Rules

| Rule | Description | Why It Exists |
|------|-------------|---------------|
| BR-001 | GoGo searches for workspaces by walking up the directory tree | Developers should not need to be in the workspace folder to discover tasks |
| BR-002 | Search stops at the project root (version control boundary) | Prevents accidentally picking up tasks from unrelated parent projects |
| BR-003 | Three folder names are recognized: `.gogo`, `gogofiles`, `magefiles` | Supports multiple conventions and Mage migration |
| BR-004 | Task listing is formatted to fit the terminal width | Output remains readable regardless of terminal size |
| BR-005 | Only the first line of a function comment becomes the short description | Keeps listings concise; full help available per-task |

### Validation Rules
- If no task workspace is found in the directory tree, the developer is informed clearly

---

## 6. Error Scenarios

| Situation | User Experience | Resolution |
|-----------|-----------------|------------|
| No task workspace found | Developer sees a message indicating no tasks were found | Initialize a workspace or navigate to the correct project |
| Task files have syntax errors | Build fails with an error pointing to the problematic file | Fix the syntax error in the indicated file |
| Empty workspace | Task listing shows no tasks | Write some task functions in the workspace |

### Common Issues
- **Tasks not showing up**: The most common cause is functions that don't meet authoring conventions — unexported names or unsupported parameter types. See the Task Authoring specification.

---

## 7. Related Domains

### Depends On
- **Workspace Setup**: Discovery searches for the workspace folder that Setup creates
- **Task Authoring**: Discovery reads and presents information from authored tasks

### Used By
- **Task Execution**: Developers typically discover a task first, then execute it

---

## 8. Access Control

Task discovery is entirely local. Any developer who can access the project files can browse available tasks.

---

## 9. Lifecycle

Task Discovery is stateless — it reads the current workspace state each time it runs. There is no persistent discovery state.

```
Developer invokes GoGo → Workspace found → Tasks parsed → Listing displayed
```

---

## 10. Success Metrics

| Metric | Target | Why It Matters |
|--------|--------|----------------|
| Time to display task listing | Under 1 second (cached) | Quick feedback keeps developers in flow |
| Task descriptions visible without reading source | 100% of documented tasks | Whole point of discovery is avoiding source reading |
| Workspace found from any project subdirectory | Always, if workspace exists | Convenience and reliability |
