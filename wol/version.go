// This file contains version specific, and usage information for
// the wol application
package main

import (
    "fmt"
    "github.com/sabhiram/go-colorize"
)

// Version represents the current Semantic Version of this application
const Version = "1.0.0"

// List of options which chloe supports
var ValidOptions = [] struct {
    short, long, description string
} {
    { `v`, `version`, `prints the application version`   },
    { `h`, `help`,    `prints this help menu`            },
    { `p`, `port`,    `udp port to send bcast packet to` },
    { `b`, `bcast`,   `broadcast IP to send packet to`   },
}

// Usage string for wol
var UsageString = `Usage:

    <cyan>wol</cyan> [<options>] [mac address]

Options:

%s

Specifying MAC addresses:
-------------------------
The following MAC addresses are valid and will match:
01-23-45-56-67-89, 89:AB:CD:EF:00:12, 89:ab:cd:ef:00:12

The supported delimiters include ":" and "-"

Version:

    <white>%s</white>

`
// Returns a tuple of commands and options which we support
func getAllOptions() string {
    options := ""
    for _, o := range ValidOptions {
        options += fmt.Sprintf("    <yellow>-%s --%-8s</yellow>    %s\n", o.short, o.long, o.description)
    }
    return options
}

// Returns the Usage string for this application
func getAppUsageString() string {
    options := getAllOptions()
    return colorize.Colorize(fmt.Sprintf(UsageString, options, Version))
}

