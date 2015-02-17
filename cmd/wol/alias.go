package main

import (
	"encoding/gob"
	"os"
	"os/user"
)

// Loads the user aliases from the aliases gob stored in
// ~/.config/go-wol/aliases
func loadUserAliases() (map[string]string, error) {
	var file *os.File
	ret := make(map[string]string)

	usr, err := user.Current()
	if err != nil {
		return ret, err
	}

	err = os.MkdirAll(usr.HomeDir+"/.config/go-wol", 0777)

	if err != nil {
		return ret, err
	}

	file, err = os.Open(usr.HomeDir + "/.config/go-wol/aliases")
	defer file.Close()

	if err != nil {
		return ret, err
	}

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&ret)
	return ret, err
}

// Flushes the user aliases to the aliases gob stored in
// ~/.config/go-wol/aliases
func flushUserAliases(m map[string]string) error {
	var file *os.File

	usr, err := user.Current()
	if err != nil {
		return err
	}

	file, err = os.Create(usr.HomeDir + "/.config/go-wol/aliases")
	defer file.Close()

	if err != nil {
		return err
	}

	encoder := gob.NewEncoder(file)
	encoder.Encode(m)
	return nil
}
