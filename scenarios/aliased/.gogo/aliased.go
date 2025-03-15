// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.
package main

import (
	"fmt"

	goCtx "github.com/2bit-software/gogo/pkg/gogo"
)

func AliasedCtxDescription(ctx goCtx.Context) error {
	ctx.ShortDescription("set a description")
	fmt.Println("AliasedCtxDescription")
	return nil
}

func AliasedCtxArgument(ctx goCtx.Context, var1 string, var2 bool) error {
	ctx.Argument(var1).Name("color")
	fmt.Printf("AliasedCtxArgument with var1: %v and var2: %v\n", var1, var2)
	return nil
}

func AliasedCtxDescriptionArgument(ctx goCtx.Context, var1 string, var2 bool) error {
	ctx.Argument(var1).Description("describe what this argument does")
	fmt.Printf("AliasedDescriptionArgument with var1: %v\n", var1)
	return nil
}

func AliasedCtxChained(ctx goCtx.Context, var1 string, var2 bool) error {
	ctx.ShortDescription("set a description, this can use any go code to set the value")
	fmt.Printf("AliasedCtxChained with var1: %v, var2: %v\n", var1, var2)
	return nil
}

func AliasedCtxArgumentChained(ctx goCtx.Context, var1 string, var2 bool) error {
	ctx.Argument(var1).
		Description("describe what this argument does").
		AllowedValues("1", "2", "3").
		RestrictedValues("4", "5", "6").
		Default("1")
	fmt.Printf("AliasedArgumentChained with var1: %v, var2: %v\n", var1, var2)
	return nil
}
