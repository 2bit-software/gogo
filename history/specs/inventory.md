# GoGo - Capability Inventory

This document provides a comprehensive list of all capabilities in GoGo.

## Summary Statistics

- **Total Capabilities**: 18
- **Breakdown by Domain**:
  - Workspace Setup: 3 capabilities
  - Task Authoring: 6 capabilities
  - Task Discovery: 4 capabilities
  - Task Execution: 5 capabilities

---

## DOMAIN: Workspace Setup

**Purpose**: Enables developers to quickly scaffold a new task workspace in any project.

| # | Capability | Description | Who Can Use |
|---|---|---|---|
| 1 | Initialize a workspace | Create a new task folder with starter configuration in the current project | Go Developer |
| 2 | Generate example tasks | Scaffold example task files that demonstrate common patterns | Go Developer |
| 3 | Auto-configure dependencies | Automatically set up the workspace so tasks can reference the GoGo toolkit | Go Developer |

---

## DOMAIN: Task Authoring

**Purpose**: Enables developers to write plain Go functions that GoGo recognizes and converts into runnable CLI commands.

### Function Conventions

| # | Capability | Description | Who Can Use |
|---|---|---|---|
| 1 | Define a task as a function | Write an exported function and GoGo automatically treats it as a runnable task | Go Developer |
| 2 | Add typed arguments | Define function parameters using text, number, decimal, or true/false types and GoGo generates typed CLI arguments | Go Developer |
| 3 | Document tasks with comments | Write a comment above a function and GoGo uses it as the task's help text and description | Go Developer |
| 4 | Signal success or failure | Return an error from a function to indicate task failure, or return nothing to indicate success | Go Developer |

### Advanced Authoring

| # | Capability | Description | Who Can Use |
|---|---|---|---|
| 5 | Access task context | Use a special first parameter to access task metadata, configure argument details, and set allowed values | Go Developer |
| 6 | Organize tasks across files | Split tasks across multiple files in the workspace folder for better organization | Go Developer |

---

## DOMAIN: Task Discovery

**Purpose**: Enables developers to browse and understand all available tasks without reading source code.

| # | Capability | Description | Who Can Use |
|---|---|---|---|
| 1 | List all available tasks | View a formatted list of all tasks with their short descriptions | Go Developer |
| 2 | View task arguments | See what arguments a task expects, including types and defaults | Go Developer |
| 3 | Discover tasks from subdirectories | Run GoGo from any subdirectory and it finds the nearest task workspace by walking up the directory tree | Go Developer |
| 4 | View version and build details | Check which version of GoGo is installed and see detailed build information | Go Developer |

---

## DOMAIN: Task Execution

**Purpose**: Enables developers to run tasks with arguments, leveraging smart caching for fast repeated execution.

| # | Capability | Description | Who Can Use |
|---|---|---|---|
| 1 | Run a task with arguments | Execute any task by name, passing arguments either by position or by name | Go Developer |
| 2 | Benefit from smart caching | GoGo tracks source file changes and reuses compiled binaries when nothing has changed | Go Developer |
| 3 | Force a fresh rebuild | Bypass the cache and force GoGo to recompile tasks from scratch | Go Developer |
| 4 | Pre-build task binaries | Compile tasks ahead of time so they execute instantly when needed | Go Developer |
| 5 | Build optimized binaries | Create smaller, optimized binaries suitable for distribution or deployment | Go Developer |

---

## Domain Groupings for Business Specification

The capabilities can be logically grouped into these business domains:

1. **Workspace Setup** - Getting started with GoGo in a new or existing project
2. **Task Authoring** - Writing and documenting tasks as Go functions
3. **Task Discovery** - Browsing, searching, and understanding available tasks
4. **Task Execution** - Running tasks, caching builds, and optimizing performance
