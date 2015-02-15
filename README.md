# go-wol

[![Build Status](https://travis-ci.org/sabhiram/go-wol.svg?branch=master)](https://travis-ci.org/sabhiram/go-wol) [![Coverage Status](https://coveralls.io/repos/sabhiram/go-wol/badge.svg)](https://coveralls.io/r/sabhiram/go-wol)

Simple wake on LAN magic packet generator for golang

## WOL in the world?

[`Wake-on-LAN (WOL)`](http://en.wikipedia.org/wiki/Wake-on-LAN) describes a simple data link layer protocol which tells a listening ethernet interface to power the target machine up.

Each target system typically exposes a setting in it's BIOS which enables or disables the system's WOL capabilities (since this can slowly consume a small amount of standby power).

### Magic Packets (of what?)

The listening interface just looks for a `Magic Packet` with it's MAC address encoded in the WOL scheme. The packet is basically 6 bytes of `0xFF` followed by 16 repetitions of the destination interface's MAC address (102 bytes total). The `Magic Packet` does not have to be exactly 102 bytes, but it's relevant payload is. This payload can be sandwiched anywhere in the broadcast's payload.

It is important to remember that since this is typically sent over the [data link layer](http://en.wikipedia.org/wiki/Data_link_layer), the target machine's IP address is irrelevant.

## Installation

```
$go get github.com/sabhiram/go-wol
$go install github.com/sabhiram/go-wol/...
$wol wake 08:BA:AD:F0:00:0D
```

## Usage

Wake up a machine with mac address `00:11:22:aa:bb:cc`
    
    wol wake 00:11:22:aa:bb:cc

Store an alias:
    
    wol alias skynet 00:11:22:aa:bb:cc

Wake up a machine using an alias

    wol wake skynet

View all aliases and corresponding MAC addresses:
    
    wol list

To delete an alias:
    
    wol remove skynet

To specify the Broadcast Port and IP

```
wol wake 00:11:22:aa:bb:cc -b 255.255.255.255 -p 7

# or

wol wake skynet --bcast 255.255.255.255 --port 7
```

#### Defaults

The default Broadcast IP is `255.255.255.255` and the UDP Port is `9`. Typically the UDP port is either `7` or `9`.

#### Alias file

The alias file is typically stored in the user's Home directory under the path of `~/.config/go-wol/aliases`. This is a binary `Gob` which is read from disk on each invocation of `wol`, and flushed to disk if any aliases are modified from the master list.

#### This is how `wol` expects MAC addresses to look

The following MAC addresses are valid and will match:
`01-23-45-56-67-89`, `89:0A:CD:EF:00:12`, `89:0a:cd:ef:00:12`

The following MAC addresses are not (yet) valid:
`1-2-3-4-5-6`, `01 23 45 56 67 89`

## Tests

All commits and PRs will get run on TravisCI and have corresponding coverage reports sent to Coveralls.io.

To run the tests:

    go test -v github.com/sabhiram/go-wol/...
