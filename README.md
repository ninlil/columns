# columns

And alternative to the golang text/tabwriter-package,
with support for colors, sorting, head and tail.

## Getting Started

### Installing

Just do a `go get` on this package (and my 'ansi' for styling) and your good to go.
```
go get github.com/ninlil/columns
go get github.com/ninlil/ansi
```

## Features

* Individual column alignment
* Auto-align numerical values on the decimal-point
* Sorting your output before printing
* Head and Tail to only show the start and/or end of your data
* Colorizeable output (only to default terminal; i.e to os.Stdout that is a CharDevice)

## Examples

See [examples](./example) for more details
```
go run example/basic.go
go run example/example.go
```

Basic example:

```go
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
```

```
| Position | Planet  | Relative radius |
| -------- | ------- | --------------- |
|     1    | Mercury |          0.3825 |
|     2    | Venus   |          0.9488 |
|     3    | Earth   |          1      |
|     4    | Mars    |          0.5326 |
|     5    | Jupiter |         11.209  |
|     6    | Saturn  |          9.449  |
|     7    | Uranus  |          4.007  |
|     8    | Neptune |          3.883  |
```

## Formatting

Cells and entire columns can be formatted

```go
style := columns.NewStyle().Color(ansi.Red).Suffix("°C")

// apply the style to the entire 3rd column
cw.Style(3, style)

// apply to the cell (overriding the column styling)
cw.Write(3, "Earth", columns.Cell(1).Style(style))
```

### Conditional formatting

The following example will color all values above 0 as Red and 0 or below as Blue, and any `nil` values using a white background.

```go
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
...
temp := columns.NewStyle().ColorFunc(tempFunc).Suffix("°C")
```

### Note 1
Any prefix or suffix set on the column will still be printed unless the cell-styling actually contains a prefix/suffix on its own.

### Note 2
Make sure to apply column styling before calling any Write to ensure prefix and suffix will fit with your values

### Note 3 - Conditional values
All numerical values are converted to `float64` before the `ColorFunc` is called.

## Footers

You can add footer containing sums of your data
```go
// will add a 'Sum' and a 'Avg' footer to the list (with a decimal precision of 1)
cw.Footer(3, columns.Sum(1), columns.Avg(1))
```

## Sort, Head & Tail
```go
cw.Sort(-1, 4) // will sort descending on the 1st column, then ascending on column 4

cw.Head(5)     // will only print the first 5 rows
cw.Tail(10)    // will only print the last 10 rows
```

`Head` and `Tail` can be used at the same time.

If any lines are excluded then a line indicating how many rows where cut will be printed.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/ninlil/columns/tags). 

## Code usage

* [Go Report Card](https://goreportcard.com/report/github.com/ninlil/columns)
* [pkg.go.dev](https://pkg.go.dev/github.com/ninlil/columns)

## Authors

* **ninlil** - *Initial work* - [ninlil](https://github.com/ninlil)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
