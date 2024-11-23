// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package _gogo

import (
	"fmt"

	"github.com/2bit-software/gogo"
)

func AdvancedFunction(ctx gogo.Context, name string, include bool, value int) error {
	ctx.
		SetDescription("set a description").
		Example("example").
		SetLongDescription("this is the long description").
		Argument(name).
		Description("this is the name").
		Default("default-value").
		Argument(include).
		Description("this is the include bool").
		Default(true).
		Argument(value).
		RestrictedValues(1, 2, 3).
		AllowedValues(5, 6, 7).
		Description("this is the value").
		Default(3)

	fmt.Printf("name: %s\n", name)
	if include {
		fmt.Printf("value: %d\n", value)
	}

	return nil
}
