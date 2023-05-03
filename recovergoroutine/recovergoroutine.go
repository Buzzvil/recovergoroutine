package recovergoroutine

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "recovergoroutine",
	Doc:  "finds goroutine code without recover",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			goStmt, ok := n.(*ast.GoStmt)
			if !ok {
				return true
			}

			if safeGoStmt(goStmt) {
				return true
			}

			pass.Report(analysis.Diagnostic{
				Pos:      goStmt.Pos(),
				End:      0,
				Category: "goroutine",
				Message:  "goroutine must have recover",
			})

			return false
		})
	}
	return nil, nil
}

func safeGoStmt(goStmt *ast.GoStmt) bool {
	fn := goStmt.Call
	result := false
	if funcLit, ok := fn.Fun.(*ast.FuncLit); ok {
		result = safeFunc(funcLit)
	}

	if ident, ok := fn.Fun.(*ast.Ident); ok {
		if ident.Obj == nil {
			return true
		}

		funcDecl, ok := ident.Obj.Decl.(*ast.FuncDecl)
		if !ok {
			return true
		}

		result = safeFunc(funcDecl)
	}

	return result
}

func safeFunc(node ast.Node) bool {
	result := false
	ast.Inspect(node, func(node ast.Node) bool {
		deferStmt, ok := node.(*ast.DeferStmt)
		if !ok {
			return true
		}

		ast.Inspect(deferStmt.Call, func(node ast.Node) bool {
			callExpr, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			if isRecover(callExpr) {
				result = true
				return false
			}

			if isCustomRecover(callExpr) {
				result = true
				return false
			}

			return true
		})

		return !result
	})

	return result
}

func isRecover(callExpr *ast.CallExpr) bool {
	ident, ok := callExpr.Fun.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "recover"
}

func isCustomRecover(callExpr *ast.CallExpr) bool {
	result := false
	if ident, ok := callExpr.Fun.(*ast.Ident); ok {
		if ident.Obj == nil {
			return true
		}

		funcDecl, ok := ident.Obj.Decl.(*ast.FuncDecl)
		if !ok {
			return true
		}

		ast.Inspect(funcDecl, func(node ast.Node) bool {
			if callExpr, ok := node.(*ast.CallExpr); ok && isRecover(callExpr) {
				result = true
				return false
			}

			return true
		})
	}

	return result
}
