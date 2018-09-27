// Package drivers contains our database driver connections, that can be used interchangeably
// by using interfaces
package drivers

// DBObj holds our database objects, along with the key, for convenience
type DBObj struct {
	Key string
	Obj interface{}
}

// SearchFeed is our channel to receive search results on
type SearchFeed chan *DBObj

// DBDriver defines our database interface, so that different drivers can be swapped out, as needed.
type DBDriver interface {
	Search(string, bool, bool) (SearchFeed, error)
	Insert(string, *DBObj, bool) error
	Delete(string) error
}
