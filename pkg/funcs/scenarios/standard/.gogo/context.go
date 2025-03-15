package _gogo

import (
	"fmt"
	"github.com/2bit-software/gogo/pkg/funcs/gogo"
)

func ShortDescriptionFunc(ctx gogo.Context) error {
	ctx.ShortDescription("set a description")
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
	ctx.Argument(var1).Optional()
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
