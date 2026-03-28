# Workspace Setup - Business Specification

## Overview

Workspace Setup enables developers to create a task workspace in any Go project with a single command. It scaffolds the folder structure, generates starter examples, and configures dependencies so that developers can write and run tasks immediately — no manual configuration required.

## Capabilities Summary

| Capability | Description | Priority |
|------------|-------------|----------|
| Initialize a workspace | Create a task folder with starter files | Core |
| Generate example tasks | Provide working examples to learn from | Core |
| Auto-configure dependencies | Set up module references automatically | Core |

---

## 1. Business Purpose

### Problem Solved
Starting a new task automation setup in a Go project requires creating folders, writing boilerplate, and configuring module dependencies. This is tedious and error-prone, especially for developers new to the tool.

### Who Benefits
- **Go Developers**: Get a working task setup in seconds instead of manually creating files and configuring dependencies.

### Business Value
Eliminates the friction of getting started. Developers can go from zero to running their first custom task in under a minute, which encourages adoption and consistent use of task automation across projects.

---

## 2. User Journey

### Typical Workflow

```
1. Developer navigates to their Go project
2. Developer runs the initialization command
3. GoGo creates a task workspace folder (defaults to ".gogo")
4. GoGo generates example task files showing common patterns
5. GoGo configures module dependencies automatically
6. Developer can immediately run the example tasks or start writing their own
```

### Common Scenarios

| Scenario | Description | Outcome |
|----------|-------------|---------|
| New project setup | Developer wants to add task automation to a fresh Go project | Task workspace created with examples and dependencies ready |
| Custom folder name | Developer wants the task folder named something other than the default | Workspace created in the specified folder |
| Mage migration | Developer has existing Mage tasks and wants to try GoGo | GoGo recognizes Mage workspace folders, easing the transition |

---

## 3. What Users Provide

### For Workspace Initialization
- **Project location**: The developer must be inside a Go project (or module) directory
- **Folder name** (optional): A custom name for the task workspace folder. Defaults to `.gogo` if not specified.

---

## 4. What Users Receive

### Confirmations
- Workspace folder is created in the project directory
- Example task files are generated inside the workspace
- Module dependencies are configured and resolved

### Information
- The developer can see the created files and immediately understand the expected structure by reading the examples

---

## 5. Business Rules

| Rule | Description | Why It Exists |
|------|-------------|---------------|
| BR-001 | Workspace must be inside a Go module | Tasks are Go code and need a module context to compile |
| BR-002 | Default folder name is `.gogo` | Convention keeps task files hidden and consistent across projects |
| BR-003 | GoGo also recognizes `gogofiles` and `magefiles` folders | Supports migration from Mage and offers naming flexibility |
| BR-004 | Dependencies are configured automatically | Developers should not need to manually wire up GoGo references |

### Validation Rules
- The target folder must not already exist (prevents overwriting existing work)
- The project must be within a Go module (a `go.mod` file must be findable)

---

## 6. Error Scenarios

| Situation | User Experience | Resolution |
|-----------|-----------------|------------|
| No Go module found | Developer is told they need to be in a Go module | Run `go mod init` first, then retry |
| Folder already exists | Developer is informed the workspace folder already exists | Use the existing folder or choose a different name |

### Common Issues
- **Running outside a Go project**: GoGo needs a module context to work. Developer should initialize a Go module first.

---

## 7. Related Domains

### Depends On
- None — Workspace Setup is the entry point for all other domains

### Used By
- **Task Authoring**: Developers write tasks in the workspace that Setup creates
- **Task Discovery**: Discovery searches for the workspace folder that Setup creates
- **Task Execution**: Execution compiles and runs tasks from the workspace

---

## 8. Access Control

GoGo is a local developer tool. All operations use the developer's file system permissions. There is no multi-user access model.

---

## 9. Lifecycle

```
No workspace → Initialized → Ready for task authoring
```

### State Definitions
- **No workspace**: The project has no GoGo task folder yet
- **Initialized**: The workspace folder, example files, and dependencies have been created
- **Ready for task authoring**: The developer can start writing and running tasks

### Transitions
- No workspace → Initialized: Developer runs the initialization command
- Initialized → Ready: Happens automatically — initialization produces a fully working workspace

---

## 10. Success Metrics

| Metric | Target | Why It Matters |
|--------|--------|----------------|
| Time from command to first runnable task | Under 30 seconds | Fast setup encourages adoption |
| Setup failures requiring manual intervention | Zero for standard Go projects | Reliability builds trust in the tool |
