// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package gadgets

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type function struct {
	Name                string
	Comment             string // The long description of the function
	Description         string // Description for the function. This is the short description of the command.
	Example             string
	Arguments           []argument
	UseGoGoCtx          bool
	GoGoCtxVariableName string
	ErrorReturn         bool // does the function return an error?
}

type argument struct {
	Name             string
	Type             string
	Long             string // Override the argument name.
	Short            byte   // The short character for the argument
	Description      string // short description of the argument
	Help             string
	Default          any
	AllowedValues    []any
	RestrictedValues []any
}

const GOGOIMPORTPATH = "github.com/2bit-software/gogo/pkg/gogo"

// parseDirectory reads in a list of files and extracts the function information, aggregating it into a single list
// TODO: this might need to return a map of files/functions instead
func parseDirectory(dir string) ([]function, error) {
	var functions []function
	items, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	// for every file, parse it
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		// make sure it's a go file
		if filepath.Ext(item.Name()) != ".go" {
			continue
		}
		// make sure it's not the MAIN_FILENAME
		if item.Name() == MAIN_FILENAME {
			continue
		}
		funcs, err := parse(path.Join(dir, item.Name()))
		if err != nil {
			return nil, err
		}
		functions = append(functions, funcs...)
	}
	return functions, nil
}

func parseAll(files []string) ([]function, error) {
	var functions []function
	for _, file := range files {
		funcs, err := parse(file)
		if err != nil {
			return nil, err
		}
		functions = append(functions, funcs...)
	}
	return functions, nil
}

// parse reads in a source document and extracts the function information
func parse(filename string) ([]function, error) {
	// read the file
	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// parse the source
	return parseSource(string(f))
}

func parseSource(src string) ([]function, error) {
	// Parse the source code
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// Find the import alias for the gogo package
	gogoAlias, _ := getGoGoImportName(file)

	var functions []function
	// For each function, extract the information
	ast.Inspect(file, func(node ast.Node) bool {
		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok {
			return true
		}
		// Check if the function is exported, if not skip it
		if !funcDecl.Name.IsExported() {
			return true
		}
		// check if all the arguments are acceptable
		if !acceptableArguments(gogoAlias, funcDecl) {
			return true
		}
		// check if the function has an acceptable return type
		if !acceptableReturnTypes(funcDecl) {
			return true
		}

		if !gogoContextCorrectPosition(gogoAlias, funcDecl) {
			return true
		}

		// extract information from the function itself
		pCtx, err := parsePlainFunc(funcDecl)
		if err != nil {
			fmt.Println(err)
			return true
		}
		// save the function to the slice when we return
		defer func(pCtx *function) {
			functions = append(functions, *pCtx)
		}(pCtx)
		// fill out the arg
		pCtx = gatherDetails(pCtx, gogoAlias, funcDecl)

		// continue parsing the rest of the functions
		return true
	})

	return functions, nil
}

// gogoContextCorrectPosition checks if the gogo.Context is in the correct position
// If the function has a gogo.Context, it must be the first argument
func gogoContextCorrectPosition(alias string, decl *ast.FuncDecl) bool {
	for i, param := range decl.Type.Params.List {
		if isGoGoCtx(alias, param) && i == 0 {
			continue
		}
		if isGoGoCtx(alias, param) {
			return false
		}
	}
	return true
}

// acceptableReturnTypes checks if the function has an acceptable return type,
// which is either nothing or an error
func acceptableReturnTypes(funcDecl *ast.FuncDecl) bool {
	if funcDecl == nil {
		return false
	}
	if funcDecl.Type == nil {
		return false
	}
	if funcDecl.Type.Results == nil {
		return true
	}
	// Check if the return value is either empty
	if len(funcDecl.Type.Results.List) == 0 {
		return true
	}
	// or more than one value
	if len(funcDecl.Type.Results.List) > 1 {
		return false
	}
	return hasErrorReturn(funcDecl)
}

// acceptableArguments checks to make sure that all the arguments are scalar types, except
// for the gogo.Context, if it exists
func acceptableArguments(alias string, funcDecl *ast.FuncDecl) bool {
	if funcDecl == nil {
		return false
	}
	if funcDecl.Type == nil {
		return false
	}
	// no args, so it's fine
	if funcDecl.Type.Params == nil {
		return true
	}
	for _, param := range funcDecl.Type.Params.List {
		if isGoGoCtx(alias, param) {
			continue
		}
		if !isScalarType(param) {
			return false
		}
	}
	return true
}

// isScalarType checks if the given parameter is a scalar type
func isScalarType(param *ast.Field) bool {
	if param == nil || param.Type == nil {
		return false
	}

	// Identify the type node
	ident, ok := param.Type.(*ast.Ident)
	if !ok {
		return false
	}

	// Check if it's one of our defined scalar types
	switch ident.Name {
	case "string", "int", "float64", "bool":
		return true
	default:
		return false
	}
}

func hasErrorReturn(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Type.Results == nil {
		return false
	}
	if len(funcDecl.Type.Results.List) == 0 {
		return false
	}
	// detect if the return type is an error
	returnType, ok := funcDecl.Type.Results.List[0].Type.(*ast.Ident)
	if !ok || returnType.Name != "error" {
		return false
	}
	return true
}

// Extract information about each exported function. This includes the function name,
// the return type, and the function comment. Arguments get parsed later.
func parsePlainFunc(stmt *ast.FuncDecl) (*function, error) {
	pCtx := &function{
		Name: stmt.Name.Name,
	}
	if stmt.Doc != nil {
		pCtx.Comment = strings.TrimSpace(stmt.Doc.Text())
	}
	return pCtx, nil
}

// gatherDetails gets the argument and ctx.<method> information
func gatherDetails(pCtx *function, gogoAlias string, funcDecl *ast.FuncDecl) *function {
	// determine if this has an error return
	if hasErrorReturn(funcDecl) {
		pCtx.ErrorReturn = true
	}

	// if there no first argument, don't bother trying to parse for a GoGoContext
	// or any other arguments
	if len(funcDecl.Type.Params.List) == 0 {
		return pCtx
	}
	var args []argument
	if funcDecl.Type.Params == nil {
		return pCtx
	}

	// determine if the first argument is a gogo|<alias>.Context
	hasGoGoCtx := isFirstArgAliasContext(gogoAlias, funcDecl)
	if hasGoGoCtx {
		// set the GoGoCtxVariableName
		pCtx.GoGoCtxVariableName = funcDecl.Type.Params.List[0].Names[0].Name
	}

	// now we can parse the rest of the arguments
	for _, param := range funcDecl.Type.Params.List {
		for i, name := range param.Names {
			if hasGoGoCtx && i == 0 && name.Name == pCtx.GoGoCtxVariableName {
				continue
			}

			typ := GetPlainType(param)
			args = append(args, argument{
				Name: name.Name,
				Type: typ,
			})
		}
	}
	pCtx.Arguments = args

	if !hasGoGoCtx {
		return pCtx
	}

	// we know we have a GoGo context, so make signal it's imported at the very least
	pCtx.UseGoGoCtx = true

	// extract information using parseGoGoCtx
	pCtx, err := parseGoGoCtx(pCtx, funcDecl)
	if err != nil {
		// we ran into an error
		fmt.Printf("Error parsing GoGoContext: %v\n", err)
		return nil
	}

	return pCtx
}

// filesHaveFunc determines if any of the given files contain the given function
func filesHaveFunc(filePaths []string, funcName string) bool {
	for _, filePath := range filePaths {
		if fileHasFunc(filePath, funcName) {
			return true
		}
	}
	return false
}

func fileHasFunc(filePath string, funcName string) bool {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return false
	}
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if funcDecl.Name.Name == funcName {
			return true
		}
	}
	return false
}

// getGoGoImportName finds the alias or default name of the GoGo import
func getGoGoImportName(file *ast.File) (string, bool) {
	for _, imp := range file.Imports {
		if imp.Path.Value == fmt.Sprintf("\"%s\"", GOGOIMPORTPATH) {
			if imp.Name != nil {
				return imp.Name.Name, true
			}
			return "gogo", true
		}
	}
	return "", false
}

// isFirstArgAliasContext checks if the first argument to the FuncDecl is <alias>.Context (e.g., ctx context.Context)
func isFirstArgAliasContext(alias string, funcDecl *ast.FuncDecl) bool {
	// Check if the function has any parameters
	if funcDecl.Type.Params == nil || len(funcDecl.Type.Params.List) == 0 {
		return false
	}

	// Get the first parameter
	firstParam := funcDecl.Type.Params.List[0]

	return isGoGoCtx(alias, firstParam)
}

func isGoGoCtx(alias string, firstParam *ast.Field) bool {
	// The type of the first parameter should be a SelectorExpr with the form: alias.Context.
	if selExpr, ok := firstParam.Type.(*ast.SelectorExpr); ok {
		// The X field is the alias (it should be an Ident)
		if ident, ok := selExpr.X.(*ast.Ident); ok {
			// Check if the identifier name matches the alias and the selector is "Context"
			return ident.Name == alias && selExpr.Sel.Name == "Context"
		}
	}
	return false
}

// GetPlainType takes an *ast.Field and returns a plain-English type description
func GetPlainType(field *ast.Field) string {
	if field == nil || field.Type == nil {
		return ""
	}
	return exprToTypeStr(field.Type)
}

// exprToTypeStr is a helper function to traverse the AST and get the string version of various Go types.
func exprToTypeStr(expr ast.Expr) string {
	switch t := expr.(type) {

	case *ast.Ident: // Identifiers, like "int", "string", "Context", etc.
		return t.Name

	case *ast.StarExpr: // Pointers, like "*context.Context", "*MyStruct"
		return "*" + exprToTypeStr(t.X)

	case *ast.SelectorExpr: // Selector, like "context.Context", "time.Time"
		return exprToTypeStr(t.X) + "." + t.Sel.Name

	case *ast.ArrayType: // Array types, like "[]int", "[]*MyStruct"
		return "[]" + exprToTypeStr(t.Elt)

	case *ast.MapType: // Map types, like "map[string]int"
		keyType := exprToTypeStr(t.Key)
		valueType := exprToTypeStr(t.Value)
		return "map[" + keyType + "]" + valueType

	case *ast.InterfaceType: // Interface types, like "interface{}"
		if t.Methods != nil && t.Methods.List == nil {
			return "interface{}"
		}
		// Construct the interface type as needed, but for now, default to `interface{}`
		return "interface{}"

	case *ast.FuncType: // Function types, like func(string) (int, error)
		params := typeListToStr(t.Params)
		results := typeListToStr(t.Results)

		if results != "" {
			return "func(" + params + ") (" + results + ")"
		}
		return "func(" + params + ")"

	case *ast.ChanType: // Chan types, like chan int, <-chan string, chan<- int
		dir := "chan "
		if t.Dir == ast.RECV {
			dir = "<-chan "
		} else if t.Dir == ast.SEND {
			dir = "chan<- "
		}
		return dir + exprToTypeStr(t.Value)

	default:
		// Other types like structs, function types, etc.
		return types.ExprString(expr) // Fallback using types.ExprString
	}
}

// typeListToStr converts a FieldList (parameters or return types) to a string representation.
func typeListToStr(list *ast.FieldList) string {
	if list == nil || len(list.List) == 0 {
		return ""
	}
	var parts []string
	for _, field := range list.List {
		typeStr := exprToTypeStr(field.Type)
		parts = append(parts, typeStr)
	}
	return strings.Join(parts, ", ")
}
