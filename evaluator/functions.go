package evaluator

import (
	"bytes"
	"strconv"
	"strings"

	m "pasza.org/sr-challenge/model"
)

var supportedFunctions map[string](func(*evalState, []CalculatedValue, int, int) CalculatedValue)

func init() {
	supportedFunctions = map[string](func(*evalState, []CalculatedValue, int, int) CalculatedValue){
		"sum":     sum,
		"bte":     bte,
		"text":    text,
		"incFrom": incFrom,
		"concat":  concat,
		"split":   split,
		"spread":  spread,
	}
}

func incFrom(es *evalState, args []CalculatedValue, rowIdx int, colIdx int) CalculatedValue {
	ec := es.evalCells[rowIdx][colIdx]
	if len(args) != 1 {
		panic("Function incFrom() requires exactly one argument")
	}
	v, ok := args[0].(intValue)
	if !ok {
		panic("Function incFrom() expects int argument")
	}

	return intValue(v + intValue(ec.copyCount))
}

func sum(es *evalState, args []CalculatedValue, rowIdx int, colIdx int) CalculatedValue {
	res := 0.0
	for _, arg := range args {
		switch v := arg.(type) {
		case intValue:
			res += float64(v)
		case floatValue:
			res += float64(v)
		case stringValue:
			f, err := strconv.ParseFloat(string(v), 64)
			if err != nil {
				panic("Couldn't convert sum() argument to a number")
			}
			res += f
		default:
			panic("Unknown argument type passed to sum()")
		}
	}
	return floatValue(res)
}

func bte(es *evalState, args []CalculatedValue, rowIdx int, colIdx int) CalculatedValue {
	if len(args) != 2 {
		panic("Function bte() expects exacltly two arguments")
	}

	return bte_(args[0], args[1])
}

func text(es *evalState, args []CalculatedValue, rowIdx int, colIdx int) CalculatedValue {
	if len(args) != 1 {
		panic("Function text() expects exactly one argument")
	}
	return stringValue(args[0].String())
}

func spread(es *evalState, args []CalculatedValue, rowIdx int, colIdx int) CalculatedValue {
	if len(args) != 1 {
		panic("Function spread() expects exactly one argument")
	}
	mv, ok := args[0].(multiValue)
	if !ok {
		panic("Function spread() expectes argument of multivalue type")
	}
	return spreadValue(mv)
}

func split(es *evalState, args []CalculatedValue, rowIdx int, colIdx int) CalculatedValue {
	if len(args) != 2 {
		panic("Function split() expects exactly two arguments")
	}
	s, ok := args[0].(stringValue)
	if !ok {
		panic("Function split() first argument must be a string value")
	}
	sep, ok := args[1].(stringValue)
	if !ok {
		panic("Function split() second argument must be a string value")
	}
	parts := strings.Split(string(s), string(sep))
	res := make([]CalculatedValue, len(parts))
	for i, part := range parts {
		res[i] = stringValue(part)
	}
	return multiValue(res)
}

func concat(es *evalState, args []CalculatedValue, rowIdx int, colIdx int) CalculatedValue {
	var buff bytes.Buffer
	for _, v := range args {
		buff.WriteString(v.String())
	}
	return stringValue(buff.String())
}

func calcFunCall(es *evalState, v m.FunCall, rowIdx int, colIdx int) CalculatedValue {
	args := make([]CalculatedValue, 0)
	for _, expr := range v.Params {
		arg := calcExpr(es, &expr, rowIdx, colIdx)
		if spreadArg, ok := arg.(spreadValue); ok {
			for _, a := range spreadArg {
				args = append(args, a)
			}
		} else {
			args = append(args, calcExpr(es, &expr, rowIdx, colIdx))
		}
	}

	f, ok := supportedFunctions[v.Name]
	if !ok {
		panic("Function not found: " + v.Name)
	}

	return f(es, args, rowIdx, colIdx)
}
