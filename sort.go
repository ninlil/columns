package columns

import (
	"math"
	"sort"
	"strings"
)

// Sort the lines base on one (or more) columns
//
// Example: Sort(-1) will sort the 1st column descending;
// Sort(2,3) will sort on column 2 and, if necessary, column 3
//
// In a column with mixed datatypes (strings & numerical), numerical values are grouped 1st, strings 2nd, nil are always last
// regardless of sorting ascending or descending
func (cw *Writer) Sort(columns ...int) {
	var sorter = func(i, j int) bool {
		for _, c := range columns {
			var asc = c > 0
			var colIndex = int(abs(int64(c))) - 1
			if colIndex < len(cw.columns) {
				a := cw.data[i][colIndex]
				b := cw.data[j][colIndex]
				switch compare(a, b, asc) {
				case compareSwap:
					return false
				case compareKeep:
					return true
				}
			}
		}
		return false
	}
	sort.SliceStable(cw.data, sorter)
}

type swap int

const (
	compareSwap  swap = 1
	compareEqual swap = 2
	compareKeep  swap = 3
)

func compareValue(a, b *CellData, asc bool) int {
	switch x := a.value.(type) {
	case string:
		switch y := b.value.(type) {
		case string:
			return strings.Compare(x, y)
		default:
			if asc {
				return 1 // swap
			}
			return -1
		}
	}

	switch b.value.(type) {
	case string:
		if asc {
			return -1 // swap
		}
		return 1
	}

	n, ok1 := getNum(a.value)
	m, ok2 := getNum(b.value)
	if ok1 && ok2 {
		diff := n - m
		if math.Abs(diff) < 1e-9 {
			return 0
		}
		if diff > 0 {
			return 1
		}
		return -1
	}
	if !ok2 {
		return 1
	}
	return 0
}

func compare(a, b *CellData, asc bool) swap {
	if a.isEmpty() && b.isEmpty() {
		return compareKeep // both empty -> dont swap
	}
	if a.isEmpty() || b.isEmpty() {
		if a.isEmpty() {
			return compareSwap // one empty -> swap if 1st is empty (promote actual values)
		}
		return compareKeep
	}

	comp := compareValue(a, b, asc)
	if comp == 0 {
		return compareEqual
	}
	if (asc && comp > 0) || (!asc && comp < 0) {
		return compareSwap
	}
	return compareKeep
}
