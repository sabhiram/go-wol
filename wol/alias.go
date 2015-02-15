package main

import (
    "os"
    "os/user"
    "encoding/gob"
)

// Loads the user aliases from the aliases gob stored in
// ~/.config/go-wol/aliases
func loadUserAliases() (map[string]string, error) {
    var file *os.File
    ret := make(map[string]string)

    usr, err := user.Current()
    if err == nil {
        err = os.MkdirAll(usr.HomeDir + "/.config/go-wol", 0777)
    }

    if err == nil {
        file, err = os.Open(usr.HomeDir + "/.config/go-wol/aliases")
        defer file.Close()
    }

    if err == nil {
        decoder := gob.NewDecoder(file)
        err = decoder.Decode(&ret)
    }

    return ret, err
}

// Flushes the user aliases to the aliases gob stored in
// ~/.config/go-wol/aliases
func flushUserAliases(m map[string]string) error {
    var file *os.File

    usr, err := user.Current()
    if err == nil {
        file, err = os.Create(usr.HomeDir + "/.config/go-wol/aliases")
        defer file.Close()
    }

    if err == nil {
        encoder := gob.NewEncoder(file)
        encoder.Encode(m)
    }

    return err
}
