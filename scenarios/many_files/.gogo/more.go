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

func MoreFunction(ctx gogo.Context, input string) {
	ctx.
		SetShortDescription("a description").
		Example("example 2").
		Argument(input).
		Description("this is the input!").
		Default("default-value-2")

	fmt.Println("input:", input)
}
