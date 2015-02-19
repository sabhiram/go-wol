package main

import (
	"fmt"
	"os"

	"errors"

	"github.com/sabhiram/go-colorize"
	wol "github.com/sabhiram/go-wol"

	"github.com/jessevdk/go-flags"
)

var (
	// Define holders for the cli arguments we wish to parse
	Options struct {
		Version             bool   `short:"v" long:"version"`
		Help                bool   `short:"h" long:"help"`
		BroadcastInterface  string `short:"i" long:"interface" default:""`
		BroadcastIP         string `short:"b" long:"bcast" default:"255.255.255.255"`
		UDPPort             string `short:"p" long:"port" default:"9"`
	}
)

// Run the alias command
func runAliasCommand(args []string, aliases map[string]MacIface) error {
	if len(args) >= 2 {
		var eth string
		if len(args) > 2 {
			eth = args[2]
		}
		// TODO: Validate mac address
		alias, mac := args[0], args[1]
		aliases[alias] = MacIface{Mac: mac, Iface: eth}
		return flushUserAliases(aliases)
	}
	return errors.New("alias command requires a <name> and a <mac>")
}

// Run the list command
func runListCommand(args []string, aliases map[string]MacIface) error {
	if len(aliases) == 0 {
		fmt.Printf("No aliases found! Add one with \"wol alias <name> <mac>\"\n")
	} else {
		for alias, mi := range aliases {
			if mi.Iface == "" {
				fmt.Printf("    %s - %s\n", alias, mi.Mac)
			} else {
				fmt.Printf("    %s - %s %s\n", alias, mi.Mac, mi.Iface)
			}
		}
	}
	return nil
}

// Run the remove command
func runRemoveCommand(args []string, aliases map[string]MacIface) error {
	if len(args) > 0 {
		alias := args[0]
		delete(aliases, alias)
		return flushUserAliases(aliases)
	}
	return errors.New("remove command requires a <name> of an alias")
}

// Run the wake command
func runWakeCommand(args []string, aliases map[string]MacIface) error {
	if len(args) <= 0 {
		return errors.New("No mac address specified to wake command")
	}

	// bcastInterface can be "eth0", "eth1", etc.. An empty string implies
	// that we use the default interface when sending the UDP packet (nil)
	bcastInterface := ""
	macAddr := args[0]

	// First we need to see if this macAddr is actually an alias, if it is
	// we set the eth interface based on the stored item, and set the macAddr
	// based on the alias entry's mac for this alias
	if val, ok := aliases[macAddr]; ok {
		macAddr = val.Mac
		bcastInterface = val.Iface
	}

	// If the command line specified an interface, we override to use that one
	// regardless of what is set in the alias map
	if len(args) > 1 {
		bcastInterface = args[1]
	}

	err := wol.SendMagicPacket(macAddr, Options.BroadcastIP + ":" + Options.UDPPort, bcastInterface)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return errors.New("Unable to send magic packet")
	}

	fmt.Printf("Magic packet sent successfully to %s\n", macAddr)
	return nil
}

// Run one of the supported commands
func runCommand(cmd string, args []string, aliases map[string]MacIface) error {
	switch cmd {

	case "alias":
		return runAliasCommand(args, aliases)

	case "list":
		return runListCommand(args, aliases)

	case "remove":
		return runRemoveCommand(args, aliases)

	case "wake":
		return runWakeCommand(args, aliases)

	default:
		panic("Invalid command passed to runCommand")

	}
	return nil
}

// Helper function to dump the usage and print an error if specified,
// it also returns the exit code requested to the function (saves me a line)
func printUsageGetExitCode(s string, e int) int {
	if len(s) > 0 {
		fmt.Printf(s)
	}
	fmt.Printf(getAppUsageString())
	return e
}

// Main entry point for binary
func main() {
	var args []string

	// Load the list of aliases from ~/.config/go-wol/aliases
	aliases, err := loadUserAliases()
	if err != nil {
		panic("Unable to load user aliases! Exiting...")
	}

	// Parse arguments which might get passed to "wol"
	parser := flags.NewParser(&Options, flags.Default & ^flags.HelpFlag)
	args, err = parser.Parse()

	exitCode := 0
	switch {

	// Parse Error, print usage
	case err != nil:
		exitCode = printUsageGetExitCode("", 1)

	// No arguments, or help requested, print usage
	case len(os.Args) == 1 || Options.Help:
		exitCode = printUsageGetExitCode("", 0)

	// "--version" requested
	case Options.Version:
		fmt.Printf("%s\n", Version)

	// Make sure we are being asked to run a command
	case len(args) == 0:
		exitCode = printUsageGetExitCode("No command specified, see usage:\n", 1)

	// All other cases go here
	case true:
		cmd, args := args[0], args[1:]
		if isValidCommand(cmd) {
			err = runCommand(cmd, args, aliases)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
				exitCode = 1
			}
		} else {
			exitCode = printUsageGetExitCode(
				fmt.Sprintf("Unknown command %s, see usage:\n", colorize.ColorString(cmd, "red")), 1)
		}

	}
	os.Exit(exitCode)
}
