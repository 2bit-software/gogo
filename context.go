// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package gogo

import (
	stdContext "context"
)

type Context interface {
	stdContext.Context
	SetDescription(string) Context     // This becomes the short description of the command.
	SetLongDescription(string) Context // The becomes the long description of the command. This is normally inferred from the comment of the function.
	Example(string) Context            // What would this go to?
	Argument(any) Argument
}

type Argument interface {
	Long(string) Argument // Override the argument name in auto-complete and arguments when calling this function. Defaults to the name of the argument.
	Short(byte) Argument  // The short character for the argument
	Default(any) Argument
	Help(string) Argument
	AllowedValues(...any) Argument
	RestrictedValues(...any) Argument
	Description(string) Argument // The short description of the argument.
	Argument(any) Argument       // Start describing a different argument.
}
