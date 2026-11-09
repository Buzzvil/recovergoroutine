package recovergoroutine

import (
	"flag"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
)

type message string

type checker struct {
	customRecover string
	pass          *analysis.Pass
}

func NewAnalyzer() *analysis.Analyzer {
	c := checker{}
	analyzer := &analysis.Analyzer{
		Name: "recovergoroutine",
		Doc:  "finds goroutine code without recover",
		Run:  func(pass *analysis.Pass) (interface{}, error) { return c.run(pass) },
	}

	analyzer.Flags.Init("recovergoroutine", flag.ExitOnError)
	analyzer.Flags.StringVar(
		&c.customRecover,
		"recover",
		"",
		"You can use this option when you want to call a method defined in a struct or use CustomRecover declared in an external package.",
	)

	return analyzer
}

func (c checker) run(pass *analysis.Pass) (interface{}, error) {
	c.pass = pass
	var runErr error
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			goStmt, ok := n.(*ast.GoStmt)
			if !ok {
				return true
			}

			ok, msg := c.safeGoStmt(goStmt)
			if ok {
				return true
			}

			pass.Report(analysis.Diagnostic{
				Pos:      goStmt.Pos(),
				End:      0,
				Category: "goroutine",
				Message:  string(msg),
			})

			return true
		})
	}

	return nil, runErr
}

func (c checker) safeGoStmt(goStmt *ast.GoStmt) (bool, message) {
	fn := goStmt.Call
	switch fun := fn.Fun.(type) {
	case *ast.FuncLit:
		if !c.safeFunc(fun) {
			return false, "goroutine must have recover"
		}
		return true, ""
	}

	return false, "use function literals when using goroutines"
}

func (c checker) safeFunc(node *ast.FuncLit) bool {
	result := false
	ast.Inspect(node.Body, func(node ast.Node) bool {
		if _, ok := node.(*ast.FuncLit); ok {
			return false
		}
		deferStmt, ok := node.(*ast.DeferStmt)
		if !ok {
			return true
		}

		ok = c.hasRecover(deferStmt.Call)
		if ok {
			result = true
			return false
		}

		return !result
	})

	return result
}

func (c checker) hasRecover(expr ast.Node) bool {
	call := expr.(*ast.CallExpr)
	if id, ok := call.Fun.(*ast.Ident); ok && id.Name == "recover" {
		return false
	}
	if f, ok := call.Fun.(*ast.FuncLit); ok {
		expr = f.Body
	} else {
		return c.isRecover(call) || c.isCustomRecover(call)
	}
	var result bool
	ast.Inspect(expr, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.CallExpr:
			if c.isRecover(n) || c.isCustomRecover(n) {
				result = true
				return false
			}

		}
		return true
	})

	return result
}

func (c checker) isRecover(callExpr *ast.CallExpr) bool {
	if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
		return sel.Sel.Name == c.customRecover
	}
	ident, ok := callExpr.Fun.(*ast.Ident)
	if !ok {
		return false
	}

	return c.isBuiltinRecover(callExpr) || ident.Name == c.customRecover
}

func (c checker) isCustomRecover(callExpr *ast.CallExpr) bool {
	result := false
	switch fun := callExpr.Fun.(type) {
	case *ast.SelectorExpr:
		if fun.Sel == nil {
			return result
		}

		result = c.checkIdent(fun.Sel)
	case *ast.Ident:
		result = c.checkIdent(fun)
	}

	return result
}

func (c checker) checkIdent(ident *ast.Ident) bool {
	result := false
	if ident.Obj == nil {
		return result
	}

	funcDecl, ok := ident.Obj.Decl.(*ast.FuncDecl)
	if !ok {
		return result
	}

	ast.Inspect(funcDecl, func(node ast.Node) bool {
		if callExpr, ok := node.(*ast.CallExpr); ok && c.isRecover(callExpr) {
			result = true
			return false
		}

		return true
	})

	return result
}

func (c checker) isBuiltinRecover(call *ast.CallExpr) bool {
	id, ok := call.Fun.(*ast.Ident)
	return ok && c.pass.TypesInfo.Uses[id] == types.Universe.Lookup("recover")
}
