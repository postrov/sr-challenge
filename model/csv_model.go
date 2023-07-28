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

func (StringCell) isCell()  {}
func (IntCell) isCell()     {}
func (FloatCell) isCell()   {}
func (LabelCell) isCell()   {}
func (FormulaCell) isCell() {}
