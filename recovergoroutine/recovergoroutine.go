package recovergoroutine

import (
	"flag"
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

type message string

var customRecover string

func NewAnalyzer() *analysis.Analyzer {
	analyzer := &analysis.Analyzer{
		Name: "recovergoroutine",
		Doc:  "finds goroutine code without recover",
		Run:  run,
	}

	analyzer.Flags.Init("recovergoroutine", flag.ExitOnError)
	analyzer.Flags.StringVar(
		&customRecover,
		"recover",
		"",
		"You can use this option when you want to call a method defined in a struct or use CustomRecover declared in an external package.",
	)

	return analyzer
}

func run(pass *analysis.Pass) (interface{}, error) {
	var runErr error
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			goStmt, ok := n.(*ast.GoStmt)
			if !ok {
				return true
			}

			ok, msg := safeGoStmt(goStmt)
			if ok {
				return true
			}

			pass.Report(analysis.Diagnostic{
				Pos:      goStmt.Pos(),
				End:      0,
				Category: "goroutine",
				Message:  string(msg),
			})

			return false
		})
	}

	return nil, runErr
}

func safeGoStmt(goStmt *ast.GoStmt) (bool, message) {
	fn := goStmt.Call
	switch fun := fn.Fun.(type) {
	case *ast.FuncLit:
		if !safeFunc(fun) {
			return false, "goroutine must have recover"
		}
		return true, ""
	}

	return false, "use function literals when using goroutines"
}

func safeFunc(node ast.Node) bool {
	result := false
	ast.Inspect(node, func(node ast.Node) bool {
		deferStmt, ok := node.(*ast.DeferStmt)
		if !ok {
			return true
		}

		ok = hasRecover(deferStmt.Call)
		if ok {
			result = true
			return false
		}

		return !result
	})

	return result
}

func hasRecover(expr ast.Node) bool {
	var result bool
	ast.Inspect(expr, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.CallExpr:
			if isRecover(n) || isCustomRecover(n) {
				result = true
				return false
			}
		case *ast.SelectorExpr:
			if n.Sel == nil {
				return true
			}

			if n.Sel.Name == customRecover {
				result = true
				return false
			}
		}
		return true
	})

	return result
}

func isRecover(callExpr *ast.CallExpr) bool {
	ident, ok := callExpr.Fun.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "recover" || ident.Name == customRecover
}

func isCustomRecover(callExpr *ast.CallExpr) bool {
	result := false
	switch fun := callExpr.Fun.(type) {
	case *ast.SelectorExpr:
		if fun.Sel == nil {
			return result
		}

		result = checkIdent(fun.Sel)
	case *ast.Ident:
		result = checkIdent(fun)
	}

	return result
}

func checkIdent(ident *ast.Ident) bool {
	result := false
	if ident.Obj == nil {
		return result
	}

	funcDecl, ok := ident.Obj.Decl.(*ast.FuncDecl)
	if !ok {
		return result
	}

	ast.Inspect(funcDecl, func(node ast.Node) bool {
		if callExpr, ok := node.(*ast.CallExpr); ok && isRecover(callExpr) {
			result = true
			return false
		}

		return true
	})

	return result
}
