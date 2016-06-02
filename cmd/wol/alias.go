package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
	"path"
	"sync"

	"github.com/boltdb/bolt"
)

const (
	MainBucketName string = "Aliases"
)

// This struct holds a pointer to a mutex which will be acquired and released
// as transactions are carried out on the `db`.
type Aliases struct {
	mtx *sync.Mutex
	db  *bolt.DB
}

// This struct holds a MAC Address to wake up, along with an optionally
// specified default interface to use when typically waking up said interface.
type MacIface struct {
	Mac   string
	Iface string
}

// Takes a byte buffer and converts decodes it using the gob package
// to a MacIface entry
func DecodeToMacIface(buf *bytes.Buffer) (MacIface, error) {
	var entry MacIface
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&entry)
	return entry, err
}

// Takes a MAC and an Iface and encodes a gob with a MacIface entry
func EncodeFromMacIface(mac, iface string) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	entry := MacIface{mac, iface}
	err := gob.NewEncoder(buf).Encode(entry)
	return buf, err
}

// This function fetches a boltDb entity at a given `dbpath`. The db just
// contains a default bucket called `Aliases` which is where the alias
// entries are stored.
func LoadAliases(dbpath string) (*Aliases, error) {

	err := os.MkdirAll(path.Dir(dbpath), os.ModePerm)
	if os.IsNotExist(err) {
		return nil, err
	}

	db, err := bolt.Open(dbpath, 0660, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		if _, lerr := tx.CreateBucketIfNotExists([]byte(MainBucketName)); lerr != nil {
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

// Add updates an alias entry or adds a new alias entry. If the alias already
// exists it is just overwritten
func (a *Aliases) Add(alias, mac, iface string) error {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	// Create a buffer to store the encoded MAC, interface pair
	buf, err := EncodeFromMacIface(mac, iface)
	if err != nil {
		return err
	}

	// We don't have to worry about the key existing, as we will just
	// update it if it is already there
	return a.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(MainBucketName))
		return bucket.Put([]byte(alias), buf.Bytes())
	})
}

// Del removes an alias from the store based on the alias string
func (a *Aliases) Del(alias string) error {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	return a.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(MainBucketName))
		return bucket.Delete([]byte(alias))
	})
}

// Get retrieves a MacIface from the store based on an alias string
func (a *Aliases) Get(alias string) (MacIface, error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	var entry MacIface
	err := a.db.View(func(tx *bolt.Tx) error {
		var err error

		bucket := tx.Bucket([]byte(MainBucketName))
		value := bucket.Get([]byte(alias))
		if value == nil {
			return errors.New("Alias not found in db")
		}

		entry, err = DecodeToMacIface(bytes.NewBuffer(value))
		return err
	})
	return entry, err
}

// List returns a map containing all alias MacIface pairs
func (a *Aliases) List() (map[string]MacIface, error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	aliasMap := make(map[string]MacIface, 1)
	err := a.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(MainBucketName))
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			if entry, err := DecodeToMacIface(bytes.NewBuffer(v)); err == nil {
				aliasMap[string(k)] = entry
			} else {
				return err
			}
		}
		return nil
	})
	return aliasMap, err
}

// Close closes the alias store
func (a *Aliases) Close() error {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	return a.db.Close()
}
