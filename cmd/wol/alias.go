package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"os/user"
)

// This struct holds a MAC Address to wake up, along with
// an optionally specified default interface to use when
// typically waking up said interface.
type MacIface struct {
	Mac   string
	Iface string
}

// Loads the user aliases from the aliases gob stored in
// ~/.config/go-wol/aliases
func loadUserAliases() (map[string]MacIface, error) {
	var file *os.File
	ret := make(map[string]MacIface)

	usr, err := user.Current()
	if err != nil {
		return ret, err
	}

	// Grab the ~ for the user, and create the go-wol folder
	// if it does not exist. If the target exists, then
	// MkdirAll will not return an err.
	err = os.MkdirAll(usr.HomeDir + "/.config/go-wol", 0777)
	if err != nil {
		return nil, err
	}

	file, err = os.Open(usr.HomeDir + "/.config/go-wol/aliases")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if decoder.Decode(&ret) != nil {
		fmt.Printf("Unable to load aliases. Resetting aliases list...\n")
		err = flushUserAliases(ret)
	}
	return ret, err
}

// Flushes the user aliases to the aliases gob stored in
// ~/.config/go-wol/aliases
func flushUserAliases(m map[string]MacIface) error {
	var file *os.File

	usr, err := user.Current()
	if err != nil {
		return err
	}

	file, err = os.Create(usr.HomeDir + "/.config/go-wol/aliases")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	encoder.Encode(m)
	return nil
}
