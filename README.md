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
3. Use `gogo FunctionName` to run it from the cli

## Goals

GoGo aims to work in the following ways:

1. Recursively search for GoGo-related folders and files
2. Search the global GoGo namespace if no local function is found
3. Provide a way to force using the global version

## Unique Features

Compared to similar tools like Magefiles and Bake, GoGo offers:

- Better CLI output and shell completions using Cobra/Viper
- Improved argument parsing and constraints
- Flexible configuration file formats (TOML, YAML, or PKL)
- Global namespace for system-wide function replacement
- Automatic recompilation of global namespace functions

## Build Modes

GoGo supports various build modes:

- Default local binary generation
- Individual function binary generation
- Global function binary generation

## Autocomplete

GoGo implements a unique autocomplete system that:

1. Lists all possible targets when no arguments are given
2. Provides target-specific completions when a target is specified

## Binary Output

GoGo can generate binaries in two modes:

1. A single binary for all functions in a directory
2. Individual binaries for each function (useful for global functions)

## Arguments and Flags

GoGo uses a specific syntax for passing arguments and flags:

```
gogo function -- --flag value
```

## Usage Scenarios

- Listing global functions
- Executing functions from within subdirectories
- Working with unique Go module configurations

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
