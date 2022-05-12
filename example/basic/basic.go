package main

import (
	"os"

	"github.com/ninlil/columns"
)

func main() {
	cw := columns.New(os.Stdout, "| ^ | < | > |")
	cw.Headers("Position", "Planet", "Relative radius")
	cw.HeaderSeparator = true
	cw.Write(1, "Mercury", 0.3825)
	cw.Write(2, "Venus", 0.9488)
	cw.Write(3, "Earth", 1)
	cw.Write(4, "Mars", 0.53260)
	cw.Write(5, "Jupiter", 11.209)
	cw.Write(6, "Saturn", 9.449)
	cw.Write(7, "Uranus", 4.007)
	cw.Write(8, "Neptune", 3.883)
	cw.Flush()
}
