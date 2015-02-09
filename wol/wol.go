package main

import (
    "fmt"
    "encoding/binary"
    "bytes"
    "net"

    wol "github.com/sabhiram/go-wol/lib"
)

const macAddr = "08:BA:DF:00:00:0D"
const udpPort = "9"
const bcastAddr = "255.255.255.255"

func main() {

    x, err := wol.NewMagicPacket(macAddr)
    if err != nil {
        fmt.Printf("Error: %s\n", err.Error())
    } else {
        // Temp code to send magic packet!
        var buf bytes.Buffer
        binary.Write(&buf, binary.BigEndian, x)

        udpAddr, err := net.ResolveUDPAddr("udp", bcastAddr + ":" + udpPort)
        if err != nil {
            fmt.Printf("Unable to get a UDP address for %s\n", bcastAddr)
            return
        }

        connection, err := net.DialUDP("udp", nil, udpAddr)
        if err != nil {
            fmt.Printf("Unable to dial UDP addr for %s\n", bcastAddr)
            return
        }
        defer connection.Close()

        bytesWritten, err := connection.Write(buf.Bytes())
        fmt.Printf("% x\n", buf.Bytes())
        if bytesWritten != 102 {
            fmt.Printf("%d bytes written, %d expected!\n", bytesWritten, 102)
        }
    }
}
