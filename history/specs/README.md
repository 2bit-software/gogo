# GoGo - Business Specification

## Overview

GoGo is a task runner that lets Go developers write plain functions and instantly run them as command-line commands. Instead of building CLI boilerplate, argument parsing, and help text by hand, developers write ordinary Go functions and GoGo handles the rest — generating a fully-featured command-line interface automatically.

GoGo solves the problem of repetitive CLI scaffolding for project automation. Whether it's build scripts, deployment tasks, or developer tooling, GoGo turns function signatures into typed, documented commands with zero ceremony.

**Business Domains**: 4

---

## Executive Summary

GoGo enables Go developers to create and run project automation tasks instantly from function definitions. The system handles:

1. **Workspace Setup** - Scaffolding new task workspaces so developers can start writing tasks immediately
2. **Task Authoring** - Converting plain Go functions into runnable commands with automatic argument handling
3. **Task Discovery** - Browsing and understanding all available tasks in a project
4. **Task Execution** - Running tasks with arguments, smart caching, and build optimization

---

## Business Domains

### [01. Workspace Setup](./01-workspace-setup.md)
Developers can initialize a new GoGo workspace in any project with a single command. GoGo creates the necessary folder structure and example tasks so developers can start writing automation immediately.
- Initialize a new task workspace
- Generate example tasks as starting templates
- Configure workspace dependencies automatically

### [02. Task Authoring](./02-task-authoring.md)
Developers write standard Go functions following simple conventions, and GoGo automatically recognizes them as runnable tasks. Function names become command names, parameters become typed arguments, and comments become help text.
- Write functions that automatically become CLI commands
- Define typed arguments (text, numbers, true/false values)
- Add descriptions and documentation through code comments
- Access task context for richer interactions

### [03. Task Discovery](./03-task-discovery.md)
Developers can browse all available tasks in a project, see their descriptions and expected arguments, and understand what each task does — all from the command line.
- List all available tasks with descriptions
- View argument requirements for any task
- See tasks from anywhere within a project directory tree

### [04. Task Execution](./04-task-execution.md)
Developers run tasks by name, passing arguments either positionally or by name. GoGo compiles tasks into fast binaries, caches them intelligently, and only rebuilds when source files change.
- Run tasks with positional or named arguments
- Automatic smart caching avoids unnecessary rebuilds
- Force rebuild when needed
- Pre-build binaries for faster execution
- Optimize builds for distribution

---

## Core Concepts

### Task
A runnable unit of automation. Each task corresponds to a single Go function that GoGo can execute from the command line. Tasks have a name, optional description, and zero or more arguments.

### Workspace
A designated folder within a project (typically `.gogo`) that contains task definitions. GoGo searches for this folder starting from the current directory and walking up the directory tree, similar to how version control tools find their configuration.

### Argument
A value that a task accepts from the command line. Arguments are strongly typed — they can be text, whole numbers, decimal numbers, or true/false values. GoGo automatically validates that provided values match the expected type.

### Task Binary
A compiled, cached version of the task runner. GoGo generates and compiles task definitions into a single binary that executes quickly. The binary is cached and reused until task definitions change.

---

## Key Workflows

### Starting a New Project with Tasks
```
1. Developer runs the initialization command in their project
2. GoGo creates a task workspace folder with example tasks
3. GoGo sets up necessary dependencies automatically
4. Developer edits the example tasks or writes new ones
5. Developer runs tasks immediately — no additional setup needed
```

### Writing and Running a New Task
```
1. Developer creates a new function in the task workspace
2. Developer gives the function a descriptive name and adds a comment
3. Developer adds typed parameters for any inputs needed
4. Developer runs the task by name from anywhere in the project
5. GoGo compiles the task, caches the binary, and executes it
```

### Running an Existing Task
```
1. Developer lists available tasks to find the right one
2. Developer sees task descriptions and argument requirements
3. Developer runs the task, passing arguments by position or name
4. GoGo checks if a cached build exists and is up to date
5. If cache is fresh, executes immediately; otherwise rebuilds first
```

---

## User Types

### Go Developer
- **Description**: A software developer working in a Go project who needs to automate repetitive tasks like building, testing, deploying, or managing infrastructure.
- **Primary Activities**: Writing task functions, running tasks from the command line, browsing available tasks.
- **Access Level**: Full access to all GoGo capabilities. All operations are local to the developer's machine.

### Mage User (Migration Path)
- **Description**: A developer currently using the Mage task runner who wants to switch to GoGo. GoGo supports Mage-compatible workspace folders for easier migration.
- **Primary Activities**: Running existing tasks through GoGo, gradually adopting GoGo-specific features.
- **Access Level**: Full access. GoGo recognizes Mage workspace folders alongside its own.

---

## Access Model

GoGo is a local developer tool with no multi-user access model. All operations run on the developer's own machine with their own file system permissions. There is no authentication, authorization, or data separation — the tool operates entirely within the developer's local project context.

---

## Files in This Specification

```
specs/
├── README.md                    # This file
├── inventory.md                 # Complete capability inventory
├── 01-workspace-setup.md        # Workspace initialization and scaffolding
├── 02-task-authoring.md         # Writing tasks as Go functions
├── 03-task-discovery.md         # Browsing and understanding tasks
└── 04-task-execution.md         # Running tasks and build management
```

Note: All files are in a flat structure (no subfolders) for compatibility with systems that don't support nested directories.

---

## Version Information

- **Generated**: 2026-03-23
- **Source**: GoGo v0.1.0 codebase audit
