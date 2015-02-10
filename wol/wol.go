package main

import (
    "os"
    "fmt"
    "strings"
    "encoding/binary"
    "bytes"
    "net"

    wol "github.com/sabhiram/go-wol/lib"

    "github.com/jessevdk/go-flags"
)

const (
    UDPPort     = "9"
    BcastAddr   = "255.255.255.255"
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
    fmt.Printf("Attempting to send a magic packet to MAC %s\n", macAddr)
    fmt.Printf("... Broadcasting to IP: %s\n", Options.BroadcastIP)
    fmt.Printf("... Using UDP Port:     %s\n", Options.UDPPort)

    magicPacket, err := wol.NewMagicPacket(macAddr)
    if err != nil {
        fmt.Printf("Error: %s\n", err.Error())
    } else {
        // Temp code to send magic packet!
        var buf bytes.Buffer
        binary.Write(&buf, binary.BigEndian, magicPacket)

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

func main() {
    // Parse arguments which might get passed to "wol"
    parser := flags.NewParser(&Options, flags.Default & ^flags.HelpFlag)
    args, error := parser.Parse()
    macAddr := strings.Join(args, " ")

    exitCode := 0
    switch {

    // Parse Error, print usage
    case error != nil:
        fmt.Printf(getAppUsageString())
        exitCode = 1

    // No arguments, or help requested, print usage
    case len(os.Args) == 1 || Options.Help:
        fmt.Printf(getAppUsageString())

    // "--version" requested
    case Options.Version:
        fmt.Printf("%s\n", Version)

    case len(macAddr) == 0:
        fmt.Printf("Error: A valid mac address is not specified!\n")
        fmt.Printf(getAppUsageString())
        exitCode = 1

    // All other cases go here!
    case true:
        err := sendMagicPacket(macAddr)
        if err != nil {
            exitCode = 1
        } else {
            fmt.Printf("Magic packet sent successfully to %s\n", macAddr)
        }
    }
    os.Exit(exitCode)
}
