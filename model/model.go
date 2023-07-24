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
