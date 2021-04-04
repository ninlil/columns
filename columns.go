package columns

import (
	"bufio"
	"io"
	"os"
)

type ColumnsWriter struct {
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
}

type column struct {
	sizeHeader int // Size of the header
	sizeValue  int // Max size of all values
	sizeI      int // Max size of Integer-part (not including decimal separator)
	sizeDot    int // 0 or 1 depending on if sizeF > 0
	sizeF      int // Max size of Decimal-part (not including decimal separator)
	sizePrefix int // Mas size of all prefixes
	sizeSuffix int // Mas size of all suffixes
	align      alignment
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

const (
	AlignLeft   alignment = '<'
	AlignRight  alignment = '>'
	AlignMiddle alignment = '^'
)

// New creates a ColumnWriter based on the 'format'
//
// <, > and ^ are columns aligned left, right and centered (respectively)
// any other characters are padding between,
// spaces as padding are not added automatically
func New(writer io.Writer, format string) *ColumnsWriter {
	useColor := false

	if writer == os.Stdout {
		if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
			useColor = true
		}
	}

	cw := ColumnsWriter{
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
	cw.spacers = append(cw.spacers, string(spacer)+"\n")
	cw.n = len(cw.columns)

	return &cw
}

// Head limits the output to the 'n' first lines (can be combined with Tail)
func (cw *ColumnsWriter) Head(n int) {
	cw.head = n
	if cw.tail < 0 {
		cw.tail = 0
	}
}

// Tail limits the output to the 'n' last lines (can be combined with Head)
func (cw *ColumnsWriter) Tail(n int) {
	cw.tail = n
	if cw.head < 0 {
		cw.head = 0
	}
}

// Separator sets the 'thousand' and 'decimal' separators
func (cw *ColumnsWriter) Separator(thousand, decimal rune) {
	cw.ThousandSeparator = thousand
	cw.DecimalSeparator = decimal
}

// Headers sets text-headers to each column
//
// More headers than columns defined in 'New' will be ignored
func (cw *ColumnsWriter) Headers(titles ...string) {
	cw.headers = make([]string, cw.n)
	for i, hdr := range titles {
		if i < len(cw.headers) {
			cw.headers[i] = hdr
			cw.columns[i].sizeHeader = len([]rune(hdr))
		}
	}
}

// Write a line/row to the ColumnWriter
//
// Sortable datatypes are string, int, int64, and float64
// other datatypes will be printed using fmt.Sprintf("%v")
//
// More values than columns defined in 'New' will be ignored
func (cw *ColumnsWriter) Write(data ...interface{}) {
	row := make([]*CellData, cw.n)
	for i, o := range data {
		if i < cw.n {
			var cell *CellData

			if c, ok := o.(*CellData); ok {
				cell = c
			} else {
				cell = Cell(o)
			}

			if cell.prefix != "" {
				l := len([]rune(cell.prefix))
				if cw.columns[i].sizePrefix < l {
					cw.columns[i].sizePrefix = l
				}
			}
			if cell.suffix != "" {
				l := len([]rune(cell.suffix))
				if cw.columns[i].sizeSuffix < l {
					cw.columns[i].sizeSuffix = l
				}
			}

			_, size, sizeI, sizeF := cw.format(cell)
			row[i] = cell

			if cw.columns[i].sizeValue < size {
				cw.columns[i].sizeValue = size
			}
			if cw.columns[i].sizeI < sizeI {
				cw.columns[i].sizeI = sizeI
			}
			if cw.columns[i].sizeF < sizeF {
				cw.columns[i].sizeF = sizeF
				cw.columns[i].sizeDot = 1
			}
		}
	}
	cw.data = append(cw.data, row)
}
