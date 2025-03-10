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

var _ Context = gogoContext{}
var _ Argument = gogoArgument{}

type Context interface {
	stdContext.Context
	SetShortDescription(short string) Context // This becomes the short description/usage of the command.
	Example(string) Context                   // What would this go to?
	Argument(any) Argument
}

type Argument interface {
	Name(string) Argument             // Override the argument name in auto-complete and arguments when calling this function. Defaults to the name of the argument.
	Short(byte) Argument              // The short character for the argument
	Default(any) Argument             // If set it's assumed the argument is also optional
	Optional() Argument               // If set the argument is optional
	Help(string) Argument             // Help for that specific argument. This is shown when inspecting the individual flag for information, or possibly when auto-completing in shell on positional/flag arguments.
	AllowedValues(...any) Argument    // Allowed values are checked in the command, and provide options for auto-complete in the shell. For now it's hard-coded values, but in the future could be regular expressions or even a go function.
	RestrictedValues(...any) Argument // Same as allowed values, but the values are not allowed. This is not used in the shell?
	Description(string) Argument      // The short description of the argument. This is used in flag descriptions
	Argument(any) Argument            // Start describing a different argument, allows for a builder pattern.
}

func NewContext() Context {
	return gogoContext{}
}

type gogoContext struct {
	stdContext.Context
}

type gogoArgument struct {
}

func (c gogoContext) SetShortDescription(short string) Context {
	return c
}

func (c gogoContext) Example(example string) Context {
	return c
}

func (c gogoContext) Argument(arg any) Argument {
	return &gogoArgument{}
}

func (a gogoArgument) Name(long string) Argument {
	return a
}

func (a gogoArgument) Short(b byte) Argument {
	return a
}

func (a gogoArgument) Default(value any) Argument {
	return a
}

func (a gogoArgument) Help(help string) Argument {
	return a
}

func (a gogoArgument) Optional() Argument {
	return a
}

func (a gogoArgument) AllowedValues(values ...any) Argument {
	return a
}

func (a gogoArgument) RestrictedValues(values ...any) Argument {
	return a
}

func (a gogoArgument) Description(description string) Argument {
	return a
}

func (a gogoArgument) Argument(arg any) Argument {
	return a
}
