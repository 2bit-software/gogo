# Advanced Usage

### Using GoGo Context
Explain why GoGo's `gogo.Context` is useful for enhancing function definitions.

Example:
```go
func AdvancedGreet(ctx gogo.Context, name string, count int) error {
    // Define function metadata
    ctx.SetShortDescription("Greet someone multiple times")
       .Example("gogo run AdvancedGreet --name 'John' --count 3")
    
    // Configure arguments
    ctx.Argument(name).
        Description("Name to greet").
        Default("World").
        Required(true)

    ctx.Argument(count).
        Description("Number of greetings").
        Default(1).
        RestrictedValues(1, 2, 3, 4, 5)

    // Function logic
    for i := 0; i < count; i++ {
        fmt.Printf("Hello, %s! (%d/%d)\n", name, i+1, count)
    }
    return nil
}
```

### GoGo Context Methods and their Usage
TODO: This

### Single binary per function
TODO: This

### Global Function Distribution
Single binary for global functions. It's useful for distribution.