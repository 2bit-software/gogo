// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package gogo

import (
	"os"

	"github.com/spf13/cobra"
)

// TODO: we need to write our own versions of these, but for now...
func GenerateShellCompletion(shell string, cmd *cobra.Command) error {
	switch shell {
	case "bash":
		return cmd.GenBashCompletion(os.Stdout)
	case "zsh":
		return cmd.GenZshCompletion(os.Stdout)
	case "fish":
		return cmd.GenFishCompletion(os.Stdout, true)
	case "powershell":
		return cmd.GenPowerShellCompletion(os.Stdout)
	default:
		return cmd.GenZshCompletion(os.Stdout)
	}
}
