package columns

import (
	"sort"
	"strings"
)

// Sort the lines base on one (or more) columns
//
// Example: Sort(-1) will sort the 1st column descending;
// Sort(2,3) will sort on column 2 and, if necessary, column 3
//
// In a column with mixed datatypes (strings & numerical), strings will be sorted after all numerical values
func (cw *ColumnsWriter) Sort(columns ...int) {
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

	switch x := a.value.(type) {
	case string:
		switch y := b.value.(type) {
		case string:
			comp := strings.Compare(x, y)
			if comp == 0 {
				return compareEqual
			}
			if (asc && comp > 0) || (!asc && comp < 0) {
				return compareSwap
			}
			return compareKeep

		case int, int64, float64:
			return compareSwap // swap strings after numerical
		}

	case int64:
		switch y := b.value.(type) {
		case int:
			if x == int64(y) {
				return compareEqual
			}
			if (asc && x > int64(y)) || (!asc && x < int64(y)) {
				return compareSwap
			}
			return compareKeep
		case int64:
			if x == y {
				return compareEqual
			}
			if (asc && x > y) || (!asc && x < y) {
				return compareSwap
			}
			return compareKeep
		case float64:
			if (asc && float64(x) > y) || (!asc && float64(x) < y) {
				return compareSwap
			}
			return compareKeep
		case string:
			return compareKeep // keep numerical before string
		}

	case int:
		switch y := b.value.(type) {
		case int:
			if x == y {
				return compareEqual
			}
			if (asc && x > y) || (!asc && x < y) {
				return compareSwap
			}
			return compareKeep
		case int64:
			if int64(x) == y {
				return compareEqual
			}
			if (asc && int64(x) > y) || (!asc && int64(x) < y) {
				return compareSwap
			}
			return compareKeep
		case float64:
			if (asc && float64(x) > y) || (!asc && float64(x) < y) {
				return compareSwap
			}
			return compareKeep
		case string:
			return compareKeep // keep numerical before string
		}

	case float64:
		switch y := b.value.(type) {
		case float64:
			if x == y {
				return compareEqual
			}
			if (asc && x > y) || (!asc && x < y) {
				return compareSwap
			}
			return compareKeep
		case int:
			if x == float64(y) {
				return compareEqual
			}
			if (asc && x > float64(y)) || (!asc && x < float64(y)) {
				return compareSwap
			}
			return compareKeep
		case int64:
			if x == float64(y) {
				return compareEqual
			}
			if (asc && x > float64(y)) || (!asc && x < float64(y)) {
				return compareSwap
			}
			return compareKeep
		case string:
			return compareKeep // keep numerical before string
		}

	default:
		//fmt.Printf("compare(%v,%v): unknown type of A\n", a.value, b.value)
	}

	return compareKeep
}
