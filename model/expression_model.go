package model

import (
	"fmt"
)

type Expr interface {
	isExpr()
}

type FunCall struct {
	Name   string
	Params []Expr
}

func (FunCall) isExpr() {}

type BinaryOperator int

func (op BinaryOperator) String() string {
	switch op {
	case MUL:
		return "*"
	case DIV:
		return "/"
	case ADD:
		return "+"
	case SUB:
		return "-"
	default:
		return "?"
	}
}

const (
	MUL BinaryOperator = 0
	DIV BinaryOperator = 1
	ADD BinaryOperator = 2
	SUB BinaryOperator = 3
)

type IntLit int

func (IntLit) isExpr() {}

type InfixOp struct {
	Lhs Expr
	Rhs Expr
	Op  BinaryOperator
}

func formatInfixOperand(e Expr) string {
	switch e.(type) {
	case InfixOp:
		return fmt.Sprintf("(%v)", e)
	default:
		return fmt.Sprint(e)
	}
}

func (op InfixOp) String() string {
	lhs, rhs := formatInfixOperand(op.Lhs), formatInfixOperand(op.Rhs)
	return fmt.Sprintf("%v %v %v", lhs, op.Op, rhs)
}

func (InfixOp) isExpr() {}

type NoResult struct{}
