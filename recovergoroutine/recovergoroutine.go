package recovergoroutine

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "recovergoroutine",
	Doc:  "finds goroutine code without recover",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	var runErr error
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			goStmt, ok := n.(*ast.GoStmt)
			if !ok {
				return true
			}

			ok, err := safeGoStmt(goStmt, pass)
			if err != nil {
				runErr = err
				return false
			}

			if ok {
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

	return nil, runErr
}

func safeGoStmt(goStmt *ast.GoStmt, pass *analysis.Pass) (bool, error) {
	fn := goStmt.Call
	result := false
	switch fun := fn.Fun.(type) {
	case *ast.SelectorExpr:
		ident, ok := fun.X.(*ast.Ident)
		if !ok {
			return false, nil
		}

		methodName := fun.Sel.Name
		objType := pass.TypesInfo.ObjectOf(ident)
		pointerType, ok := objType.Type().(*types.Pointer)
		if !ok {
			return false, nil
		}

		named, ok := pointerType.Elem().(*types.Named)
		if !ok {
			return false, nil
		}

		for i := 0; i < named.NumMethods(); i++ {
			if named.Method(i).Name() != methodName {
				continue
			}

			fset := pass.Fset
			position := fset.Position(named.Method(i).Pos())
			file, err := parser.ParseFile(fset, position.Filename, nil, 0)
			if err != nil {
				return false, fmt.Errorf("parse file: %w", err)
			}

			for _, decl := range file.Decls {
				if funcDecl, ok := decl.(*ast.FuncDecl); ok {
					if funcDecl.Name.Name == methodName {
						result = safeFunc(funcDecl)
					}
				}
			}
		}
	case *ast.FuncLit:
		result = safeFunc(fun)
	case *ast.Ident:
		if fun.Obj == nil {
			return false, nil
		}

		funcDecl, ok := fun.Obj.Decl.(*ast.FuncDecl)
		if !ok {
			return false, nil
		}

		result = safeFunc(funcDecl)
	}

	return result, nil
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

			if isRecover(callExpr) || isCustomRecover(callExpr) {
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
