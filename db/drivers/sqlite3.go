package drivers

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3" // SQL Driver
)

var (
	dbFilename = "dfind.sq3"
	dbTable    = "files"
)

// SQ3Driver is our DBDriver for sqlite3
type SQ3Driver struct {
	Location string
	DB       *sql.DB
}

// NewSqlite allocates a new sqlite3 driver connection.
func NewSqlite(dataDir string) (*SQ3Driver, error) {
	path := filepath.Join(dataDir, dbFilename)

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Create our structure
	sq := &SQ3Driver{
		Location: path,
		DB:       db,
	}

	// Check for provisioning
	err = sq.checkProvision()
	if err != nil {
		sq.DB.Close()
		sq = nil
	}

	return sq, err
}

// Search is our search function for sqlite3
func (sq *SQ3Driver) Search(pattern string, ignoreCase, seaarchContent bool) ([]*DBObj, error) {
	return nil, nil
}

// Insert adds a new record to the DB, it will return an error if the object already
// exists, unless force is specified.
func (sq *SQ3Driver) Insert(key string, obj *DBObj, force bool) error {
	// For now, just insert the key (and make sure the key matches the obj.Key
	if key != obj.Key {
		return fmt.Errorf("Error: key does not match obj.key:\n    key: %v\nobj.Key:%v", key, obj.Key)
	}

	// If force is NOT set, check for a duplicate
	found, err := sq.queryExists(fmt.Sprintf("SELECT key FROM %v WHERE key=$1", dbTable), key)
	if err != nil {
		return err
	}

	action := "INSERT"

	// If we have one, either error out, or delete
	if found {
		if force {
			// We could also update, but, I'm lazy
			err = sq.Delete(key)
			if err != nil {
				return err
			}
			action = "UPDATE"
		} else {
			return fmt.Errorf("Error: can not insert: item already exists")
		}
	}

	// And, do our insert
	_, err = sq.DB.Exec(fmt.Sprintf("INSERT INTO %v(key, lc_key) VALUES (?, ?)", dbTable), key, strings.ToLower(key))
	if err != nil {
		return err
	}

	fmt.Printf("%v: %v\n", action, key)
	return nil
}

// Delete also returns the object it deleted. The key must be an exact match
func (sq *SQ3Driver) Delete(key string) error {
	_, err := sq.DB.Exec(fmt.Sprintf("DELETE FROM %v WHERE key=$1", dbTable), key)
	return err
}

// queryExists will run the query, and return whether data exists or not
func (sq *SQ3Driver) queryExists(sql string, args ...interface{}) (bool, error) {

	// Run the query
	rows, err := sq.DB.Query(sql, args...)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

// checkProvision will check if the database has been provisioned, and, if not,
// adds our tables and indeces
// TODO: Eventually we should probably use github.com/golang-migrate/migrate
func (sq *SQ3Driver) checkProvision() error {

	// First, let's see if we are a new database by checking the existence of our table(s)
	found, err := sq.queryExists("SELECT name FROM sqlite_master WHERE type='table' AND name=$1", dbTable)
	if err != nil {
		return err
	}
	if !found {
		// We need to create our table
		return sq.provision()
	}

	return nil
}

// provision will create our schema (create tables and indices
func (sq *SQ3Driver) provision() error {

	_, err := sq.DB.Exec(fmt.Sprintf("CREATE TABLE '%v' (key TEXT NOT NULL PRIMARY KEY, lc_key TEXT NOT NULL)",
		dbTable))
	if err != nil {
		return err
	}

	// And, add our case insensitive index
	_, err = sq.DB.Exec(fmt.Sprintf("CREATE INDEX lc_key_idx ON '%v' (lc_key)", dbTable))
	if err != nil {
		return fmt.Errorf("Error Creating Index: %v", err)
	}

	return nil
}
