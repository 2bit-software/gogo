// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package main

import (
	"fmt"
	"time"

	"github.com/2bit-software/gogo/pkg/gogo"
	"github.com/morganhein/gogo/_gogo/pkg"
)

// PrintShortSha prints the sha of the current/latest commit using only the first 6 characters
func PrintShortSha() {
	sha, err := pkg.GetCurrentShortSha()
	if err != nil {
		fmt.Println("could not get current SHA")
		return
	}
	fmt.Print(sha)
}

// CompileGo builds a Go binary from the given input folder path.
// Optionally, you can provide an outputFolderPath, tags, binaryName, versionPath, strOs, and strArch.
// The 'versionPath' is the path to use to set the build time arguments, which is set via LDFLAGS.
func CompileGo(ctx gogo.Context, inputFolderPath, outputFolderPath, tags, binaryName, versionPath, strOs, strArch string) error {
	ctx.ShortDescription("Compiles a Go binary, assuming filePath contains a main package.")
	return pkg.CompileGo(inputFolderPath, outputFolderPath, tags, binaryName, versionPath, strOs, strArch)
}

// Hello prints a greeting to the console
func Hello(input string, count int) {
	fmt.Printf("Hello, %v. Can you count to %v?", input, count)
}

func LongRunning() {
	for i := 0; i < 10; i++ {
		fmt.Printf("Counting: %d\n", i)
		time.Sleep(1 * time.Second)
	}
}
