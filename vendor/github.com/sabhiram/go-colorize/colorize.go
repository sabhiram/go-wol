// Package colorize returns an ascii colorized string version of an input string
//
// Here is a table of ASCII to color values:
//
//     Intensity   0       1      2       3       4       5       6       7
//     Normal      Black   Red    Green   Yellow  Blue    Magenta Cyan    White
//     Bright      Black   Red    Green   Yellow  Blue    Magenta Cyan    White
//
// Table sourced from: http://en.wikipedia.org/wiki/ANSI_escape_code
//
package colorize

import (
	"fmt"
	"regexp"
	"strings"
)

// colorToValueMap can convert the string values of an input ASCII
// color to its appropriate integer counterpart. See table above for
// mapping information
var colorToValueMap = map[string]int{
	"black":   0,
	"red":     1,
	"green":   2,
	"yellow":  3,
	"blue":    4,
	"magenta": 5,
	"cyan":    6,
	"white":   7,
}

// ColorString is a function which returns the "input" string after surrounding
// it by the appropriate ASCII escape sequence for the requested "color".
// If the "color" input is not a member of the map, we simply return
// the input without any change
func ColorString(input, color string) string {
	color = strings.ToLower(color)
	if colorIndex, valid := colorToValueMap[color]; valid {
		return fmt.Sprintf("\033[3%dm%s\033[0m", colorIndex, input)
	}
	return input
}

// Colorize is a function whichtakes in a string embedded with
// "<color>...</color>" tags. Any un-matched tags will be ignored and not
// stripped from the input string
func Colorize(input string) string {
	for color, index := range colorToValueMap {
		var (
			searchString  = fmt.Sprintf("(<%s>)(.*?)(</%s>)", color, color)
			replaceString = fmt.Sprintf("\033[3%dm$2\033[0m", index)
			match         = regexp.MustCompile(searchString)
		)
		input = match.ReplaceAllString(input, replaceString)
	}
	return input
}
