package main

import (
	"fmt"
	"os"

	"github.com/ninlil/ansi"
	"github.com/ninlil/columns"
)

func main() {
	styled()
	fmt.Println("=================")
	headNtail()
}

func tempFunc(v interface{}) (c ansi.Style, ok bool) {
	if v == nil {
		return (ansi.White).Background(), true
	}

	if n, ok := v.(float64); ok {
		if n > 0 {
			return ansi.Red | ansi.Bright, true
		}
	}
	return ansi.Blue | ansi.Bright, true
}

func styled() {
	cw := columns.New(os.Stdout, "| ^ | < | ^ | > | > |")
	cw.HeaderSeparator = true
	cw.Headers("Position", "Planet", "Relative radius", "Orbital period", "Avg. temp")
	cw.Footer(1, columns.Sum(1), columns.Avg(1))
	cw.Footer(3, columns.Sum(1), columns.Avg(1))
	cw.Footer(4, columns.Sum(1), columns.Avg(1))
	cw.Footer(5, columns.Sum(1), columns.Avg(1))

	temp := columns.NewStyle().ColorFunc(tempFunc).Suffix("°C")
	earth := columns.NewStyle().Color(ansi.Green).Suffix("°C")
	cw.Style(5, temp)

	cw.Write(1, "Mercury", 0.3825, 87.969, nil) //columns.Cell(167).Style(temp))
	cw.Write(2, "Venus", 0.9488, 224.701, 457)
	cw.Write(3, "Earth", 1, 365.256363, columns.Cell(13.85).Style(earth))
	cw.Write(4, "Mars", 0.53260, 686.971, -46)
	cw.Write(5, "Jupiter", 11.209, 4_332.59, -121)
	cw.Write(6, "Saturn", 9.449, 10_759.22, -139)
	cw.Write(7, "Uranus", 4.007, 30_688.5, -197)
	cw.Write(8, "Neptune", 3.883, 60_182, -200)

	cw.Sort(-1)

	cw.Flush()
}

func headNtail() {
	cw := columns.New(os.Stdout, "<")

	for i := 1; i <= 10; i++ {
		cw.Write(i)
	}
	cw.Head(3)
	cw.Tail(2)
	cw.Flush()
}
