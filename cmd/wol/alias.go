package main

import (
	"encoding/gob"
	"os"
	"os/user"
)

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

	err = os.MkdirAll(usr.HomeDir+"/.config/go-wol", 0777)

	if err != nil {
		return ret, err
	}

	file, err = os.Open(usr.HomeDir + "/.config/go-wol/aliases")
	if err != nil {
		return ret, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&ret)
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
