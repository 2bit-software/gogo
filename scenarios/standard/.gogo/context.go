package main

import (
	"fmt"
	"github.com/2bit-software/gogo/pkg/gogo"
)

func ContextWithNoUsage(ctx gogo.Context) {
	fmt.Println("ContextWithNoUsage")
}

func ShortDescriptionFunc(ctx gogo.Context) error {
	ctx.ShortDescription("this is a short description set specifically for the BasicShortDescription function")
	return nil
}

func ExampleFunc(ctx gogo.Context) error {
	ctx.Example("example")
	return nil
}

func ArgumentNameFunc(ctx gogo.Context, var1 string) error {
	ctx.Argument(var1).Name("personsName")
	return nil
}

func ArgumentShortFunc(ctx gogo.Context, var1 string) error {
	ctx.Argument(var1).Short('p')
	return nil
}

func ArgumentDefaultFunc(ctx gogo.Context, var1 string) error {
	ctx.Argument(var1).Default("default-value")
	fmt.Println(var1)
	return nil
}

func ArgumentOptionalFunc(ctx gogo.Context, var1 string) error {
	ctx.Argument(var1).Required()
	if var1 == "" {
		fmt.Println("var1 is empty")
		return nil
	}
	fmt.Println(var1)
	return nil
}

func ArgumentHelpFunc(ctx gogo.Context, var1 string) error {
	ctx.Argument(var1).Help("help text")
	return nil
}

func ArgumentAllowedValuesFunc(ctx gogo.Context, var1 int) error {
	ctx.Argument(var1).AllowedValues(8, 9, 10)
	return nil
}

func ArgumentRestrictedValuesFunc(ctx gogo.Context, var1 int) error {
	ctx.Argument(var1).RestrictedValues(1, 2, 3)
	return nil
}

func ArgumentDescriptionFunc(ctx gogo.Context, var1 string) error {
	ctx.Argument(var1).Description("this is the var 1 description")
	return nil
}

// BasicShortDescription is a function that uses the ShortDescription method to set the short description.
func BasicShortDescription(ctx gogo.Context) error {
	ctx.ShortDescription("this is a short description set specifically for the BasicShortDescription function")
	return nil
}

// BasicArgument is the builder argument that signifies the following methods
// are chained to the argument. By itself, it does nothing.
func BasicArgument(ctx gogo.Context, var1 string, var2 bool) error {
	ctx.Argument(var1)
	return nil
}

// BasicDescriptionArgument sets the description of the argument. This will show up in
// --help of the function.
func BasicDescriptionArgument(ctx gogo.Context, var1 string, var2 bool) error {
	ctx.Argument(var1).Description("describe what this argument does")
	return nil
}

func BasicCtxChained(ctx gogo.Context, var1 string, var2 bool) error {
	ctx.ShortDescription("set a description, this can use any go code to set the value")

	fmt.Println(var1, var2)

	return nil
}

func BasicArgumentChained(ctx gogo.Context, var1 string, var2 bool) error {
	ctx.Argument(var1).
		Description("describe what this argument does").
		AllowedValues("1", "2", "3").
		RestrictedValues("4", "5", "6").
		Default("1")

	fmt.Println(var1, var2)

	return nil
}
