# Colorize

[![Build Status](https://travis-ci.org/sabhiram/go-colorize.svg)](https://travis-ci.org/sabhiram/go-colorize) [![Coverage Status](https://coveralls.io/repos/sabhiram/go-colorize/badge.png?branch=master)](https://coveralls.io/r/sabhiram/go-colorize?branch=master)

A `Go` library to fetch colorized ascii text

For more information regarding ASCII Escape Codes, and Colors see [this link](http://en.wikipedia.org/wiki/ANSI_escape_code).

Table of supported colors for reference (from the above wiki entry):

![](https://raw.githubusercontent.com/sabhiram/public-images/master/colorize/ascii_color_table.png)

## Install

```shell
go get github.com/sabhiram/go-colorize
```

## (Sample) Usage

```shell
mkdir color_me && cd color_me
touch color_me.go
```

*color_me/color_me.go*:
```go
package main

import (
    "fmt"

    "github.com/sabhiram/go-colorize"
)

func main() {
    fmt.Println("Hello - " + colorize.ColorString("black",   "black"))
    fmt.Println("Hello - " + colorize.ColorString("red",     "red"))
    fmt.Println("Hello - " + colorize.ColorString("green",   "green"))
    fmt.Println("Hello - " + colorize.ColorString("yellow",  "yellow"))
    fmt.Println("Hello - " + colorize.ColorString("blue",    "blue"))
    fmt.Println("Hello - " + colorize.ColorString("magenta", "magenta"))
    fmt.Println("Hello - " + colorize.ColorString("cyan",    "cyan"))
    fmt.Println("Hello - " + colorize.ColorString("white",   "white"))

    fmt.Println(colorize.Colorize(`

<red>This text will be red</red> and this is default...

            <blue>This is blue!</blue>

<red>0</red><green>1</green><yellow>3</yellow><blue>4</blue><magenta>5</magenta><cyan>6</cyan><white>7</white>

`))
}
```

#### Install and run:

```shell
cd color_me
go install
color_me
```

#### Outputs:

![](https://raw.githubusercontent.com/sabhiram/public-images/master/colorize/colorize_sample.png)

