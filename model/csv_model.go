package model

type Cell interface {
	isCell()
}

type IntCell struct {
	Value int
}

type FloatCell struct {
	Value float64
}

type StringCell struct {
	Value string
}

type LabelCell struct {
	Label string
}

type FormulaCell struct {
	Formula Expr
}

type Formula interface {
	isFormula()
}

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

func (StringCell) isCell()  {}
func (IntCell) isCell()     {}
func (FloatCell) isCell()   {}
func (LabelCell) isCell()   {}
func (FormulaCell) isCell() {}
func (EmptyCell) isCell()   {}
