# columns

And alternative to the golang text/tabwriter-package,
with support for colors, sorting, head and tail.

## Getting Started

### Installing

Just do a `go get` on this package and your good to go.
```
go get github.com/ninlil/columns
```

## Features

* Individual column alignment
* Auto-align numerical values on the decimal-point
* Sorting your output before printing
* Head and Tail to only show the start and/or end of your data
* Colorizeable output (only to default terminal)

## Examples

See [examples](./example) for more details
```
go run example/basic.go
go run example/example.go
```

Basic example:

```
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

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/ninlil/columns/tags). 

## Authors

* **ninlil** - *Initial work* - [ninlil](https://github.com/ninlil)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
