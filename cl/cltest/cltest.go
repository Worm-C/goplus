package cltest

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/qiniu/goplus/ast"
	"github.com/qiniu/goplus/cl"
	"github.com/qiniu/goplus/parser"
	"github.com/qiniu/goplus/token"
	"github.com/qiniu/x/log"

	exec "github.com/qiniu/goplus/exec/bytecode"
	_ "github.com/qiniu/goplus/lib" // libraries
)

// -----------------------------------------------------------------------------

func getPkg(pkgs map[string]*ast.Package) *ast.Package {
	for _, pkg := range pkgs {
		return pkg
	}
	return nil
}

func testFrom(t *testing.T, pkgDir, sel, exclude string) {
	if sel != "" && !strings.Contains(pkgDir, sel) {
		return
	}
	if exclude != "" && strings.Contains(pkgDir, exclude) {
		return
	}
	log.Debug("Compiling", pkgDir)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgDir, nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseDir failed:", err, len(pkgs))
	}

	bar := getPkg(pkgs)
	b := exec.NewBuilder(nil)
	_, err = cl.NewPackage(b.Interface(), bar, fset, cl.PkgActClMain)
	if err != nil {
		if err == cl.ErrNotAMainPackage {
			return
		}
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
}

// FromTestdata - run test cases from a directory
func FromTestdata(t *testing.T, dir, sel, exclude string) {
	cl.CallBuiltinOp = exec.CallBuiltinOp
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal("ReadDir failed:", err)
	}
	for _, fi := range fis {
		testFrom(t, dir+"/"+fi.Name(), sel, exclude)
	}
}

// -----------------------------------------------------------------------------
