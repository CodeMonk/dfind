// Package db contains our database interface for dfind.
package db

import (
	"github.com/CodeMonk/dfind/db/drivers"
)

// DB holds our database connection. It is used by both the
// searcher and scanner to access the database.
type DB struct {
	ReadOnly bool             // the searcher should only connect read only
	Driver   drivers.DBDriver // our actual driver
}

// New will connect to our database(s) and return a DB instance that
// can be used to access the database
func New(readOnly bool, dataDir string) (*DB, error) {

	sqlite, err := drivers.NewSqlite(dataDir)
	if err != nil {
		return nil, err
	}

	db := &DB{
		ReadOnly: readOnly,
		Driver:   sqlite,
	}

	return db, nil
}
