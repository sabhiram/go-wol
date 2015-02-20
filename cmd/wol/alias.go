package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/boltdb/bolt"
	"os/user"
	"sync"
)

type Aliases struct {
	mtx *sync.Mutex
	db  *bolt.DB
}

const (
	MainBucketName string = "Aliases"
)

// This struct holds a MAC Address to wake up, along with
// an optionally specified default interface to use when
// typically waking up said interface.
type MacIface struct {
	Mac   string
	Iface string
}

func OpenDB(path string) (*Aliases, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(usr.HomeDir+"/.config/go-wol", 0660, nil)
	if err != nil {
		return nil, err
	}
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, lerr := tx.CreateBucketIfNotExists([]byte(MainBucketName)); err != nil {
			return lerr
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &Aliases{
		mtx: &sync.Mutex{},
		db:  db,
	}, nil
}

//Add updates an alias entry or adds a new alias entry
//If the alias already exists it is just overwritten
func (a *Aliases) Add(alias, mac, iface string) error {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	//we don't have to worry about the key existing, as we will just
	//update it if it is already there
	bb := bytes.NewBuffer(nil)
	genc := gob.NewEncoder(bb)
	if err := genc.Encode(MacIface{Mac: mac, Iface: iface}); err != nil {
		return err
	}
	if err := a.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(MainBucketName))
		if err := b.Put([]byte(alias), bb.Bytes()); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

//Del removes an alias from the store based on the alias string
func (a *Aliases) Del(alias string) error {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	if err := a.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(MainBucketName))
		if err := b.Delete([]byte(alias)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

//Get retrieves a MacIface from the store based on an alias string
func (a *Aliases) Get(alias string) (MacIface, error) {
	var mi MacIface
	a.mtx.Lock()
	defer a.mtx.Unlock()
	if err := a.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(MainBucketName))
		v := b.Get([]byte(alias))
		if v == nil {
			return errors.New("Alias not found")
		}
		bb := bytes.NewBuffer(v)
		gdec := gob.NewDecoder(bb)
		if err := gdec.Decode(&mi); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return mi, err
	}
	return mi, nil
}

//Close closes the alias store
func (a *Aliases) Close() error {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	return a.db.Close()
}

//List returns a map containing all alias MacIface pairs
func (a *Aliases) List() (map[string]MacIface, error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	mp := make(map[string]MacIface, 1)
	if err := a.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(MainBucketName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var mi MacIface
			bb := bytes.NewBuffer(v)
			gdec := gob.NewDecoder(bb)
			if err := gdec.Decode(&mi); err != nil {
				return err
			}
			mp[string(k)] = mi
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return mp, nil
}
