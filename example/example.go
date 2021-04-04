package main

import (
	"fmt"
	"os"

	"github.com/ninlil/columns"
)

func main() {
	styled()
	fmt.Println("=================")
	headNtail()
}

func styled() {
	cw := columns.New(os.Stdout, "| ^ | < | ^ | > | > |")
	cw.HeaderSeparator = true
	cw.Headers("Position", "Planet", "Relative radius", "Orbital period", "Avg. temp (C)")

	cw.Write(1, "Mercury", 0.3825, 87.969, columns.Cell(167).Color(columns.ColorRed, false).Suffix("°C"))
	cw.Write(2, "Venus", 0.9488, 224.701, columns.Cell(457).Color(columns.ColorRed, false).Suffix("°C"))
	cw.Write(3, "Earth", 1, 365.256363, columns.Cell(13.85).Suffix("°C"))
	cw.Write(4, "Mars", 0.53260, 686.971, columns.Cell(-46).Color(columns.ColorBlue, false))
	cw.Write(5, "Jupiter", 11.209, 4_332.59, columns.Cell(-121).Color(columns.ColorBlue, true))
	cw.Write(6, "Saturn", 9.449, 10_759.22, columns.Cell(-139).Color(columns.ColorBlue, true).Suffix("°C"))
	cw.Write(7, "Uranus", 4.007, 30_688.5, columns.Cell(-197).Color(columns.ColorBlue, true))
	cw.Write(8, "Neptune", 3.883, 60_182, columns.Cell(-200).Color(columns.ColorBlue, true))

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
