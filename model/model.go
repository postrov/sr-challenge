package model

type Cell interface {
	isCell()
}

type StringCell struct {
	Value string
}

type FormulaCell struct {
	Formula Formula
}

type Formula interface {
	isFormula()
}

type DummyFormula struct {
	RawValue string
}

func (DummyFormula) isFormula() {}

type Row struct {
	Len   int
	Cells []Cell
}

type RowGroup struct {
	Len    int      // number of rows in this row group
	RowLen int      // length of each individual row
	Labels []string // Column labels
}

type EmptyCell struct{}

func (StringCell) isCell() {}

func (FormulaCell) isCell() {}

func (EmptyCell) isCell() {}

type Expression interface {
	isExpression()
}

type FunCall struct {
	name   string
	params []Expression
}

type BinOp int

const (
	MUL BinOp = 0
	DIV       = 1
	ADD       = 2
	SUB       = 3
)

func (op BinOp) String() string {
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

type InfixOp struct {
	lhs Expression
	rhs Expression
	op  BinOp
}

// expr: float | int | cellref | unary op
