package gogo

// HiddenArgsKey for the context value that will store the hidden args in the context
type HiddenArgsKey struct{}

// ParseHiddenFlags extracts flags before -- and returns the remaining args
func ParseHiddenFlags(args []string) ([]string, []string) {
	var separatorIndex = -1

	// Find the -- separator
	// TODO: this needs to be smarter, we may want to parse the flags using
	// 	real flag parsing, so we don't accidentally parse a "--" in a string
	for i, arg := range args {
		if arg == "--" {
			separatorIndex = i
			break
		}
	}

	// If no separator found, return original args as all subCmd args
	if separatorIndex == -1 {
		return nil, args
	}

	gogoFlags := args[separatorIndex+1:]
	subCmdArgs := args[:separatorIndex]

	return gogoFlags, subCmdArgs
}
