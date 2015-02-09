# go-wol

Simple wake on LAN magic packet generator for golang

## WOL in the world?

[`Wake-on-LAN (WOL)`](http://en.wikipedia.org/wiki/Wake-on-LAN) describes a simple data link layer protocol which tells a listening ethernet interface to power the target machine up.

Each target system typically exposes a setting in it's BIOS which enables or disables the system's WOL capabilities (since this can slowly consume a small amount of standby power).

### Magic Packets (of what?)

The listening interface just looks for a `Magic Packet` with it's MAC address encoded in the WOL scheme. The packet is basically 6 bytes of `0xFF` followed by 16 repetitions of the destination interface's MAC address (102 bytes total). The `Magic Packet` does not have to be exactly 102 bytes, but it's relevant payload is. This payload can be sandwiched anywhere is the broadcast's payload.

It is important to remember that since this is typically sent over the [data link layer](http://en.wikipedia.org/wiki/Data_link_layer), the target machine's IP address is irrelevant.

## Installation

```
$go get github.com/sabhiram/go-wol
$go install github.com/sabhiram/go-wol/...
$wol 08:BA:AD:F0:00:0D
```

## Usage

TODO

## Tests

TODO
