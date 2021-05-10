package columns

import (
	"fmt"
	"math"
)

// Aggregation is the interface for describing a calculated value in a footer
type Aggregation interface {
	AddValue(v interface{}) error
	Name() string
	Result() float64
}

// ErrInvalidType error for saying it's a value we can't "aggregate"
var ErrInvalidType = fmt.Errorf("not a numerical value")

func round(n float64, precision int) float64 {
	decimals := math.Pow10(precision)
	return math.Round(n*decimals) / decimals
}

func getNum(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int64:
		return float64(n), true
	case int32:
		return float64(n), true
	case int16:
		return float64(n), true
	case int8:
		return float64(n), true
	case int:
		return float64(n), true
	case uint64:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint:
		return float64(n), true
	case bool:
		if n {
			return 1, true
		}
		return 0, true
	}
	return 0, false
}

// SUM
type aggSum struct {
	precision int
	value     float64
}

// Sum creates an 'Aggregation' that 'sums' all valid values
func Sum(prec int) Aggregation {
	return &aggSum{
		precision: prec,
	}
}

func (sum *aggSum) Name() string {
	return "Sum"
}

func (sum *aggSum) Result() float64 {
	return round(sum.value, sum.precision)
}

func (sum *aggSum) AddValue(v interface{}) error {
	if n, ok := getNum(v); ok {
		sum.value += n
	}
	return ErrInvalidType
}

// AVG
type aggAvg struct {
	precision int
	total     float64
	count     int
}

// Avg creates an 'Aggregation' that 'averages' all valid values
func Avg(prec int) Aggregation {
	return &aggAvg{
		precision: prec,
	}
}

func (avg *aggAvg) Name() string {
	return "Average"
}

func (avg *aggAvg) Result() float64 {
	return round(avg.total/float64(avg.count), avg.precision)
}

func (avg *aggAvg) AddValue(v interface{}) error {
	if n, ok := getNum(v); ok {
		avg.total += n
		avg.count++
	}
	return ErrInvalidType
}
