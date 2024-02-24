package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"math/rand"
	"os"
	"path"

	"github.com/tsuna/gorewrite"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// tmpFileName creates a new file
func tmpFileName(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getValuesForRange(stmt *ast.ForStmt) (string, string) {
	variable := ""

	// Step 1) Is the init correct?
	assignStmt, ok := stmt.Init.(*ast.AssignStmt)
	//fmt.Println(assignStmt.Tok)
	if !ok {
		return "", ""
	}

	// LHS of statement must be a single variable

	if len(assignStmt.Lhs) != 1 {
		return "", ""
	}

	variableStmt, ok := assignStmt.Lhs[0].(*ast.Ident)
	if !ok {
		return "", ""
	}
	variable = variableStmt.Name

	// RHS of statement must be ":= 0"

	if assignStmt.Tok != token.DEFINE {
		return "", ""
	}
	if len(assignStmt.Rhs) != 1 {
		return "", ""
	}
	assignToStmt, ok := assignStmt.Rhs[0].(*ast.BasicLit)
	if !ok {
		return "", ""
	}
	if assignToStmt.Kind != token.INT {
		return "", ""
	}
	if assignToStmt.Value != "0" {
		return "", ""
	}

	// Step 2) Is the condition correct?

	binaryExpr, ok := stmt.Cond.(*ast.BinaryExpr)
	if !ok {
		return "", ""
	}
	if binaryExpr.Op != token.LSS {
		return "", ""
	}
	lessThanExpr, ok := binaryExpr.X.(*ast.Ident)
	if !ok {
		return "", ""
	}
	if lessThanExpr.Name != variable {
		return "", ""
	}
	rangeExpr, ok := binaryExpr.Y.(*ast.BasicLit)
	if !ok {
		return "", ""
	}
	if rangeExpr.Kind != token.INT {
		return "", ""
	}
	rangeValue := rangeExpr.Value

	// Step 3) Is post correct?
	incExp, ok := stmt.Post.(*ast.IncDecStmt)
	if !ok {
		return "", ""
	}
	if incExp.Tok != token.INC {
		return "", ""
	}
	variableIncremented, ok := incExp.X.(*ast.Ident)
	if !ok {
		return "", ""
	}
	if variableIncremented.Name != variable {
		return "", ""
	}

	return variable, rangeValue
}

func fixOneFile(filename string) error {
	_, err := os.Stat(filename)
	if err != nil {
		return err
	}
	dir := path.Dir(filename)

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	replacer := forLoopWithIntReplacer{}
	gorewrite.Rewrite(&replacer, node)

	// No updates, no need to write file
	if !replacer.updated {
		log.Printf("No updates needed for %s", filename)
		return nil
	}
	log.Printf("Updating %s", filename)

	outputFile, err := os.CreateTemp(dir, fmt.Sprintf("%s.new.", path.Base(filename)))
	if err != nil {
		log.Fatal(err)
	}
	outputFileName := outputFile.Name()

	defer func() {
		outputFile.Close()
		_, err := os.Stat(outputFileName)
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		if err != nil {
			panic(err)
		}
		fmt.Println("DELETING")
		err = os.Remove(outputFileName)
		if err != nil {
			panic(err)
		}
	}()

	// Format the modified AST and write it to the file.
	if err := format.Node(outputFile, fset, node); err != nil {
		return err
	}
	err = os.Rename(outputFileName, filename)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := fixOneFile("foo/example.go")
	if err != nil {
		log.Fatal(err)
	}

}

type forLoopWithIntReplacer struct {
	updated bool
}

func (v *forLoopWithIntReplacer) Rewrite(n ast.Node) (ast.Node, gorewrite.Rewriter) {
	forStmt, ok := n.(*ast.ForStmt)
	if !ok {
		return n, v
	}
	variable, numToRange := getValuesForRange(forStmt)
	if variable == "" || numToRange == "" {
		return n, v
	}
	v.updated = true
	return &ast.RangeStmt{
		Key: &ast.Ident{
			Name: variable,
		},
		Tok: token.DEFINE,
		X: &ast.BasicLit{
			Value: numToRange,
		},
		Body: forStmt.Body,
	}, v

}
