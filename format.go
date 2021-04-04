package columns

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	space = " "
)

func (cw *ColumnsWriter) format(cell *CellData) (txt string, size int, sizeI int, sizeF int) {
	sizeI = 0
	sizeF = 0
	switch v := cell.value.(type) {
	case string:
		txt = v
		size = len([]rune(txt))

	case int:
		txt, size, sizeI, sizeF = cw.formatNumeric(strconv.FormatInt(abs(int64(v)), 10), v < 0)

	case int64:
		txt, size, sizeI, sizeF = cw.formatNumeric(strconv.FormatInt(abs(v), 10), v < 0)

	case float64:
		txt, size, sizeI, sizeF = cw.formatNumeric(fmt.Sprintf("%v", math.Abs(v)), v < 0)

	default:
		txt = fmt.Sprintf("%v", v)
		size = len([]rune(txt))
	}

	return txt, size, sizeI, sizeF
}

func abs(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

func (cw *ColumnsWriter) formatNumeric(input string, neg bool) (txt string, size int, sizeI int, sizeF int) {
	parts := strings.Split(input, ".")

	var txtI, txtF string
	txtI = separate(parts[0], -len(parts[0]), cw.ThousandSeparator)
	if neg {
		txtI = "- " + txtI
	}
	if len(parts) > 1 {
		txtF = parts[1]
	}

	sizeI = len([]rune(txtI))
	sizeF = len([]rune(txtF))

	txt = txtI
	size = sizeI
	if sizeF > 0 {
		size += 1 + sizeF
		txt += string(cw.DecimalSeparator) + txtF
	}
	return txt, size, sizeI, sizeF
}

func separate(txt string, i int, sep rune) string {
	result := make([]rune, 0, 10)
	lastWasSeparator := false
	for _, ch := range txt {
		lastWasSeparator = false
		result = append(result, ch)
		i++
		if i%3 == 0 && sep != rune(0) {
			lastWasSeparator = true
			result = append(result, sep)
		}
	}
	if lastWasSeparator {
		result = result[:len(result)-1]
	}
	return string(result)
}
