package columns

import "strings"

type CellData struct {
	value  interface{}
	color  Color
	bright bool
	prefix string
	suffix string
}

type Color byte

const (
	ColorDefault Color = 0
	ColorBlack   Color = '0'
	ColorRed     Color = '1'
	ColorGreen   Color = '2'
	ColorYellow  Color = '3'
	ColorBlue    Color = '4'
	ColorMagenta Color = '5'
	ColorCyan    Color = '6'
	ColorWhite   Color = '7'
)

// Create a cell that can be colored, prefixed or suffixed
func Cell(value interface{}) *CellData {
	return &CellData{
		value:  value,
		color:  ColorDefault,
		bright: false,
	}
}

// Prefix the cell value with a text
func (cell *CellData) Prefix(text string) *CellData {
	cell.prefix = text
	return cell
}

// Suffix the cell value with a text
func (cell *CellData) Suffix(text string) *CellData {
	cell.suffix = text
	return cell
}

// Color the cell value
func (cell *CellData) Color(color Color, bright bool) *CellData {
	if cell != nil {
		cell.color = color
		cell.bright = bright
	}
	return cell
}

func (cell *CellData) isEmpty() bool {
	return cell == nil || cell.value == nil
}

func (cw *ColumnsWriter) writeCell(c *CellData, col *column) {

	if c == nil || c.value == nil {
		_, _ = cw.bufwr.WriteString(strings.Repeat(space, col.outerSize()))
		return
	}

	var colors = []byte{27, '[', '0', ';', '3', '7', 'm'}
	txt, _, sizeI, sizeF := cw.format(c)
	if col.sizeDot > 0 && sizeF == 0 { // add space to integer values where other rows have a decimal separator
		txt += " "
	}
	if sizeI > 0 && col.sizeI > 0 {
		txt = strings.Repeat(space, col.sizeI-sizeI) + txt
	}
	if (sizeI > 0 || sizeF > 0) && col.sizeF > 0 {
		txt = txt + strings.Repeat(space, col.sizeF-sizeF)
	}
	if c.color != ColorDefault && cw.useColor {
		if c.bright {
			colors[2] = '1'
		} else {
			colors[2] = '0'
		}
		colors[5] = byte(c.color)
		_, _ = cw.bufwr.Write(colors)
	}

	lenPrefix := len([]rune(c.prefix))
	lenSuffix := len([]rune(c.suffix))

	if col.sizePrefix > 0 {
		_, _ = cw.bufwr.WriteString(c.prefix)
		_, _ = cw.bufwr.WriteString(strings.Repeat(space, col.sizePrefix-lenPrefix))
	}

	_, _ = cw.bufwr.WriteString(pad(txt, col.innerSize(), col.align, ' '))

	if col.sizeSuffix > 0 {
		_, _ = cw.bufwr.WriteString(c.suffix)
		_, _ = cw.bufwr.WriteString(strings.Repeat(space, col.sizeSuffix-lenSuffix))
	}

	if c.color != ColorDefault && cw.useColor {
		colors[2] = '0'
		colors[5] = '7'
		_, _ = cw.bufwr.Write(colors)
	}
}
