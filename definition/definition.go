package definition

import (
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/pkg/errors"
)

// GetExportedFields returns field names that is exported
func GetExportedFields(packageName, structName string) ([]string, error) {
	fset := token.NewFileSet()
	v, err := build.Default.Import(packageName, "", 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to import package")
	}

	files := []*ast.File{}
	for _, s := range v.GoFiles {
		f, err := parser.ParseFile(fset, v.Dir+"/"+s, nil, parser.Mode(0))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse file %s", s)
		}
		files = append(files, f)
	}

	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check(v.Name, fset, files, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check")
	}
	s := pkg.Scope().Lookup(structName)
	if s == nil {
		return nil, errors.Errorf("%s is nil", structName)
	}
	internal, ok := s.Type().Underlying().(*types.Struct)
	if !ok {
		return nil, errors.Errorf("%s is not struct", structName)
	}

	fields := []string{}
	for i := 0; i < internal.NumFields(); i++ {
		field := internal.Field(i)
		if field.Exported() {
			fields = append(fields, field.Name())
		}
	}
	return fields, nil
}
