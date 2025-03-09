package _gogo

import (
	"fmt"
	"github.com/2bit-software/gogo/pkg/funcs/gogo"
)

// None of these functions should be parsed as GoGo functions. They can be *used in a function*,
// but they are not valid GoGo functions themselves.

func basicArgument(ctx gogo.Context, var1 string, var2 bool) error {
	ctx.Argument(var1)
	return nil
}

func basicDescriptionArgument(ctx gogo.Context, var1 string, var2 bool) error {
	ctx.Argument(var1).Description("describe what this argument does")
	return nil
}

func basicCtxChained(ctx gogo.Context, var1 string, var2 bool) error {
	ctx.SetShortDescription("set a description, this can use any go code to set the value")

	fmt.Println(var1, var2)

	return nil
}

func basicArgumentChained(ctx gogo.Context, var1 string, var2 bool) error {
	ctx.Argument(var1).
		Description("describe what this argument does").
		AllowedValues("1", "2", "3").
		RestrictedValues("4", "5", "6").
		Default("1")

	fmt.Println(var1, var2)

	return nil
}
