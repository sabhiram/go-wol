package main

import (
    "os"
    "fmt"
    "encoding/binary"
    "bytes"
    "net"

    "errors"

    "github.com/sabhiram/go-colorize"
    wol "github.com/sabhiram/go-wol/lib"

    "github.com/jessevdk/go-flags"
)

var (
    // Define holders for the cli arguments we wish to parse
    Options struct {
        Version     bool   `short:"v" long:"version"`
        Help        bool   `short:"h" long:"help"`
        BroadcastIP string `short:"b" long:"bcast" default:"255.255.255.255"`
        UDPPort     string `short:"p" long:"port" default:"9"`
    }
)

// Function to send a magic packet to a given mac address
func sendMagicPacket(macAddr string) error {
    magicPacket, err := wol.NewMagicPacket(macAddr)
    if err == nil {
        var buf bytes.Buffer
        binary.Write(&buf, binary.BigEndian, magicPacket)

        fmt.Printf("Attempting to send a magic packet to MAC %s\n", macAddr)
        fmt.Printf("... Broadcasting to IP: %s\n", Options.BroadcastIP)
        fmt.Printf("... Using UDP Port:     %s\n", Options.UDPPort)


        udpAddr, err := net.ResolveUDPAddr("udp", Options.BroadcastIP + ":" + Options.UDPPort)
        if err != nil {
            fmt.Printf("Unable to get a UDP address for %s\n", Options.BroadcastIP)
            return err
        }

        connection, err := net.DialUDP("udp", nil, udpAddr)
        if err != nil {
            fmt.Printf("Unable to dial UDP addr for %s\n", Options.BroadcastIP)
            return err
        }
        defer connection.Close()

        bytesWritten, err := connection.Write(buf.Bytes())
        if bytesWritten != 102 {
            fmt.Printf("%d bytes written, %d expected!\n", bytesWritten, 102)
        }
    }
    return err
}

// Run one of the supported commands
func runCommand(cmd string, args []string, aliases map[string]string) error {
    var err error

    switch cmd {

    case "alias":
        fmt.Printf("%s\n", cmd)
        fmt.Printf("ALIASES: %v\n", aliases);

    case "list":
        if len(aliases) == 0 {
            fmt.Printf("No aliases found! Add one with \"wol alias <name> <mac>\"\n")
        } else {
            for alias, mac := range aliases {
                fmt.Printf("    %s        %s\n", mac, alias)
            }
        }

    case "remove":
        fmt.Printf("%s\n", cmd)
        fmt.Printf("ALIASES: %v\n", aliases);

    case "wake":
        if len(args) > 0 {
            err = sendMagicPacket(args[0])
            if err == nil {
                fmt.Printf("Magic packet sent successfully to %s\n", args[0])
            }
        } else {
            err = errors.New("No mac address specified to wake command")
        }
    default:
        panic("Invalid command passed to runCommand")

    }
    return err
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

func loadUserAliases() (map[string]string, error) {
    var ret map[string]string

    err := os.MkdirAll("~/.config/go-wol/", 0777)
    if err != nil {
        return ret, err
    }

    fmt.Printf("TODO: ... loading user aliases\n");

    return ret, err
}

func flushUserAliases(m map[string]string) error {
    fmt.Printf("TODO: ... flush user aliases\n");
    return nil
}

func main() {
    var args []string

    // Load the list of aliases from ~/.config/go-wol/aliases
    aliases, err := loadUserAliases()

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
                exitCode = printUsageGetExitCode(
                    fmt.Sprintf("%s\n", err.Error()), 1)
            }
        } else {
            exitCode = printUsageGetExitCode(
                fmt.Sprintf("Unknown command %s, see usage:\n", colorize.ColorString(cmd, "red")), 1)
        }

    }
    os.Exit(exitCode)
}
