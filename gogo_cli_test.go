package scripts

import "testing"

func TestRunCommands(t *testing.T) {
	tests := []struct {
		name     string
		folder   string
		command  string
		args     []string
		expected string
	}{
		{
			name:     "NoArgumentsNoReturns",
			folder:   "standard",
			command:  "NoArgumentsNoReturns",
			args:     []string{},
			expected: "NoArgumentsNoReturns",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunCommands(test.folder, test.command, test.args, test.expected)
		})
	}
}

/*
NoArgumentsNoReturns,
DescriptionOnly,
ErrorReturn,
SingleArgument,
SingleArgumentAndErrorReturn,
TwoDifferentArguments,
TwoDifferentArgumentsAndErrorReturn,
*/
