# GoGo: A Flexible Function Execution Tool

GoGo is a versatile tool designed to simplify the execution of Go functions across different contexts. It offers a unique approach to function management, execution, and discovery.

## Table of Contents

1. [General Process](#general-process)
2. [Goals](#goals)
3. [Unique Features](#unique-features)
4. [Build Modes](#build-modes)
5. [Autocomplete](#autocomplete)
6. [Binary Output](#binary-output)
7. [Arguments and Flags](#arguments-and-flags)
8. [Usage Scenarios](#usage-scenarios)
9. [Configuration](#configuration)
10. [Future Improvements](#future-improvements)
11. [Credits](#credits)

## General Process

1. Create a local or global .gogo folder
2. Write an exported Go function in a .go file
3. Use `gogo run FunctionName` to run it from the cli

## Magefiles
This takes heavy influence from magefiles. In fact, since both tools just run go functions, GoGo supports all magefiles.
While mage is a great tool, it has some limitations that GoGo aims to solve. Much love for mage, but personally I needed more. 

## Justification

GoGo aims to improve upon mage in the following ways:
1. Better CLI output and shell completions.
2. Better DX while inside repositories with GoGo/mage functions.
3. Also provide a way to run global functions from anywhere on the system.
4. Have the ability to build a single binary for all functions in a directory.

## Unique Features

### GoGo Context in function signatures
An opt-in feature of GoGo functions is using a `gogo.Context` as the first argument. This allows for more advanced features, like providing shell completions for arguments, including defaults, allowed, disallowed, required, and optional values.
Here's an example:
```go
func AdvancedFunction(ctx gogo.Context, name string, include bool, value int) error {
	ctx.
		SetShortDescription("set a short description used when listing functions").
		Example("example").
		Argument(name).
		Description("this is the name").
		Default("default-value").
		Argument(include).
		Description("this is the include bool").
		Default(true).
		Argument(value).
		RestrictedValues(1, 2, 3).
		AllowedValues(5, 6, 7).
		Description("this is the value").
		Default(3)

	fmt.Printf("name: %s\n", name)
	if include {
		fmt.Printf("value: %d\n", value)
	}

	return nil
}
```

### Allow functions to be run with flags instead of positional arguments.
Positional arguments provide no details about their usage when reading the command. Flags provide a way to provide more information about the arguments.
Imagine the following
```bash
gogo run CompileGoToFolderWithTags cmd/api/mocked/ . mocks,walternate mocked_api linux arm64
```

This is hard to read and understand. Instead, we can use flags.
```bash
gogo run CompileGoToFolderWithTags --source cmd/api/mocked/ --output . --tags mocks,walternate --package mocked_api --os linux --arch arm64
```
This is easier to read and understand the intent after the fact.

### Walk up the tree
GoGo will walk up the tree to find the nearest .gogo folder. This allows for a single .gogo folder to be used in a repository. This is especially useful for monorepos.

### Enhanced shell-completion
GoGo provides enhanced shell completion for functions. This includes providing completions for arguments, flags, and functions. The functions themselves continue to self-identify the function signature, and then you can opt-in to more enhanced features by utilizing the gogo.Context features.

### Gradual adoption of features
GoGo is meant to facilitate replacing BASH scripts with a type safe and more powerful language. However, one reason many people use BASH is how easy it is to get started. GoGo optimizes for this, requiring only a few small steps to setup (a folder and a go.mod), and then you can write simple scripts to get work done.
As you need more advanced features, like self-documenting functions, enhanced autocomplete, and single-function binaries, you can opt-in to these features.

### Global functions
It's possible to have global functions that can be run from anywhere on the system. This is useful for things like `gogo run FormatCode` or `gogo run BuildAll`. These functions can be run from anywhere on the system, and are not tied to a specific repository. 

### Configuration
GoGo is entirely data-driven. This means all the actual values are defined in a configuration. This is normally done in TOML, which GoGo supports. However, GoGo also supports PKL, that way we can provide a way to validate your configuration, and hopefully increase the obviousness of the config.

## Build Modes

GoGo supports various build modes:

- Default local binary generation
- Individual function binary generation
- Global function binary generation
- Optimized build mode, where all normal debug symbols are stripped

## Autocomplete

GoGo implements a unique autocomplete system that:

1. Lists all possible targets when no arguments are given
2. Provides target-specific completions when a target is specified

## Arguments and Flags

GoGo uses a specific syntax for passing arguments and flags:

```
gogo function -- --flag value
```

## Configuration

GoGo uses configuration files to customize its behavior. These can be in TOML, YAML, or PKL format.

## How it works
1. GoGo searches for nearby .gogo folders. Nearby means either in the current directory, or as a sibling in one of the parent directories.
2. Parse all .go files that are in the .gogo folder for functions that can be run.
3. Generate the binary with all functions
4. Find the function that matches the requested function
5. Parse arguments from the command line or input into the desired types
6. Pass the arguments to the built binary/function

## Future Improvements

- Implement GOOS and GOARCH support
- Add dependency management
- Allow setting build tags
- Improve variable naming and collision prevention
- Provide options for preserving generated files
- Support functions with error returns

## Credits

ASCII image generation done by https://www.asciiart.eu/image-to-ascii
