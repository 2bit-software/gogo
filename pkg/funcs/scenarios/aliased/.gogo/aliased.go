// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.
package _gogo

import (
	"fmt"

	goCtx "github.com/2bit-software/gogo/pkg/funcs/gogo"
)

func AliasedDescription(ctx goCtx.Context) error {
	ctx.ShortDescription("set a description")
	return nil
}

func AliasedArgument(ctx goCtx.Context, var1 string, var2 bool) error {
	ctx.Argument(var1)
	return nil
}

func AliasedDescriptionArgument(ctx goCtx.Context, var1 string, var2 bool) error {
	ctx.Argument(var1).Description("describe what this argument does")
	return nil
}

func AliasedCtxChained(ctx goCtx.Context, var1 string, var2 bool) error {
	ctx.ShortDescription("set a description, this can use any go code to set the value")
	fmt.Println(var1, var2)
	return nil
}

func AliasedArgumentChained(ctx goCtx.Context, var1 string, var2 bool) error {
	ctx.Argument(var1).
		Description("describe what this argument does").
		AllowedValues("1", "2", "3").
		RestrictedValues("4", "5", "6").
		Default("1")

	fmt.Println(var1, var2)
	return nil
}
