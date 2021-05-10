package columns

import "strings"

// CellData contains the cell value and it's assigned styling
type CellData struct {
	value interface{}
	style *Style
}

// Cell create a cell that can be colored, prefixed or suffixed
func Cell(value interface{}) *CellData {
	return &CellData{
		value: value,
	}
}

// Style applies a style (from the NewStyle-function) to a specific cell
func (cell *CellData) Style(style *Style) *CellData {
	cell.style = style
	return cell
}

func (cell *CellData) prefix(style *Style) string {
	if cell.style != nil && cell.style.prefix != nil {
		return *cell.style.prefix
	}
	if style != nil && style.prefix != nil {
		return *style.prefix
	}
	return ""
}

func (cell *CellData) suffix(style *Style) string {
	if cell.style != nil && cell.style.suffix != nil {
		return *cell.style.suffix
	}
	if style != nil && style.suffix != nil {
		return *style.suffix
	}
	return ""
}

func (cell *CellData) isEmpty() bool {
	return cell == nil || cell.value == nil
}

func (cw *Writer) writeCell(c *CellData, col *column) {

	var txt, prefix, suffix string
	var sizeI, sizeF, lenPrefix, lenSuffix int

	if c == nil || c.value == nil {
		txt = strings.Repeat(space, col.outerSize())
		lenPrefix = col.sizePrefix
		lenSuffix = col.sizeSuffix
	} else {
		txt, _, sizeI, sizeF = cw.format(c)
		if col.sizeDot > 0 && sizeF == 0 { // add space to integer values where other rows have a decimal separator
			txt += " "
		}
		if sizeI > 0 && col.sizeI > 0 {
			txt = strings.Repeat(space, col.sizeI-sizeI) + txt
		}
		if (sizeI > 0 || sizeF > 0) && col.sizeF > 0 {
			txt = txt + strings.Repeat(space, col.sizeF-sizeF)
		}

		prefix = c.prefix(col.style)
		suffix = c.suffix(col.style)

		lenPrefix = len([]rune(prefix))
		lenSuffix = len([]rune(suffix))
	}

	if cw.useColor {
		_, _ = cw.bufwr.WriteString(c.beginStyle(col))
	}

	if col.sizePrefix > 0 {
		_, _ = cw.bufwr.WriteString(prefix)
		_, _ = cw.bufwr.WriteString(strings.Repeat(space, col.sizePrefix-lenPrefix))
	}

	_, _ = cw.bufwr.WriteString(pad(txt, col.innerSize(), col.align, ' '))

	if col.sizeSuffix > 0 {
		_, _ = cw.bufwr.WriteString(suffix)
		_, _ = cw.bufwr.WriteString(strings.Repeat(space, col.sizeSuffix-lenSuffix))
	}

	if cw.useColor {
		_, _ = cw.bufwr.WriteString(c.endStyle(col))
	}
}
