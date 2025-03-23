// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package gadgets

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"slices"
	"strconv"
	"strings"
)

type call struct {
	FuncName string
	Args     []any
	Next     *call
	Previous *call
}

// parseGoGoCtx assumes that there are at least 1 calls to the pCtx.GoGoCtxVariableName in the function.
// If none are found, the original function is returned. This can happen if they specify a gogo.Context in the function
// signature but don't end up using it. Eventually, this should parse *all* uses of the context, but for now, it just
// parses the first one.
func parseGoGoCtx(pCtx *function, funcDecl *ast.FuncDecl) (*function, error) {
	stmnt, found := findFirstUsageOfChain(funcDecl, pCtx.GoGoCtxVariableName)
	if !found {
		// first argument is ctx, but it's not used
		return pCtx, nil
	}

	invertedChain := invertCallChain(stmnt)
	if invertedChain == nil {
		return nil, errors.New("could not invert call chain")
	}

	// We now walk the chain for "stmnt", which is "ctx" and its subsequent method calls:
	err := processGoGoChain(invertedChain, pCtx)
	if err != nil {
		return nil, err
	}

	pCtx.UseGoGoCtx = true
	return pCtx, nil
}

// processGoGoChain walks through the root of the expression chain and identifies method calls on ctx
func processGoGoChain(ctx *call, pCtx *function) error {
	current := ctx
	var err error
	for current != nil {
		// If the function is a method on the context, process it
		// Process the method
		current, err = processMethodOnContext(current, pCtx)
		if err != nil {
			return err
		}
	}
	return nil
}

func processMethodOnContext(current *call, ctx *function) (*call, error) {
	// Process the method
	switch current.FuncName {
	case "ShortDescription":
		if len(current.Args) == 1 {
			ctx.Description = current.Args[0].(string)
		}
	case "Example":
		if len(current.Args) == 1 {
			ctx.Example = current.Args[0].(string)
		}
	case "Argument":
		if len(current.Args) == 1 {
			argName := current.Args[0].(string)
			var err error
			current, ctx.Arguments, err = processMethodOnArgument(current.Next, argName, ctx.Arguments)
			if err != nil {
				return nil, err
			}
			return current, nil
		}
	}
	return current.Next, nil
}

func processMethodOnArgument(current *call, argName string, args []argument) (*call, []argument, error) {
	// check if this argument exists in the arg map already, as it should
	argIndex := slices.IndexFunc(args, func(a argument) bool {
		return a.Name == argName
	})
	if argIndex == -1 {
		return nil, nil, fmt.Errorf("argument %q not found in argument map", argName)
	}
	arg := args[argIndex]
	for current != nil {
		// begin walking the current call chain
		switch current.FuncName {
		case "Name":
			if len(current.Args) == 1 {
				arg.Long = current.Args[0].(string)
			}
		case "Short":
			if len(current.Args) == 1 {
				// convert the string to a byte
				x := current.Args[0].(string)
				// strip the single quotes
				x = strings.Replace(x, "'", "", -1)
				if len(x) != 1 {
					return nil, nil, fmt.Errorf("expected a single byte, got %q", x)
				}
				arg.Short = x[0]
			}
		case "Default":
			if len(current.Args) == 1 {
				arg.Default = current.Args[0]
			}
		case "Help":
			if len(current.Args) == 1 {
				arg.Help = current.Args[0].(string)
			}
		case "AllowedValues":
			if len(current.Args) > 0 {
				arg.AllowedValues = current.Args
			}
		case "RestrictedValues":
			if len(current.Args) > 0 {
				arg.RestrictedValues = current.Args
			}
		case "Description":
			if len(current.Args) == 1 {
				arg.Description = current.Args[0].(string)
			}
		case "Argument":
			// It's a new argument, return and let the caller handle it
			args[argIndex] = arg
			return current, args, nil
		}
		current = current.Next
	}
	// find the argument with this name, and set it
	args[argIndex] = arg
	return current, args, nil
}

func invertCallChain(expr ast.Expr) *call {
	var root *call

	for {
		switch node := expr.(type) {
		case *ast.CallExpr:
			// Create a new call node
			newCall := &call{
				FuncName: extractFuncName(node.Fun),
				Args:     extractArgs(node.Args),
			}

			if root != nil {
				newCall.Next = root
				root.Previous = newCall
			}

			root = newCall

			// Move up the chain
			if sel, ok := node.Fun.(*ast.SelectorExpr); ok {
				expr = sel.X
			} else {
				return root
			}

		case *ast.SelectorExpr:
			expr = node.X

		default:
			// We've reached the top of the chain
			return root
		}
	}
}

func extractFuncName(expr ast.Expr) string {
	switch node := expr.(type) {
	case *ast.Ident:
		return node.Name
	case *ast.SelectorExpr:
		return node.Sel.Name
	default:
		return ""
	}
}

func extractArgs(args []ast.Expr) []any {
	var result []any
	for _, arg := range args {
		switch node := arg.(type) {
		case *ast.BasicLit:
			if node.Kind == token.STRING {
				// Remove the surrounding quotes from the string literal
				unquoted, err := strconv.Unquote(node.Value)
				if err != nil {
					result = append(result, node.Value)
				} else {
					result = append(result, unquoted)
				}
			} else {
				result = append(result, node.Value)
			}
		case *ast.Ident:
			result = append(result, node.Name)
		// Add more cases as needed for other types of arguments
		default:
			result = append(result, fmt.Sprintf("%T", node))
		}
	}
	return result
}

// findFirstUsageOfChain locates the leaf method chain rooted in the provided context variable argument name.
func findFirstUsageOfChain(funcDecl *ast.FuncDecl, argName string) (ast.Expr, bool) {
	// Traverse the function body statements
	for _, stmt := range funcDecl.Body.List {
		// Recursively inspect the statement to find method chains using the specified argument.
		if foundExpr := findUsageInStmt(stmt, argName); foundExpr != nil {
			// convert stmt to expr
			switch s := stmt.(type) {
			case *ast.ExprStmt:
				return s.X, true
			}
		}
	}
	// No usage found
	return nil, false
}

// findUsageInStmt inspects a given statement to find method calls rooted in the specified argument name.
func findUsageInStmt(stmt ast.Stmt, argName string) ast.Expr {
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		// Expression statements (i.e., expressions as statements)
		return findUsageInExpr(s.X, argName)
	case *ast.AssignStmt:
		// Assignment statements: check RHS (the expressions being assigned)
		for _, rhs := range s.Rhs {
			if expr := findUsageInExpr(rhs, argName); expr != nil {
				return expr
			}
		}
	case *ast.ReturnStmt:
		// Return statements: check returned expressions
		for _, result := range s.Results {
			if expr := findUsageInExpr(result, argName); expr != nil {
				return expr
			}
		}
	case *ast.IfStmt:
		// If statements: check condition, body, and else parts
		if expr := findUsageInExpr(s.Cond, argName); expr != nil {
			return expr
		}
		// Check body of if
		for _, stmt := range s.Body.List {
			if expr := findUsageInStmt(stmt, argName); expr != nil {
				return expr
			}
		}
		// Check else part (either a block or another statement)
		if s.Else != nil {
			if expr := findUsageInStmt(s.Else, argName); expr != nil {
				return expr
			}
		}
	case *ast.ForStmt:
		// Check in for loop: init, condition, and post expressions
		if expr := findUsageInStmt(s.Init, argName); expr != nil {
			return expr
		}
		if expr := findUsageInExpr(s.Cond, argName); expr != nil {
			return expr
		}
		for _, stmt := range s.Body.List {
			if expr := findUsageInStmt(stmt, argName); expr != nil {
				return expr
			}
		}
	}
	return nil
}

// findUsageInExpr checks for method chains or expressions rooted in the context variable name.
func findUsageInExpr(expr ast.Expr, argName string) ast.Expr {
	switch e := expr.(type) {
	case *ast.CallExpr:
		// Check if a function call (method) expression is rooted in the desired argument name
		if selExpr, ok := e.Fun.(*ast.SelectorExpr); ok {
			if rootIsArg(selExpr.X, argName) {
				return e // Return the method call if it's rooted in the context variable
			}
			return findUsageInExpr(selExpr.X, argName) // Recursively check the method receiver
		}
	case *ast.SelectorExpr:
		// Check a selector expression (e.g., obj.Method)
		if rootIsArg(e.X, argName) {
			return e // Return the selector expression if rooted in the context variable
		}
		return findUsageInExpr(e.X, argName) // Recursively check the receiver
	}
	return nil
}

// rootIsArg walks up the tree of expressions to check if the root is the specified argument name.
func rootIsArg(expr ast.Expr, argName string) bool {
	if ident, ok := expr.(*ast.Ident); ok && ident.Name == argName {
		return true
	}
	return false
}
