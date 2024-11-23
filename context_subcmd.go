package gogo

import (
	"context"
)

type subCmdContext struct {
	context.Context
}

type subCmdArgument struct {
}

func (c subCmdContext) SetDescription(short string) Context {
	return c
}

func (c subCmdContext) SetLongDescription(long string) Context {
	return c
}

func (c subCmdContext) Example(example string) Context {
	return c
}

func (c subCmdContext) Argument(arg any) Argument {
	return &subCmdArgument{}
}

func (a subCmdArgument) Long(long string) Argument {
	return a
}

func (a subCmdArgument) Short(b byte) Argument {
	return a
}

func (a subCmdArgument) Default(value any) Argument {
	return a
}

func (a subCmdArgument) Help(help string) Argument {
	return a
}

func (a subCmdArgument) AllowedValues(values ...any) Argument {
	return a
}

func (a subCmdArgument) RestrictedValues(values ...any) Argument {
	return a
}

func (a subCmdArgument) Description(description string) Argument {
	return a
}

func (a subCmdArgument) Argument(arg any) Argument {
	return a
}
