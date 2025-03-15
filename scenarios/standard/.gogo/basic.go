// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package main

import (
	"fmt"
)

// WARNING: This file contains comments and magic strings that are validated against in tests!

// This file tests the basic functionality of the gogo package.
// Primarily, functions that do not use the GoGo context to set values.

// The function has no description, no arguments, and no return values. KEEP THE EXTRA NEWLINE FOLLOWING THIS COMMENT!

func NoArgumentsNoReturns() {
	fmt.Println("NoArgumentsNoReturns")
}

// DescriptionOnly This is the description for the function. Without any other arguments to the ctx,
// this will show up in the list view and the --help output.
func DescriptionOnly() {
	fmt.Println("DescriptionOnly")
}

// ErrorReturn requires no arguments, but returns an error.
func ErrorReturn() error {
	fmt.Println("ErrorReturn")
	return nil
}

// SingleArgument tests a single argument.
func SingleArgument(arg1 string) {
	fmt.Printf("SingleArgument with arg1: %v\n", arg1)
}

// SingleArgumentAndErrorReturn tests a single argument and returns an error.
func SingleArgumentAndErrorReturn(arg1 string) error {
	fmt.Printf("SingleArgumentAndErrorReturn with arg1: %v\n", arg1)
	return nil
}

// TwoDifferentArguments tests two different arguments.
func TwoDifferentArguments(arg1 string, arg2 bool) {
	fmt.Printf("TwoDifferentArguments with arg1: %v, arg2: %v\n", arg1, arg2)
}

// TwoDifferentArgumentsAndErrorReturn tests two different arguments and returns an error.
func TwoDifferentArgumentsAndErrorReturn(arg1 string, arg2 bool) error {
	fmt.Printf("TwoDifferentArgumentsAndErrorReturn with arg1: %v, arg2: %v\n", arg1, arg2)
	return nil
}
