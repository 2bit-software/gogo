# GoGo 🏃

## Why GoGo?
- 🚀 Run Go functions as CLI commands
- 🌳 Smart directory traversal
- 🔍 Enhanced shell discoverability and completion
- 🌍 Global function support
- 🎯 CI/CD/Automation friendly

## Quick Start
### Install
```bash
# go
go install github.com/2bit-software/gogo/cmd/gogo@latest
or 
go install github.com/2bit-software/gogo/cmd/wizard@latest
```

### Usage

```bash
# Init gogo within a repo
gogo init --local

# Or init the global function cache
# TODO: global cache not yet finished
gogo init --global

# Then create a Go function in that folder
echo 'func Hello() { fmt.Println("Hello, World!") }' > hello.go

# Run it
gogo gadget Hello
```

## Documentation
- [Getting Started](./getting-started.md)
- [Advanced Usage](./advanced-usage.md)
- [Comparison to Mage](./mage-comparison.md)
- [Configuration](./configuration.md) (TODO)
- [CI/CD Integration](./cicd.md) (TODO)
- [Misc notes](./notes.md)