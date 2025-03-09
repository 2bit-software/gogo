//go:build broken

// this file is inherently broken, so we ignore it
// but still used it for testing

package _gogo

import "fmt"

func BasicDescription(ctx templates.Context) error {
	ctx.SetShortDescription("set a description")
	return nil
}

func BasicArgument(ctx templates.Context, var1 string, var2 bool) error {
	ctx.Argument(var1)
	return nil
}

func BasicDescriptionArgument(ctx templates.Context, var1 string, var2 bool) error {
	ctx.Argument(var1).Description("describe what this argument does")
	return nil
}

func BasicCtxChained(ctx templates.Context, var1 string, var2 bool) error {
	ctx.SetShortDescription("set a description, this can use any go code to set the value")

	fmt.Println(var1, var2)

	return nil
}

func BasicArgumentChained(ctx templates.Context, var1 string, var2 bool) error {
	ctx.Argument(var1).
		Description("describe what this argument does").
		AllowedValues("1", "2", "3").
		RestrictedValues("4", "5", "6").
		Default("1")

	fmt.Println(var1, var2)

	return nil
}
