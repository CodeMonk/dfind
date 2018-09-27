// This file contains all of our database search functions

package db

import (
	"fmt"

	"github.com/CodeMonk/dfind/db/drivers"
)

// Search is our normal search mechanism, that can be called from any DB instance.
// There is also a package-level Search function, but, that allocates a DB instance on the
// fly, and should only be used when you don't need to do anything with the DB except for search.
// Returns: a channel that returns matches, and/or error.
func (db *DB) Search(pattern string, ignoreCase, searchContent bool) (drivers.SearchFeed, error) {
	return db.Driver.Search(pattern, ignoreCase, searchContent)
}

// Search is our top level search function. it should only be used when you do not
// need to do anything except a single search. It is also a good example of how to
// implement channel based searches
func Search(pattern string, ignoreCase, searchContent bool, dataDir string) ([]*drivers.DBObj, error) {
	db, err := New(true, dataDir) // Allocate read-only instance
	if err != nil {
		return nil, err
	}

	ch, err := db.Search(pattern, ignoreCase, searchContent)
	if err != nil {
		return nil, err
	}

	// process our channel
	objs := make([]*drivers.DBObj, 10)
	for {
		obj, ok := <-ch
		if !ok {
			fmt.Printf("Done")
			return objs, nil
		}
		objs = append(objs, obj)
	}
}
