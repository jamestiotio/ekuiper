package ast

import (
	"reflect"
)

type Visitor interface {
	Visit(Node) bool
}

func Walk(v Visitor, node Node) {
	if node == nil || reflect.ValueOf(node).IsNil() {
		return
	}

	if !v.Visit(node) {
		return
	}

	switch n := node.(type) {
	case *SelectStatement:
		Walk(v, n.Fields)
		Walk(v, n.Sources)
		Walk(v, n.Joins)
		Walk(v, n.Condition)
		Walk(v, n.Dimensions)
		Walk(v, n.Having)
		Walk(v, n.SortFields)

	case Fields:
		for _, f := range n {
			Walk(v, &f)
		}

	case *Field:
		Walk(v, n.Expr)
		if fr, ok := n.Expr.(*FieldRef); ok && fr.IsAlias() {
			Walk(v, fr.Expression)
		}

	case Sources:
		for _, s := range n {
			Walk(v, s)
		}

	//case *Table:

	case Joins:
		for _, s := range n {
			Walk(v, &s)
		}

	case *Join:
		Walk(v, n.Expr)

	case Dimensions:
		Walk(v, n.GetWindow())
		for _, dimension := range n.GetGroups() {
			Walk(v, dimension.Expr)
		}

	case *Window:
		Walk(v, n.Length)
		Walk(v, n.Interval)
		Walk(v, n.Filter)

	case SortFields:
		for _, sf := range n {
			Walk(v, &sf)
		}

	//case *SortField:

	case *BinaryExpr:
		Walk(v, n.LHS)
		Walk(v, n.RHS)

	case *Call:
		for _, expr := range n.Args {
			Walk(v, expr)
		}

	case *ParenExpr:
		Walk(v, n.Expr)

	case *CaseExpr:
		Walk(v, n.Value)
		for _, w := range n.WhenClauses {
			Walk(v, w)
		}
		Walk(v, n.ElseClause)

	case *WhenClause:
		Walk(v, n.Expr)
		Walk(v, n.Result)
	}
}

// WalkFunc traverses a node hierarchy in depth-first order.
func WalkFunc(node Node, fn func(Node) bool) {
	Walk(walkFuncVisitor(fn), node)
}

type walkFuncVisitor func(Node) bool

func (fn walkFuncVisitor) Visit(n Node) bool { return fn(n) }