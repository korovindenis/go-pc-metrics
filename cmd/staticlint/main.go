// Staticlint provides essential static code analysis checks.
package main

import (
	"go/ast"
	"regexp"

	"github.com/gostaticanalysis/emptycase"
	"github.com/masibw/goone"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func main() {
	exitCheckAnalyzer := &analysis.Analyzer{
		Name: "exitcheck",
		Doc:  "custom check",
		Run:  run,
	}

	checks := []*analysis.Analyzer{
		printf.Analyzer,          // detect inconsistency of printf format strings and arguments.
		shadow.Analyzer,          // detect shadowed variables.
		structtag.Analyzer,       // detect error in structure tags.
		assign.Analyzer,          // detect useless assignments.
		atomic.Analyzer,          // check correct use sync package.
		bools.Analyzer,           // check error in boolean func.
		composite.Analyzer,       // check unkeyed composite literal.
		copylock.Analyzer,        // check error in copy some objects (Mutex, GroupWait and etc).
		deepequalerrors.Analyzer, // check incorrect use of deepequals.
		defers.Analyzer,          // check for defers.
		directive.Analyzer,       // check correct use directives.
		errorsas.Analyzer,        // check second args is error in errors.Is().
		fieldalignment.Analyzer,  // detect struct with less using memory if their fields is sorted.
		nilfunc.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unusedresult.Analyzer,
		goone.Analyzer, // check sql query in loop

	}

	for _, v := range staticcheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	for _, v := range stylecheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}
	checks = append(checks, emptycase.Analyzer) // check empty case body

	checks = append(checks, exitCheckAnalyzer) // check osExit
	multichecker.Main(
		checks...,
	)
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if skipFile(file) {
			continue
		}
		ast.Inspect(file, func(n ast.Node) bool {
			if x, ok := n.(*ast.CallExpr); ok {
				if exitCheck(pass, x) {
					pass.Reportf(x.Pos(), "call to os.Exit() found in main()")
				}
			}

			return true
		})
	}

	return nil, nil
}

func exitCheck(pass *analysis.Pass, x *ast.CallExpr) bool {
	if selExpr, ok := x.Fun.(*ast.SelectorExpr); ok {
		//nolint:todo
		if idExpr, ok := selExpr.X.(*ast.Ident); ok && idExpr.Name == "os" && selExpr.Sel.Name == "Exit" && pass.Pkg.Name() == "main" {
			for _, f := range pass.Files {
				if f.Name.Name == "main" {
					return true
				}
			}
		}
	}

	return false
}

func skipFile(file *ast.File) bool {
	regexp := regexp.MustCompile("^// Code generated .* DO NOT EDIT")
	for _, cg := range file.Comments {
		for _, comment := range cg.List {
			if regexp.MatchString(comment.Text) {
				return true
			}
		}
	}

	return false
}
