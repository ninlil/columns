package columns

import (
	"bufio"
	"fmt"
	"strings"
)

// Flush writes the completed columns to the output
func (cw *ColumnsWriter) Flush() {

	if len(cw.aggOrder) > 0 {
		for _, col := range cw.columns {
			for _, agg := range col.aggregations {
				col.ensureSize(cw, Cell(agg.Result()), col.style)
			}
		}
	}

	for _, c := range cw.columns {
		if c.sizeI > 0 || c.sizeF > 0 {
			size := c.sizeI + c.sizeDot + c.sizeF
			if c.sizeValue < size {
				c.sizeValue = size
			}
		}
	}

	cw.dump()

	var sep []string

	cw.bufwr = bufio.NewWriter(cw.writer)
	if len(cw.headers) > 0 {
		cw.writeStrings(cw.headers, "\n")

		if cw.HeaderSeparator {
			sep = make([]string, cw.n)
			for i, col := range cw.columns {
				sep[i] = strings.Repeat("-", col.outerSize())
			}
			cw.writeStrings(sep, "\n")
		}
	}
	if cw.head < 0 {
		cw.head = len(cw.data)
	}
	cw.tail = len(cw.data) - cw.tail

	cutmsg := false
	for i, row := range cw.data {
		if i < cw.head || i >= cw.tail {
			cw.writeCells(row, "\n")
		} else {
			if !cutmsg {
				_, _ = cw.bufwr.WriteString(fmt.Sprintf("--- cut %d lines ---\n", cw.tail-cw.head))
				cutmsg = true
			}
		}
	}

	if len(cw.aggOrder) > 0 {
		if len(sep) > 0 {
			cw.writeStrings(sep, "\n")
		}
		for _, aggName := range cw.aggOrder {
			aggline := make([]*CellData, cw.n)
			for i, col := range cw.columns {
				if agg, ok := col.aggregations[aggName]; ok {
					aggline[i] = Cell(agg.Result())
				}
			}
			cw.writeCells(aggline, " ", aggName, "\n")
		}
	}

	_ = cw.bufwr.Flush()
}

func (cw *ColumnsWriter) writeCells(data []*CellData, suffix ...string) {
	for i := 0; i < cw.n; i++ {
		col := cw.columns[i]

		if len(cw.spacers[i]) > 0 {
			_, _ = cw.bufwr.WriteString(cw.spacers[i])
		}
		if i < len(data) {
			cw.writeCell(data[i], col)
		}
	}
	_, _ = cw.bufwr.WriteString(cw.spacers[cw.n])
	_, _ = cw.bufwr.WriteString(strings.Join(suffix, ""))
}

func (cw *ColumnsWriter) writeStrings(data []string, suffix ...string) {

	for i := 0; i < cw.n; i++ {
		col := cw.columns[i]

		if len(cw.spacers[i]) > 0 {
			_, _ = cw.bufwr.WriteString(cw.spacers[i])
		}
		if i < len(data) {
			_, _ = cw.bufwr.WriteString(pad(data[i], col.outerSize(), col.align, ' '))
		}
	}
	_, _ = cw.bufwr.WriteString(cw.spacers[cw.n])
	_, _ = cw.bufwr.WriteString(strings.Join(suffix, ""))
}

func pad(txt string, size int, align alignment, ch rune) string {
	l := len([]rune(txt))
	if l >= size {
		return txt
	}

	after := 0
	before := size - l
	switch align {
	case AlignLeft:
		after = before
		before = 0
	case AlignMiddle:
		after = before / 2
		before -= after
	}
	return strings.Repeat(string(ch), before) + txt + strings.Repeat(string(ch), after)
}
