// This file contains all of our database search functions

package db

import (
	"fmt"

	"github.com/CodeMonk/dfind/db/drivers"
)

// Insert adds or updates an item in our database.
func (db *DB) Insert(key string, obj *drivers.DBObj, updateIfExists bool) error {
	if db.ReadOnly {
		return fmt.Errorf("Error: Can not insert when ReadOnly")
	}
	return db.Driver.Insert(key, obj, updateIfExists)
}

// Delete removes an item in our database.
func (db *DB) Delete(key string) error {
	if db.ReadOnly {
		return fmt.Errorf("Error: Can not delete when ReadOnly")
	}
	return db.Driver.Delete(key)
}
