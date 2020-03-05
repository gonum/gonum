// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package autofd

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Func struct {
	Path  string // import path of the package holding the function
	Name  string // function or method name
	Deriv string // name of the output derivative function
}

// Derivative generates the derivative code from the given function declaration.
func Derivative(w io.Writer, f Func) error {
	gen, err := newGenerator(w, f)
	if err != nil {
		return fmt.Errorf("could not create derivative generator: %w", err)
	}
	err = gen.generate()
	if err != nil {
		return fmt.Errorf("could not generate derivative: %w", err)
	}
	return nil
}

type generator struct {
	w   io.Writer
	pkg *packages.Package
	fct *types.Func
	der string
	err error
}

func newGenerator(w io.Writer, f Func) (*generator, error) {
	path := f.Path
	name := f.Name

	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo,
	}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		return nil, fmt.Errorf("could not load package of %q: %w", f, err)
	}

	pkg := func() *packages.Package {
		for _, p := range pkgs {
			if p.PkgPath == path {
				return p
			}
		}
		return nil
	}()

	if pkg == nil || len(pkg.Errors) > 0 {
		return nil, fmt.Errorf("could not find package %q", path)
	}

	var (
		fct   *types.Func
		scope = pkg.Types.Scope()
	)
	switch {
	case strings.Contains(name, "."):
		idx := strings.Index(name, ".")
		obj := scope.Lookup(name[:idx])
		if obj == nil {
			return nil, fmt.Errorf("could not lookup %q in package %q", name[:idx], path)
		}
		typ, ok := obj.Type().(*types.Named)
		if !ok {
			return nil, fmt.Errorf(
				"object %q in package %q is not a named type (%T)",
				name[:idx], path, obj,
			)
		}
		fct = func() *types.Func {
			for i := 0; i < typ.NumMethods(); i++ {
				m := typ.Method(i)
				if m.Name() == name[idx+1:] {
					return m
				}
			}
			return nil
		}()
		if fct == nil {
			return nil, fmt.Errorf("could not find %q in package %q", name, path)
		}

	default:
		obj := scope.Lookup(name)
		if obj == nil {
			return nil, fmt.Errorf("could not lookup %q in package %q", name, path)
		}
		var ok bool
		fct, ok = obj.(*types.Func)
		if !ok {
			return nil, fmt.Errorf("object %q in package %q is not a func (%T)", name, path, obj)
		}
	}

	if !types.Identical(fct.Type(), f1x.Type()) {
		return nil, fmt.Errorf("invalid function signature for %s", name)
	}

	der := f.Deriv
	if der == "" {
		der = "Deriv" + strings.Replace(f.Name, ".", "_", -1)
	}

	return &generator{w: w, pkg: pkg, fct: fct, der: der}, nil
}

func (g *generator) generate() error {
	fct := func() *ast.FuncDecl {
		for _, f := range g.pkg.Syntax {
			for i := range f.Decls {
				decl, ok := f.Decls[i].(*ast.FuncDecl)
				if !ok {
					continue
				}

				if decl.Name.Name == g.fct.Name() {
					return decl
				}
			}
		}
		return nil
	}()

	var (
		ret     *ast.ReturnStmt
		returns = 0
	)
	ast.Inspect(fct.Body, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.ReturnStmt:
			returns++
			ret = stmt
		}
		return true
	})

	if returns > 1 {
		return fmt.Errorf("can not handle functions with multiple return statements")
	}
	if returns == 0 {
		return fmt.Errorf("could not find a return statement")
	}
	switch len(ret.Results) {
	case 0:
		return fmt.Errorf("naked returns not supported")
	case 1:
		// ok
	default:
		return fmt.Errorf("too many return values")
	}

	args := g.fct.Type().Underlying().(*types.Signature).Params()
	g.printf("func %s(%s float64) float64 {\n",
		g.der,
		args.At(0).Name(),
	)
	g.printf("\treturn ")
	g.expr(ret.Results[0])
	g.printf(".Emag\n")
	g.printf("}\n")

	return g.err
}

func (g *generator) expr(expr ast.Expr) {
	if g.err != nil {
		return
	}

	switch expr := expr.(type) {
	default:
		panic(fmt.Errorf("invalid expr type: %#v (%T)", expr, expr))
	case *ast.BasicLit:
		g.printf("dual.Number{Real:%s}", expr.Value)
	case *ast.Ident:
		g.printf("dual.Number{Real:%s, Emag:1}", expr.Name)
	case *ast.ParenExpr:
		g.printf("(")
		g.expr(expr.X)
		g.printf(")")
	case *ast.UnaryExpr:
		switch expr.Op {
		default:
			panic(fmt.Errorf("invalid binary expression token %v", expr.Op))
		case token.ADD:
			// no op
		case token.SUB:
			g.printf("dual.Mul(dual.Number{Real:1, Emag:1}, ")
			g.expr(expr.X)
			g.printf(")")
		}
	case *ast.BinaryExpr:
		switch expr.Op {
		default:
			panic(fmt.Errorf("invalid binary expression token %v", expr.Op))
		case token.ADD:
			g.printf("dual.Add(")
			g.expr(expr.X)
			g.printf(", ")
			g.expr(expr.Y)
			g.printf(")")
		case token.SUB:
			g.printf("dual.Sub(")
			g.expr(expr.X)
			g.printf(", ")
			g.expr(expr.Y)
			g.printf(")")
		case token.MUL:
			g.printf("dual.Mul(")
			g.expr(expr.X)
			g.printf(", ")
			g.expr(expr.Y)
			g.printf(")")
		case token.QUO:
			g.printf("dual.Mul(")
			g.expr(expr.X)
			g.printf(", ")
			g.printf("dual.Inv(")
			g.expr(expr.Y)
			g.printf("))")
		}

	case *ast.CallExpr:
		g.expr(expr.Fun)
		g.printf("(")
		for i, arg := range expr.Args {
			if i > 0 {
				g.printf(", ")
			}
			g.expr(arg)
		}
		g.printf(")")

	case *ast.SelectorExpr:
		x, ok := expr.X.(*ast.Ident)
		if !ok || x.Name != "math" {
			panic(fmt.Errorf("invalid selector expression %#v", expr))
		}
		switch expr.Sel.Name {
		case "Abs",
			"Acos", "Acosh",
			"Asin", "Asinh",
			"Atan", "Atanh",
			"Cos", "Cosh",
			"Exp", "Log",
			"Pow",
			"Sin", "Sinh",
			"Sqrt",
			"Tan", "Tanh":
			g.printf("dual.%s", expr.Sel.Name)
		case "E", "Pi", "Phi",
			"Sqrt2", "SqrtE", "SqrtPi", "SqrtPhi",
			"Ln2", "Log2E", "Ln10", "Log10E":
			g.printf("dual.Number{Real: math.%s}", expr.Sel.Name)
		default:
			panic(fmt.Errorf("invalid selector expression %#v", expr))
		}
	}
}

func (g *generator) printf(format string, args ...interface{}) {
	fmt.Fprintf(g.w, format, args...)
}

var (
	f1x *types.Func
)

func init() {
	const variadic = false
	f64 := types.NewParam(0, nil, "x", types.Typ[types.Float64])

	sig := types.NewSignature(nil, types.NewTuple(f64), types.NewTuple(f64), variadic)
	f1x = types.NewFunc(0, nil, "dxf", sig)
}
