package main

////////////////////////////////////////////////////////////////////////////////

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/sabhiram/go-colorize"
	"github.com/sabhiram/go-wol/wol"

	flags "github.com/jessevdk/go-flags"
)

////////////////////////////////////////////////////////////////////////////////

const (
	defaultDBDir = "/.config/go-wol"
)

var (
	// Define holders for the cli arguments we wish to parse.
	cliFlags struct {
		Version            bool   `short:"v" long:"version"`
		DBDir              string `short:"d" long:"db-dir" default:""`
		DBName             string `short:"a" long:"db-name" default:"bolt.db"`
		Help               bool   `short:"h" long:"help"`
		NoColor            bool   `short:"n" long:"no-color"`
		BroadcastInterface string `short:"i" long:"interface" default:""`
		BroadcastIP        string `short:"b" long:"bcast" default:"255.255.255.255"`
		UDPPort            string `short:"p" long:"port" default:"9"`
	}
)

////////////////////////////////////////////////////////////////////////////////

// ipFromInterface returns a `*net.UDPAddr` from a network interface name.
func ipFromInterface(iface string) (*net.UDPAddr, error) {
	ief, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, err
	}

	addrs, err := ief.Addrs()
	if err == nil && len(addrs) <= 0 {
		err = fmt.Errorf("no address associated with interface %s", iface)
	}
	if err != nil {
		return nil, err
	}

	// Validate that one of the addrs is a valid network IP address.
	for _, addr := range addrs {
		switch ip := addr.(type) {
		case *net.IPNet:
			if !ip.IP.IsLoopback() && ip.IP.To4() != nil {
				return &net.UDPAddr{
					IP: ip.IP,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("no address associated with interface %s", iface)
}

////////////////////////////////////////////////////////////////////////////////

// Run the alias command.
func aliasCmd(args []string, aliases *Aliases) error {
	if len(args) >= 2 {
		var eth string
		if len(args) > 2 {
			eth = args[2]
		}
		// TODO: Validate mac address
		alias, mac := args[0], args[1]
		return aliases.Add(alias, mac, eth)
	}
	return errors.New("alias command requires a <name> and a <mac>")
}

// Run the list command.
func listCmd(args []string, aliases *Aliases) error {
	mp, err := aliases.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get list of aliases: %v\n", err)
		return err
	}
	if len(mp) == 0 {
		fmt.Printf("No aliases found! Add one with \"wol alias <name> <mac>\"\n")
	} else {
		for alias, mi := range mp {
			fmt.Printf("    %s - %s %s\n", alias, mi.Mac, mi.Iface)
		}
	}
	return nil
}

// Run the remove command.
func removeCmd(args []string, aliases *Aliases) error {
	if len(args) > 0 {
		alias := args[0]
		return aliases.Del(alias)
	}
	return errors.New("remove command requires a <name> of an alias")
}

// Run the wake command.
func wakeCmd(args []string, aliases *Aliases) error {
	if len(args) <= 0 {
		return errors.New("No mac address specified to wake command")
	}

	// bcastInterface can be "eth0", "eth1", etc.. An empty string implies
	// that we use the default interface when sending the UDP packet (nil).
	bcastInterface := ""
	macAddr := args[0]

	// First we need to see if this macAddr is actually an alias, if it is:
	// we set the eth interface based on the stored item, and set the macAddr
	// based on the alias of the entry.
	mi, err := aliases.Get(macAddr)
	if err == nil {
		macAddr = mi.Mac
		bcastInterface = mi.Iface
	}

	// Always use the interface specified in the command line, if it exists.
	if cliFlags.BroadcastInterface != "" {
		bcastInterface = cliFlags.BroadcastInterface
	}

	// Populate the local address in the event that the broadcast interface has
	// been set.
	var localAddr *net.UDPAddr
	if bcastInterface != "" {
		localAddr, err = ipFromInterface(bcastInterface)
		if err != nil {
			return err
		}
	}

	// The address to broadcast to is usually the default `255.255.255.255` but
	// can be overloaded by specifying an override in the CLI arguments.
	bcastAddr := fmt.Sprintf("%s:%s", cliFlags.BroadcastIP, cliFlags.UDPPort)
	udpAddr, err := net.ResolveUDPAddr("udp", bcastAddr)
	if err != nil {
		return err
	}

	// Build the magic packet.
	mp, err := wol.New(macAddr)
	if err != nil {
		return err
	}

	// Grab a stream of bytes to send.
	bs, err := mp.Marshal()
	if err != nil {
		return err
	}

	// Grab a UDP connection to send our packet of bytes.
	conn, err := net.DialUDP("udp", localAddr, udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Printf("Attempting to send a magic packet to MAC %s\n", macAddr)
	fmt.Printf("... Broadcasting to: %s\n", bcastAddr)
	n, err := conn.Write(bs)
	if err == nil && n != 102 {
		err = fmt.Errorf("magic packet sent was %d bytes (expected 102 bytes sent)", n)
	}
	if err != nil {
		return err
	}

	fmt.Printf("Magic packet sent successfully to %s\n", macAddr)
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type cmdFnType func([]string, *Aliases) error

var cmdMap = map[string]cmdFnType{
	"alias":  aliasCmd,
	"list":   listCmd,
	"remove": removeCmd,
	"wake":   wakeCmd,
}

////////////////////////////////////////////////////////////////////////////////

// Helper function to dump the usage and print an error if specified,
// it also returns the exit code requested to the function (saves me a line).
func printUsageGetExitCode(s string, e int) int {
	if len(s) > 0 {
		fmt.Printf(s)
	}
	fmt.Printf(getAppUsageString())
	return e
}

func fatalOnError(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

// Main entry point for binary.
func main() {
	var args []string
	var err error

	// Parse arguments which might get passed to "wol".
	parser := flags.NewParser(&cliFlags, flags.Default & ^flags.HelpFlag)
	args, err = parser.Parse()

	// Disable color if needed.
	if cliFlags.NoColor {
		colorize.DisableColor = true
	}

	ec := 0
	switch {

	// Parse Error, print usage.
	case err != nil:
		fmt.Print(err.Error())
		ec = printUsageGetExitCode("", 1)

	// No arguments, or help requested, print usage.
	case len(os.Args) == 1 || cliFlags.Help:
		ec = printUsageGetExitCode("", 0)

	// "--version" requested.
	case cliFlags.Version:
		fmt.Printf("%s\n", wol.Version)

	// Make sure we are being asked to run a something.
	case len(args) == 0:
		ec = printUsageGetExitCode("No command specified, see usage:\n", 1)

	// All other cases go here.
	case true:
		// Detect the current user to figure out what their ~ is.
		usr, err := user.Current()
		fatalOnError(err)

		// If the user provided a `--db-dir` we expect an existing bolt db
		// at the appropriate path.
		dbDir := path.Join(usr.HomeDir, defaultDBDir)
		if len(cliFlags.DBDir) != 0 {
			dbDir = cliFlags.DBDir
		}

		// Allow the name for the `db` to also be customized. Default is
		// `bolt.db`.
		dbPath := path.Join(dbDir, cliFlags.DBName)

		// Load the list of aliases from the file at dbPath.
		aliases, err := LoadAliases(dbPath)
		fatalOnError(err)
		defer aliases.Close()

		cmd, cmdArgs := strings.ToLower(args[0]), args[1:]
		if fn, ok := cmdMap[cmd]; ok {
			err = fn(cmdArgs, aliases)
		} else {
			err = wakeCmd(args, aliases)
		}
		fatalOnError(err)
	}
	os.Exit(ec)
}
