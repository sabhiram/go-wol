package main

////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"

	"github.com/sabhiram/go-wol/wol"

	"github.com/sabhiram/go-colorize"
)

////////////////////////////////////////////////////////////////////////////////

var (
	validCommands = []struct {
		name, description string
	}{
		{`wake`, `wakes up a machine by mac address or alias`},
		{`list`, `lists all mac addresses and their aliases`},
		{`alias`, `stores an alias to a mac address`},
		{`remove`, `removes an alias or a mac address`},
	}

	validOptions = []struct {
		short, long, description string
	}{
		{`v`, `version`, `prints the application version`},
		{`h`, `help`, `prints this help menu`},
		{`d`, `db-dir`, `directory to store alias db`},
		{`a`, `db-name`, `bold db file name (default "bolt.db")`},
		{`c`, `no-color`, `disables ANSI color`},
		{`p`, `port`, `udp port to send bcast packet to`},
		{`b`, `bcast`, `broadcast IP to send packet to`},
		{`i`, `interface`, `outbound interface to broadcast using`},
	}

	usageString = `Usage:

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
)

////////////////////////////////////////////////////////////////////////////////

// Build a command string from the above valid ones.
func getAllCommands() string {
	commands := ""
	for _, c := range validCommands {
		commands += fmt.Sprintf("    <yellow>%-16s</yellow> %s\n", c.name, c.description)
	}
	return commands
}

// Build an option string from the above valid ones.
func getAllOptions() string {
	options := ""
	for _, o := range validOptions {
		options += fmt.Sprintf("    <yellow>-%s --%-10s</yellow>    %s\n", o.short, o.long, o.description)
	}
	return options
}

// Returns the Usage string for this application.
func getAppUsageString() string {
	return colorize.Colorize(fmt.Sprintf(usageString, getAllCommands(), getAllOptions(), wol.Version))
}
