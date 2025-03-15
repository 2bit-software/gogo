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

// WARNING: This file contains comments and magic strings that are validated against in tests!

// Description This is the description for the function. Without any other arguments to the ctx,
// this will show up in the list view and the --help output.
func Description(ctx gogo.Context) error {
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
