package columns

import (
	"bufio"
	"fmt"
	"strings"
)

func (cw *ColumnsWriter) dump() {
	var a, b, c string
	for i, col := range cw.columns {
		a += fmt.Sprint(cw.spacers[i], strings.Repeat("O", col.outerSize()))
		b += fmt.Sprint(cw.spacers[i], strings.Repeat("+", col.sizePrefix), strings.Repeat("I", col.innerSize()), strings.Repeat("-", col.sizeSuffix))

		// //fmt.Printf("> Column %d: size:%d sizeI:%d sizeF:%d\n", i, c.size, c.sizeI, c.sizeF)
		if col.sizeI > 0 {
			var fillSize int
			var fillAfter string
			if col.sizeF > 0 {
				fillSize--
			}
			fillSize += col.sizeValue - col.sizeI - col.sizeF
			fillBefore := strings.Repeat(space, fillSize)
			if col.align == AlignLeft {
				fillAfter = fillBefore
				fillBefore = ""
			}
			c += fmt.Sprint(cw.spacers[i],
				strings.Repeat("+", col.sizePrefix),
				fillBefore,
				strings.Repeat("N", col.sizeI),
				strings.Repeat(".", col.sizeDot),
				strings.Repeat("D", col.sizeF),
				fillAfter,
				strings.Repeat("-", col.sizeSuffix))
		} else {
			c += fmt.Sprint(cw.spacers[i],
				strings.Repeat("+", col.sizePrefix),
				strings.Repeat("T", col.sizeValue),
				strings.Repeat("-", col.sizeSuffix))
		}
	}
	a += cw.spacers[len(cw.columns)]
	b += cw.spacers[len(cw.columns)]
	c += cw.spacers[len(cw.columns)]

	fmt.Print(a, b, c)
}

func (cw *ColumnsWriter) Flush() {
	for _, c := range cw.columns {
		if c.sizeI > 0 || c.sizeF > 0 {
			size := c.sizeI + c.sizeDot + c.sizeF
			if c.sizeValue < size {
				c.sizeValue = size
			}
		}
	}

	//cw.dump()

	cw.bufwr = bufio.NewWriter(cw.writer)
	if len(cw.headers) > 0 {
		cw.writeStrings(cw.headers)

		if cw.HeaderSeparator {
			sep := make([]string, cw.n)
			for i, col := range cw.columns {
				sep[i] = strings.Repeat("-", col.outerSize())
			}
			cw.writeStrings(sep)
		}
	}
	if cw.head < 0 {
		cw.head = len(cw.data)
	}
	cw.tail = len(cw.data) - cw.tail

	cutmsg := false
	for i, row := range cw.data {
		if i < cw.head || i >= cw.tail {
			cw.writeCells(row)
		} else {
			if !cutmsg {
				_, _ = cw.bufwr.WriteString(fmt.Sprintf("--- cut %d lines ---\n", cw.tail-cw.head))
				cutmsg = true
			}
		}
	}

	_ = cw.bufwr.Flush()
}

func (cw *ColumnsWriter) writeCells(data []*CellData) {
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
}

func (cw *ColumnsWriter) writeStrings(data []string) {

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
