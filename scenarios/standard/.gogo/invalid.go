package main

import (
	"fmt"
	"github.com/2bit-software/gogo/pkg/gogo"
)

// the below are not valid GoGo functions

// not valid because the second parameter is a gogo.Context, not the first
func GoGoCtxInWrongPos(var1 string, ctx gogo.Context, var2 bool) error {
	return nil
}

// not valid because we can only return an error or nothing
func WrongReturnType(ctx gogo.Context, var1 string, var2 bool) string {
	return ""
}

type MyType struct{}

// WrongCompoundType Compound types which should be ignored
func WrongCompoundType(ctx gogo.Context, var1 MyType) error {
	return nil
}

// InvalidExportedWithPointerArgs should not be valid
// This should not show up.
func InvalidExportedWithPointerArgs(ctx gogo.Context, var1 *string, var2 *bool) error {
	return nil
}

// the following should not be valid because they are not exported functions,
// regardless of their signature

func basicArgument(ctx gogo.Context, var1 string, var2 bool) error {
	ctx.Argument(var1)
	return nil
}

func basicDescriptionArgument(ctx gogo.Context, var1 string, var2 bool) error {
	ctx.Argument(var1).Description("describe what this argument does")
	return nil
}

func basicCtxChained(ctx gogo.Context, var1 string, var2 bool) error {
	ctx.ShortDescription("set a description, this can use any go code to set the value")

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
