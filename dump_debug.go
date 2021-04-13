// +build debug

package columns

import (
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
			fillSize += col.outerSize() - col.sizeI - col.sizeF - col.sizePrefix - col.sizeSuffix
			fillBefore := strings.Repeat("_", fillSize)
			if col.align == AlignLeft { // not working on AlignMiddle
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

	fmt.Print(a, "\n", b, "\n", c, "\n")
}
