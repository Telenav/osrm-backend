// Package nodes2wayblotdb stores `fromNodeID,toNodeID -> wayID` mapping in blotdb.
package nodes2wayblotdb

import (
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

// DB stores `fromNodeID,toNodeID -> wayID` in blotdb.
type DB struct {
	db *bolt.DB
}

const (
	defaultBucket = "bucket"
)

var (
	errEmptyDB = errors.New("empty db")
)

// Open creates/opens a DB structure to store the nodes2way mapping.
func Open(dbFilePath string, readOnly bool) (*DB, error) {
	var db DB

	options := bolt.Options{

		// Default Bolt Options
		Timeout:      0,
		NoGrowSync:   false,
		FreelistType: bolt.FreelistArrayType,

		// set readonly
		ReadOnly: readOnly,
	}

	var err error
	db.db, err = bolt.Open(dbFilePath, 0666, &options)
	if err != nil {
		return nil, err
	}

	if readOnly {
		return &db, nil
	}

	// for write, make sure bucket available
	err = db.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(defaultBucket))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &db, nil
}

// Close releases all database resources.
func (db *DB) Close() error {
	if db.db == nil {
		return errEmptyDB
	}

	return db.db.Close()
}

// Write writes wayID and its nodeIDs into cache or storage.
// wayID: is undirected when input, so will always be positive.
func (db *DB) Write(wayID int64, nodeIDs []int64) error {
	if db.db == nil {
		return errEmptyDB
	}
	if wayID < 0 {
		return fmt.Errorf("expect undirected wayID")
	}
	if len(nodeIDs) < 2 {
		return fmt.Errorf("at least 2 nodes for a way")
	}

	err := db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucket))

		for i := 0; i < len(nodeIDs)-1; i++ {
			if err := b.Put(key(nodeIDs[i], nodeIDs[i+1]), value(wayID)); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
