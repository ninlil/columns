package columns

import (
	"bufio"
	"io"
	"os"
)

// Writer is the main class capable of printing columns with headers, footers, dynamic styling and more
type Writer struct {
	HeaderSeparator   bool // Add a header-separator between header and values
	ThousandSeparator rune // Change (or remove) the automatic thousand-separator (default ' ', disable when 0)
	DecimalSeparator  rune // Change the decimal separator (default '.')

	writer  io.Writer
	bufwr   *bufio.Writer
	n       int
	columns []*column
	spacers []string
	headers []string
	data    [][]*CellData

	head int
	tail int

	useColor bool

	aggOrder []string
}

type column struct {
	sizeHeader   int // Size of the header
	sizeValue    int // Max size of all values
	sizeI        int // Max size of Integer-part (not including decimal separator)
	sizeDot      int // 0 or 1 depending on if sizeF > 0
	sizeF        int // Max size of Decimal-part (not including decimal separator)
	sizePrefix   int // Max size of all prefixes
	sizeSuffix   int // Max size of all suffixes
	align        alignment
	style        *Style
	aggregations map[string]Aggregation
}

func (col *column) outerSize() int {
	vSize := col.sizeValue + col.sizePrefix + col.sizeSuffix
	if vSize > col.sizeHeader {
		return vSize
	}
	return col.sizeHeader
}

func (col *column) innerSize() int {
	return col.outerSize() - col.sizePrefix - col.sizeSuffix
}

type alignment rune

// Alignment symbols
const (
	AlignLeft   alignment = '<'
	AlignRight  alignment = '>'
	AlignMiddle alignment = '^'
)

// New creates a Writer based on the 'format'
//
// <, > and ^ are columns aligned left, right and centered (respectively)
// any other characters are padding between,
// spaces as padding are not added automatically
func New(writer io.Writer, format string) *Writer {
	useColor := false

	if writer == os.Stdout {
		if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
			useColor = true
		}
	}

	cw := Writer{
		writer:            writer,
		ThousandSeparator: ' ',
		DecimalSeparator:  '.',
		useColor:          useColor,
		head:              -1,
		tail:              -1,
	}
	var spacer []rune
	for _, ch := range format {
		switch ch {
		case rune(AlignLeft), rune(AlignMiddle), rune(AlignRight):
			cw.spacers = append(cw.spacers, string(spacer))
			spacer = spacer[:0]
			cw.columns = append(cw.columns, &column{align: alignment(ch)})

		default:
			spacer = append(spacer, ch)
		}
	}
	cw.spacers = append(cw.spacers, string(spacer))
	cw.n = len(cw.columns)

	return &cw
}

// Head limits the output to the 'n' first lines (can be combined with Tail)
func (cw *Writer) Head(n int) {
	cw.head = n
	if cw.tail < 0 {
		cw.tail = 0
	}
}

// Tail limits the output to the 'n' last lines (can be combined with Head)
func (cw *Writer) Tail(n int) {
	cw.tail = n
	if cw.head < 0 {
		cw.head = 0
	}
}

// Separator sets the 'thousand' and 'decimal' separators
func (cw *Writer) Separator(thousand, decimal rune) {
	cw.ThousandSeparator = thousand
	cw.DecimalSeparator = decimal
}

// Headers sets text-headers to each column
//
// More headers than columns defined in 'New' will be ignored
func (cw *Writer) Headers(titles ...string) {
	cw.headers = make([]string, cw.n)
	for i, hdr := range titles {
		if i < len(cw.headers) {
			cw.headers[i] = hdr
			cw.columns[i].sizeHeader = len([]rune(hdr))
		}
	}
}

// Footer creates a footer for column 'i' (1-based) with the supplied aggregations
func (cw *Writer) Footer(i int, aggrs ...Aggregation) {
	i--
	if i >= 0 && i < cw.n {
		if cw.columns[i].aggregations == nil {
			cw.columns[i].aggregations = make(map[string]Aggregation)
		}

		for _, agg := range aggrs {
			name := agg.Name()
			found := false
			for _, n := range cw.aggOrder {
				if n == name {
					found = true
				}
			}
			if !found {
				cw.aggOrder = append(cw.aggOrder, name)
			}
			cw.columns[i].aggregations[name] = agg
		}
	}
}

// Style applies a single style to an entire column (1-based)
func (cw *Writer) Style(i int, style *Style) {
	i--
	if i >= 0 && i < cw.n {
		cw.columns[i].style = style
	}
}

// Write a line/row to the Writer
//
// Sortable datatypes are string, int, int64, and float64
// other datatypes will be printed using fmt.Sprintf("%v")
//
// More values than columns defined in 'New' will be ignored
func (cw *Writer) Write(data ...interface{}) {
	row := make([]*CellData, cw.n)
	for i, o := range data {
		if i < cw.n {
			var cell *CellData

			if c, ok := o.(*CellData); ok {
				cell = c
			} else {
				cell = Cell(o)
			}

			for _, agg := range cw.columns[i].aggregations {
				_ = agg.AddValue(cell.value)
			}

			cw.columns[i].ensureSize(cw, cell, cw.columns[i].style)
			row[i] = cell
		}
	}
	cw.data = append(cw.data, row)
}

func (col *column) ensureSize(cw *Writer, cell *CellData, style *Style) {

	if txt := cell.prefix(style); txt != "" {
		l := len([]rune(txt))
		if col.sizePrefix < l {
			col.sizePrefix = l
		}
	}
	if txt := cell.suffix(style); txt != "" {
		l := len([]rune(txt))
		if col.sizeSuffix < l {
			col.sizeSuffix = l
		}
	}

	_, size, sizeI, sizeF := cw.format(cell)

	if col.sizeValue < size {
		col.sizeValue = size
	}
	if col.sizeI < sizeI {
		col.sizeI = sizeI
	}
	if col.sizeF < sizeF {
		col.sizeF = sizeF
		col.sizeDot = 1
	}
}
