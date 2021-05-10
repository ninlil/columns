package columns

import (
	"github.com/ninlil/ansi"
)

const colorFunc ansi.Style = -1

// ColorFunc should return a ansi.Style value to set a style/color depending on data-value
//
// Please note: all numerical values are converted to float64 before called
type ColorFunc func(v interface{}) (color ansi.Style, ok bool)

func getValueForColorFunc(v interface{}) interface{} {
	n, ok := getNum(v)
	if ok {
		return n
	}
	return v
}

// Style of a cell or column
type Style struct {
	color   ansi.Style
	colorFn ColorFunc

	prefix *string
	suffix *string
}

// NewStyle creates a new style
func NewStyle() *Style {
	return &Style{}
}

// Color assigns a ansi.Style to the cell/columns style
func (s *Style) Color(color ansi.Style) *Style {
	s.color = color
	return s
}

// ColorFunc sets a formatting function to style depending on value
func (s *Style) ColorFunc(fn ColorFunc) *Style {
	s.color = colorFunc
	s.colorFn = fn
	return s
}

// Prefix adds a prefix to the data-value for the cell/column
func (s *Style) Prefix(text string) *Style {
	s.prefix = &text
	return s
}

// Suffix adds a suffix to the data-value for the cell/column
func (s *Style) Suffix(text string) *Style {
	s.suffix = &text
	return s
}

func (c *CellData) getStyle(col *column) *Style {
	if c != nil && c.style != nil {
		return c.style
	}
	return col.style
}

func (c *CellData) beginStyle(col *column) string {
	var style = c.getStyle(col)
	if style == nil || style.color == ansi.Default {
		return ""
	}

	var color ansi.Style = ansi.White
	var ok bool

	if style.color == colorFunc {
		color, ok = style.colorFn(getValueForColorFunc(c.value))
		if !ok {
			return ""
		}
	} else {
		color = style.color
	}

	return color.String()
}

func (c *CellData) endStyle(col *column) string {
	var style = c.getStyle(col)
	if style == nil || style.color == ansi.Default {
		return ""
	}

	return ansi.Default.String()
}
