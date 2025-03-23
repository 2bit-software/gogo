# Getting Started with GoGo

## Installation
[Installation instructions moved to main README]

## Basic Usage

### Local Functions (Project-Specific)
```bash
# Create a new GoGo project
mkdir -p myproject
gogo init --local

# Create your first function
cat << EOF > hello.go
package main

func Hello(name string) {
    fmt.Printf("Hello, %s!\n", name)
}
EOF

# Run it
gogo run Hello "World"
```

### Global Functions (System-Wide)
TO BE DECIDED: how is this configured? is it loaded from a .config?
Or should it assume whatever directory you run `gogo init --global` from?
```bash
# Create global GoGo directory
mkdir -p ~/.gogo
gogo init --global

# Create a utility function
cat << EOF > format.go
package main

func FormatCode() error {
    return exec.Command("go", "fmt", "./...").Run()
}
EOF

# Run it from anywhere
gogo gadget FormatCode
```

## Discovering Functions

### Listing Available Functions
```bash
# List all functions
gogo

# List with descriptions
gogo list --verbose

# List global functions only
gogo list --global

# List local functions only
gogo list --local
```

### Function Location Priority
1. Local `.gogo` directory in current path
2. Parent directories' `.gogo` folders (walks up)
3. Global `.gogo` directory

## Running Functions

### Different Ways to Run
```bash
# Basic execution
gogo gadget FunctionName

# With positional arguments
gogo gadget Greet "John" true 42

# With flags (recommended)
gogo gadget Greet --name "John" --verbose --count 42

# Running global function
gogo gadget g:FormatCode

```

## Enhanced Function Definitions

### Basic Function
```go
// Simple function
func SimpleGreet(name string) {
    fmt.Printf("Hello, %s!\n", name)
}
```

## Function Signatures
GoGo supports various function signatures:
```go
// Basic functions
func Basic()
func WithError() error
func WithArgs(name string, count int)
func WithArgsAndError(name string) error

// Context-aware functions
func WithContext(ctx gogo.Context)
func WithContextAndError(ctx gogo.Context) error
func WithContextAndArgs(ctx gogo.Context, name string)
func WithContextArgsAndError(ctx gogo.Context, name string) error
```

## Directory Structure
```
├── .gogo/                  # Local GoGo directory
│   ├── go.mod
│   ├── go.sum
│   ├── hello.go           # Function definitions
│   └── build.go           # More functions
│
~/.gogo/                   # Global GoGo directory
    ├── go.mod
    ├── go.sum
    └── utils.go           # Global utility functions
```

## Next Steps
- Learn about [advanced usage](./advanced-usage.md)
- Explore [configuration options](./configuration.md)
- Set up [CI/CD integration](./cicd.md)
- Check out [comparison with Mage](./comparison-to-mage.md)

## Common Issues and Solutions

### Function Not Found
- Ensure you're in a directory at or below your `.gogo` folder
- Check if the function name is exported (starts with capital letter)
- Verify the function is in a `.go` file inside a `.gogo` directory

### Build Errors
- Ensure your `go.mod` is properly initialized
- Check for missing dependencies
- Verify function signature matches supported formats
