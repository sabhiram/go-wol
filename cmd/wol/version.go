/*
wol is a command line Magic Packet generator to enable remote Wake-on-LAN.

// TODO: Better package comment

Here are some examples of common usage:

Wake up a machine with mac address 00:11:22:aa:bb:cc
    wol wake 00:11:22:aa:bb:cc

Store an alias
    wol alias skynet 00:11:22:aa:bb:cc

Wake up a machine using an alias
    wol wake skynet

View all aliases and corresponding MAC addresses
    wol list

Delete an alias
    wol remove skynet

Specify the Broadcast Port and IP
    wol wake 00:11:22:aa:bb:cc -b 255.255.255.255 -p 7
    wol wake skynet --bcast 255.255.255.255 --port 7

*/
package main

// Version represents the current Semantic Version of this application
const Version = "1.0.1"
