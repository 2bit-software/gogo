// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package sh

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func DetermineWidth(verbose bool) int {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		if verbose {
			fmt.Println("DEBUG: Running in a shell")
		}
		width, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			return -1
		}
		if verbose {
			fmt.Printf("DEBUG: Terminal width: %d\n", width)
		}
		return width
	}
	return -1
}
