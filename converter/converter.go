package converter

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"

	"github.com/knsh14/go-sisy/definition"
	"github.com/pkg/errors"
)

// Convert returns fixed file and file set.
func Convert(filePath string) (*ast.File, *token.FileSet, error) {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to parse file")
	}
	importMap := map[string]string{}
	for _, v := range f.Imports {
		p, err := strconv.Unquote(v.Path.Value)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to unquote")
		}
		key := filepath.Base(p)
		if v.Name != nil {
			key = v.Name.Name
		}
		importMap[key] = p
	}

	// TODO: for じゃなくて walk で実装する
	for _, v := range f.Decls {
		fnc, ok := v.(*ast.FuncDecl)
		if ok {
			for _, stmt := range fnc.Body.List {
				assign, ok := stmt.(*ast.AssignStmt)
				if ok {
					for _, rh := range assign.Rhs {
						switch t := rh.(type) {
						case *ast.UnaryExpr:
							v, err := convertLiteral(filepath.Dir(filePath), t.X.(*ast.CompositeLit), importMap)
							if err != nil {
								return nil, nil, errors.Wrap(err, "failed to convert")
							}
							t.X = v
						case *ast.CompositeLit:
							v, err := convertLiteral(filepath.Dir(filePath), t, importMap)
							if err != nil {
								return nil, nil, errors.Wrap(err, "failed to convert")
							}
							t = v
						}
					}
				}
			}
		}
	}
	return f, fset, nil
}

func convertLiteral(f string, lit *ast.CompositeLit, dict map[string]string) (*ast.CompositeLit, error) {
	_, ok := lit.Elts[0].(*ast.KeyValueExpr)
	if ok {
		return lit, nil
	}
	fields := []string{}
	switch t := lit.Type.(type) {
	case *ast.Ident:
		if t.Obj == nil {
			// TODO どうにかして自分のパッケージパスを取ってこないと
			ast.Print(token.NewFileSet(), f)
			fs, err := definition.GetExportedFields(f, t.Name)
			if err != nil {
				return lit, errors.Wrap(err, "failed to get strct exported fields")
			}
			fields = fs
		} else {
			ts, ok := t.Obj.Decl.(*ast.TypeSpec)
			if !ok {
				return lit, errors.Errorf("%s is not TypeSpec", t.Obj.Name)
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				return lit, errors.Errorf("not StructType")
			}
			if st.Incomplete {
				return lit, errors.Errorf("source incomplete")
			}
			for _, f := range st.Fields.List {
				for _, n := range f.Names {
					fields = append(fields, n.Name)
				}
			}
		}
	case *ast.SelectorExpr:
		fs, err := definition.GetExportedFields(dict[t.X.(*ast.Ident).Name], t.Sel.Name)
		if err != nil {
			return lit, errors.Wrap(err, "failed to get strct exported fields")
		}
		fields = fs
	}

	if len(lit.Elts) != len(fields) {
		return lit, errors.Errorf("field num is not same. fields = %d, names = %d", len(lit.Elts), len(fields))
	}
	literals := []ast.Expr{}
	for i := range lit.Elts {
		literals = append(literals, &ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: fields[i],
			},
			Value: lit.Elts[i],
		})
	}
	lit.Elts = literals
	return lit, nil
}
