// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package _gogo

import (
	"fmt"

	"github.com/2bit-software/gogo/pkg/funcs/gogo"
)

func BasicDescription(ctx gogo.Context) error {
	ctx.SetShortDescription("set a description")
	return nil
}

func BasicArgument(ctx gogo.Context, var1 string, var2 bool) error {
	ctx.Argument(var1)
	return nil
}

func BasicDescriptionArgument(ctx gogo.Context, var1 string, var2 bool) error {
	ctx.Argument(var1).Description("describe what this argument does")
	return nil
}

func BasicCtxChained(ctx gogo.Context, var1 string, var2 bool) error {
	ctx.SetShortDescription("set a description, this can use any go code to set the value")

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

// is this a thing?
func BasicSetHelp(ctx gogo.Context, var1 string, var2 bool) error {
	return nil
}

// the below are not valid GoGo functions

// not valid because the first parameter is not a gogo.Context
func NoGoGoCtx() error {
	return nil
}

// not valid because the first parameter is not a gogo.Context
func GoGoCtxInWrongPos(var1 string, ctx gogo.Context, var2 bool) error {
	return nil
}

// not valid because we can only return an error
func WrongReturnType(ctx gogo.Context, var1 string, var2 bool) string {
	return ""
}
