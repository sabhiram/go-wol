// This file contains version specific, and usage information for
// the wol application
package main

import (
    "fmt"
    "github.com/sabhiram/go-colorize"
)

// Version represents the current Semantic Version of this application
const Version = "1.0.0"

// List of strings which contain allowed commands
var ValidCommands = [] struct {
    name, description string
} {
    { `wake`,     `wakes up a machine by mac address or alias` },
    { `list`,     `lists all mac addresses and their aliases`  },
    { `alias`,    `stores an alias to a mac address`           },
    { `remove`,   `removes an alias or a mac address`          },
}

// List of options which wol supports
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

    To wake up a machine:
        <cyan>wol</cyan> [<options>] <yellow>wake</yellow> <mac address | alias>

    To store an alias:
        <cyan>wol</cyan> [<options>] <yellow>alias</yellow> <alias> <mac address>

    To view aliases:
        <cyan>wol</cyan> [<options>] <yellow>list</yellow>

    To delete aliases:
        <cyan>wol</cyan> [<options>] <yellow>remove</yellow> <mac address | alias>

    The following MAC addresses are valid and will match:
    01-23-45-56-67-89, 89:AB:CD:EF:00:12, 89:ab:cd:ef:00:12

Commands:
%s
Options:
%s
Version:
    <white>%s</white>

`

// Build a command string from the above valid ones
func getAllCommands() string {
    commands := ""
    for _, c := range ValidCommands {
        commands += fmt.Sprintf("    <yellow>%-16s</yellow> %s\n", c.name, c.description)
    }
    return commands
}

// Build an option string from the above valid ones
func getAllOptions() string {
    options := ""
    for _, o := range ValidOptions {
        options += fmt.Sprintf("    <yellow>-%s --%-8s</yellow>    %s\n", o.short, o.long, o.description)
    }
    return options
}

// Returns the Usage string for this application
func getAppUsageString() string {
    return colorize.Colorize(fmt.Sprintf(UsageString, getAllCommands(), getAllOptions(), Version))
}

