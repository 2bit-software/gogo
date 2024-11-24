# Comparing GoGo to Mage

> First, credit where credit is due: GoGo draws significant inspiration from Mage, and we're grateful for their pioneering work in Go-based task running. In fact, GoGo is fully compatible with existing Magefiles!

> Second, none of this is meant to diss or throw stones at Mage or any other tool. If something in here is wrong, please let me know and i'll fix it. This document was made very quickly and definitely needs a review.

## Key Differences

### 1. Enhanced CLI Experience
| Feature | Mage | GoGo                                                                    |
|---------|------|-------------------------------------------------------------------------|
| Shell Completion | Basic target completion | Full context-aware completion for targets, flags, and arguments         |
| Command Discovery | Lists targets | Rich function discovery with descriptions, argument types, and examples |
| Argument Handling | Positional only | Supports both positional, named, and environment flags                  |

### 2. Function Context
```go
// Mage
func Build() error {
    // Basic function
}

// GoGo
func Build(ctx gogo.Context) error {
    ctx.SetShortDescription("Build the project")
       .Example("gogo run Build --arch arm64")
       .Argument(arch).
       .Description("target architecture")
       .Default("amd64")
    // Enhanced function with self-documenting capabilities
}
```

### 3. Scope and Location
- **Mage**: Project-specific, requires magefile in current directory
- **GoGo**:
    - Walks up directory tree to find nearest `.gogo` folder
    - Supports global functions accessible from anywhere
    - Can be used both project-specific and system-wide

### 4. Build Capabilities
| Feature | Mage | GoGo |
|---------|------|------|
| Single Binary | ❌ | ✅ |
| Individual Function Binaries | ❌ | ✅ |
| Global Function Distribution | ❌ | ✅ |

### 5. Configuration
- **Mage**: Environment configuration options
- **GoGo**:
    - Flexible configuration using TOML or PKL
    - Validated configurations
    - Environment-specific settings

### 6. Function Arguments
```bash
# Mage
mage deployToEnv prod us-west-2 high-mem

# GoGo
gogo run deployToEnv --env prod --region us-west-2 --instance-type high-mem
```

### 7. Gradual Adoption Path
GoGo is designed to make the transition from shell scripts to Go as smooth as possible:
1. Start with simple Go functions (similar to Mage)
2. Gradually adopt enhanced features (context, completions, etc.)
3. Build standalone binaries when needed
4. Integrate with CI/CD systems

## When to Use What?

### Choose Mage if:
- You need a simple, proven task runner
- You prefer minimal configuration
- You're working on a single project
- You don't need advanced CLI features

### Choose GoGo if:
- You want rich CLI interactions and completions
- You need global function accessibility
- You're building complex automation tools
- You want to distribute standalone binaries
- You need advanced configuration options
- You're replacing complex shell scripts
