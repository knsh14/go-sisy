package main

import (
	"bytes"
	"flag"
	"go/ast"
	"go/format"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/knsh14/go-sisy/converter"
	"github.com/pkg/errors"
)

var (
	overwrite bool
)

func init() {
	flag.BoolVar(&overwrite, "w", false, "overwrite to fixed code")
}
func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		switch dir, err := os.Stat(path); {
		case err != nil:
			log.Fatal(err)
		case dir.IsDir():
			filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				f, fset, err := converter.Convert(path)
				if err != nil {
					return errors.Wrap(err, "failed to convert")
				}
				err = write(path, fset, f)
				if err != nil {
					return errors.Wrap(err, "failed to write")
				}
				return nil
			})
		default:
			f, fset, err := converter.Convert(path)
			if err != nil {
				log.Fatal(err)
			}
			err = write(path, fset, f)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func write(path string, fset *token.FileSet, f *ast.File) error {
	var w bytes.Buffer
	err := format.Node(&w, fset, f)
	if err != nil {
		return errors.Wrap(err, "failed to write to buffer")
	}
	if overwrite {
		info, _ := os.Stat(path)
		return ioutil.WriteFile(path, w.Bytes(), info.Mode())
	}
	_, err = w.WriteTo(os.Stdout)
	return err
}
