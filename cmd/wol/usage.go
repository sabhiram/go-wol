package main

import (
	"fmt"
	"strings"

	"github.com/sabhiram/go-colorize"
)

// List of strings which contain allowed commands
var ValidCommands = []struct {
	name, description string
}{
	{`wake`, `wakes up a machine by mac address or alias`},
	{`list`, `lists all mac addresses and their aliases`},
	{`alias`, `stores an alias to a mac address`},
	{`remove`, `removes an alias or a mac address`},
}

// List of options which wol supports
var ValidOptions = []struct {
	short, long, description string
}{
	{`v`, `version`, `prints the application version`},
	{`h`, `help`, `prints this help menu`},
	{`p`, `port`, `udp port to send bcast packet to`},
	{`b`, `bcast`, `broadcast IP to send packet to`},
	{`i`, `interface`, `outbound interface to broadcast using`},
}

// Usage string for wol
var UsageString = `Usage:

    To wake up a machine:
        <cyan>wol</cyan> [<options>] <yellow>wake</yellow> <mac address | alias> <optional interface>

    To store an alias:
        <cyan>wol</cyan> [<options>] <yellow>alias</yellow> <alias> <mac address> <optional interface>

    To view aliases:
        <cyan>wol</cyan> [<options>] <yellow>list</yellow>

    To delete aliases:
        <cyan>wol</cyan> [<options>] <yellow>remove</yellow> <alias>

    The following MAC addresses are valid and will match:
    01-23-45-56-67-89, 89:AB:CD:EF:00:12, 89:ab:cd:ef:00:12

    The following MAC addresses are not (yet) valid:
    1-2-3-4-5-6, 01 23 45 56 67 89

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

// Returns true if the ValidCommands struct contains an entry with the
// input string "s"
func isValidCommand(s string) bool {
	for _, c := range ValidCommands {
		if strings.ToLower(s) == c.name {
			return true
		}
	}
	return false
}
