package evaluator

import (
	"bytes"
	"fmt"
	"strconv"

	m "pasza.org/sr-challenge/model"
)

type CalculatedValue interface {
	isCalculatedValue()
	String() string
}

type intValue int
type floatValue float64
type stringValue string
type boolValue bool
type multiValue []CalculatedValue
type spreadValue []CalculatedValue

func (intValue) isCalculatedValue() {}
func (v intValue) String() string {
	return strconv.Itoa(int(v))
}

func (floatValue) isCalculatedValue() {}
func (v floatValue) String() string {
	return strconv.FormatFloat(float64(v), 'f', 3, 64)
}

func (stringValue) isCalculatedValue() {}
func (v stringValue) String() string {
	return string(v)
}

func (boolValue) isCalculatedValue() {}
func (v boolValue) String() string {
	if v {
		return "true"
	} else {
		return "false"
	}
}

func multipleValueStringer(v []CalculatedValue) string {
	var buff bytes.Buffer
	buff.WriteRune('[')
	for i, value := range v {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(value.String())
	}
	buff.WriteRune(']')
	return string(buff.String())
}

func (multiValue) isCalculatedValue() {}
func (v multiValue) String() string {
	return multipleValueStringer(v)
}

func (spreadValue) isCalculatedValue() {}
func (v spreadValue) String() string {
	return multipleValueStringer(v)
}

type evalCell struct {
	done      bool
	copyCount int
	formula   *m.Expr
	value     CalculatedValue
}

type labelDef struct {
	rowIdx int // row that the label was defined in
	colIdx int
}

// given label name, which column is it (if any)
type labelMap map[string]labelDef

type evalState struct {
	evalCells   [][]evalCell
	csvCells    CSVCells
	labelsOnRow []labelMap
}

type CSVCells [][]m.Cell

func finishLabelCell(cell *evalCell, label string) {
	cell.done = true
	cell.value = stringValue(fmt.Sprintf("!%s", label))
}

func copyMap(m labelMap) labelMap {
	res := make(labelMap)
	for k, v := range m {
		res[k] = v
	}
	return res
}

// first pass to memorize labels, and create output structure
func initState(cells CSVCells) evalState {
	height := len(cells)
	labelsOnRow := make([]labelMap, height)
	evalCells := make([][]evalCell, height)

	currentLabels := make(labelMap)
	for rowIdx, row := range cells {
		newLabelMapNeeded := true
		width := len(row)
		evalRow := make([]evalCell, width)
		for colIdx, cell := range row {
			if labelCell, ok := cell.(m.LabelCell); ok {
				if newLabelMapNeeded {
					currentLabels = copyMap(currentLabels)
					newLabelMapNeeded = false
				}
				label := labelCell.Label
				currentLabels[label] = labelDef{
					rowIdx: rowIdx,
					colIdx: colIdx,
				}
				finishLabelCell(&evalRow[colIdx], label)
			}
		}
		evalCells[rowIdx] = evalRow
		labelsOnRow[rowIdx] = currentLabels
	}
	return evalState{
		evalCells:   evalCells,
		csvCells:    cells,
		labelsOnRow: labelsOnRow,
	}
}

func calcCell(es *evalState, rowIdx int, colIdx int) {
	cell := es.csvCells[rowIdx][colIdx]
	esCell := &es.evalCells[rowIdx][colIdx]
	switch v := cell.(type) {
	case m.IntCell:
		esCell.value = intValue(v.Value)
	case m.FloatCell:
		esCell.value = floatValue(v.Value)
	case m.StringCell:
		esCell.value = stringValue(v.Value)
	case m.EmptyCell:
		esCell.value = stringValue("")
	case m.FormulaCell:
		calcFormulaCell(es, rowIdx, colIdx, &v)
	default:
		// label cells should be calculated beforehand
		panic("Cannot evaluate unknown cell type")
	}
	esCell.done = true
}

func calcInfixOp(es *evalState, v m.InfixOp, rowIdx int, colIdx int) CalculatedValue {
	lhs, rhs := calcExpr(es, &v.Lhs, rowIdx, colIdx), calcExpr(es, &v.Rhs, rowIdx, colIdx)
	switch v.Op {
	case m.MUL:
		return calcMul(lhs, rhs)
	case m.DIV:
		return calcDiv(lhs, rhs)
	case m.ADD:
		return calcAdd(lhs, rhs)
	case m.SUB:
		return calcSub(lhs, rhs)
	default:
		panic("Cannot evaluate unknown infix operation")
	}
}

func calcExpr(es *evalState, expr *m.Expr, rowIdx int, colIdx int) CalculatedValue {
	switch v := (*expr).(type) {
	case m.IntLit:
		return intValue(v)
	case m.FloatLit:
		return floatValue(v)
	case m.StringLit:
		return stringValue(v)
	case m.InfixOp:
		return calcInfixOp(es, v, rowIdx, colIdx)
	case m.FunCall:
		return calcFunCall(es, v, rowIdx, colIdx)
	case m.CellRef:
		return calcCellRef(es, v, rowIdx, colIdx)
	case m.CopyAbove:
		return calcCopyAbove(es, v, rowIdx, colIdx)
	case m.CopyColumnAbove:
		return calcCopyColumnAbove(es, v, rowIdx, colIdx)
	case m.CopyLastInColumn:
		return calcCopyLastInColumn(es, v, rowIdx, colIdx)
	case m.LabelRelativeRowRef:
		return calcLabelRelativeRowRef(es, v, rowIdx, colIdx)
	default:
		panic("Cannot evaluate unknown expression")
	}
}

func colNameToIdx(colName string) int {
	name := int(colName[0]) // fixme: make it safer
	return name - 'A'
}

func calcLabelRelativeRowRef(es *evalState, v m.LabelRelativeRowRef, rowIdx, colIdx int) CalculatedValue {
	label := v.Label
	relativeRow := v.RelativeRow
	labelAnchor, found := es.labelsOnRow[rowIdx][label]
	if !found {
		panic("Invalid relative row reference")
	}
	targetColIdx := labelAnchor.colIdx
	targetRowIdx := labelAnchor.rowIdx + relativeRow
	return getTargetValue(es, targetRowIdx, targetColIdx)
}

func calcCopyLastInColumn(es *evalState, v m.CopyLastInColumn, rowIdx, colIdx int) CalculatedValue {
	targetColIdx := colNameToIdx(v.Col)
	// search upwards for available column
	targetRowIdx := rowIdx - 1
	for targetRowIdx >= 0 {
		if len(es.csvCells[targetRowIdx]) > targetColIdx {
			return getTargetValue(es, targetRowIdx, targetColIdx)
		}
		targetRowIdx--
	}
	panic("Copy last available column calculation failed")
}

func getTargetValue(es *evalState, targetRowIdx int, targetColIdx int) CalculatedValue {
	target := &es.evalCells[targetRowIdx][targetColIdx]
	if !target.done {
		calcCell(es, targetRowIdx, targetColIdx)
	}
	return target.value
}

// e.g. C^
func calcCopyColumnAbove(es *evalState, v m.CopyColumnAbove, rowIdx, colIdx int) CalculatedValue {
	targetColIdx := colNameToIdx(v.Col)
	return getTargetValue(es, rowIdx-1, targetColIdx)
}

func calcCellRef(es *evalState, v m.CellRef, rowIdx, colIdx int) CalculatedValue {
	targetColIdx := colNameToIdx(v.Col)
	targetRowIdx := v.Row - 1 // we use 0-based indexing

	return getTargetValue(es, targetRowIdx, targetColIdx)
}

func calcCopyAbove(es *evalState, v m.CopyAbove, rowIdx, colIdx int) CalculatedValue {
	if rowIdx <= 0 {
		panic("Attempted to copy above in row 0")
	}
	above := &es.evalCells[rowIdx-1][colIdx]
	if !above.done {
		calcCell(es, rowIdx-1, colIdx)
	}
	ec := &es.evalCells[rowIdx][colIdx]
	ec.copyCount = above.copyCount + 1
	ec.formula = above.formula
	if ec.formula != nil {
		ec.value = calcExpr(es, ec.formula, rowIdx, colIdx)
	} else {
		ec.value = above.value
	}
	return ec.value
}

func calcFormulaCell(es *evalState, rowIdx, colIdx int, cell *m.FormulaCell) {
	esCell := &es.evalCells[rowIdx][colIdx]
	esCell.formula = &cell.Formula
	esCell.value = calcExpr(es, &cell.Formula, rowIdx, colIdx)
}

func calculateAll(es *evalState, cells CSVCells) {
	for rowIdx, row := range cells {
		for colIdx := range row {
			if !es.evalCells[rowIdx][colIdx].done {
				calcCell(es, rowIdx, colIdx)
			}
		}
	}
}

func Evaluate(cells CSVCells) [][]CalculatedValue {
	evalState := initState(cells)
	calculateAll(&evalState, cells)

	// rewrite just calculated values and return
	res := make([][]CalculatedValue, len(evalState.evalCells))
	for rowIdx, row := range evalState.evalCells {
		resRow := make([]CalculatedValue, len(row))
		for colIdx, cell := range row {
			resRow[colIdx] = cell.value
		}
		res[rowIdx] = resRow
	}
	return res
}
