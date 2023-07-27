package model

type Expr interface {
	isExpr()
}

type FunCall struct {
	Name   string
	Params []Expr
}

type BinaryOperator int

const (
	MUL BinaryOperator = 0
	DIV BinaryOperator = 1
	ADD BinaryOperator = 2
	SUB BinaryOperator = 3
)

type IntLit int
type FloatLit float64
type StringLit string

type CellRef struct {
	Col string
	Row int
}

type CopyAbove struct{}
type CopyLastInColumn struct {
	Col string
}
type CopyColumnAbove struct {
	Col string
}

type LabelRelativeRowRef struct {
	Label       string
	RelativeRow int
}

type InfixOp struct {
	Lhs Expr
	Rhs Expr
	Op  BinaryOperator
}

type NoResult struct{}

// Make sure all the expression variants implement Expr
func (IntLit) isExpr()    {}
func (FloatLit) isExpr()  {}
func (StringLit) isExpr() {}
func (InfixOp) isExpr()   {}
func (FunCall) isExpr()   {}

func (CellRef) isExpr()             {}
func (CopyAbove) isExpr()           {}
func (CopyLastInColumn) isExpr()    {}
func (LabelRelativeRowRef) isExpr() {}
func (CopyColumnAbove) isExpr()     {}
