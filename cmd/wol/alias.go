package main

////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"sync"

	bolt "go.etcd.io/bbolt"
)

////////////////////////////////////////////////////////////////////////////////

const (
	bucketName = "Aliases"
)

////////////////////////////////////////////////////////////////////////////////

// MacIface holds a MAC Address to wake up, along with an optionally specified
// default interface to use when typically waking up said interface.
type MacIface struct {
	Mac   string
	Iface string
}

// DecodeToMacIface takes a byte buffer and converts decodes it using the gob
// package to a MacIface entry.
func DecodeToMacIface(buf *bytes.Buffer) (MacIface, error) {
	var entry MacIface
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&entry)
	return entry, err
}

// EncodeFromMacIface takes a MAC and an Iface and encodes a gob with a MacIface
// entry.
func EncodeFromMacIface(mac, iface string) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	entry := MacIface{mac, iface}
	err := gob.NewEncoder(buf).Encode(entry)
	return buf, err
}

////////////////////////////////////////////////////////////////////////////////

// Aliases holds a pointer to a mutex which will be acquired and released as
// transactions are carried out on the `db`.
type Aliases struct {
	mtx *sync.Mutex
	db  *bolt.DB
}

// LoadAliases fetches a boltDb entity at a given `dbpath`. The db just contains
// a default bucket called `Aliases` which is where the alias entries are
// stored.
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
		if _, lerr := tx.CreateBucketIfNotExists([]byte(bucketName)); lerr != nil {
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
// exists it is just overwritten.
func (a *Aliases) Add(alias, mac, iface string) error {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	// Create a buffer to store the encoded MAC, interface pair.
	buf, err := EncodeFromMacIface(mac, iface)
	if err != nil {
		return err
	}

	// We don't have to worry about the key existing, as we will update it
	// provided it exists.
	return a.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return bucket.Put([]byte(alias), buf.Bytes())
	})
}

// Del removes an alias from the store based on the alias string.
func (a *Aliases) Del(alias string) error {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	return a.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return bucket.Delete([]byte(alias))
	})
}

// Get retrieves a MacIface from the store based on an alias string.
func (a *Aliases) Get(alias string) (MacIface, error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	var entry MacIface
	err := a.db.View(func(tx *bolt.Tx) error {
		var err error

		bucket := tx.Bucket([]byte(bucketName))
		value := bucket.Get([]byte(alias))
		if value == nil {
			return fmt.Errorf("alias (%s) not found in db", alias)
		}

		entry, err = DecodeToMacIface(bytes.NewBuffer(value))
		return err
	})
	return entry, err
}

// List returns a map containing all alias MacIface pairs.
func (a *Aliases) List() (map[string]MacIface, error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	aliasMap := make(map[string]MacIface, 1)
	err := a.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
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

// Close closes the alias store.
func (a *Aliases) Close() error {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	return a.db.Close()
}
