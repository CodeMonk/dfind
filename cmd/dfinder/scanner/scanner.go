// Package scanner will scan for files, and send the fileinfo onto a channel
// for DB insertion
//
package scanner

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CodeMonk/dfind/db"
	"github.com/CodeMonk/dfind/db/drivers"
)

// FileType holds our file type
type FileType string

var (
	// FileTypeFile is a regular file
	FileTypeFile = FileType("FILE")
	// FileTypeSymLink is a symbolic link
	FileTypeSymLink = FileType("SYMLINK")
	// FileTypeDirectory is a regular directory
	FileTypeDirectory = FileType("DIRECTORY")
	// FileTypeArchive is an archive
	FileTypeArchive = FileType("ARCHIVE")
	// FileTypeDevice is a Device File
	FileTypeDevice = FileType("DEVICE")
	// FileTypeOther are Named Pipes, etc -- other non regular files
	FileTypeOther = FileType("OTHER")
)

// FileFeed is returned from Scan, and can be read to receive scanned objects.
type FileFeed chan *ScannedObject

// Scanner holds a scanner's state, etc
type Scanner struct {
	RootPath string
	db       *db.DB

	FollowSymlinks bool // If true, will follow symlinks
	OneFilesystem  bool // If true, will not jump filesystems
	fileChan       FileFeed
}

// ScannedObject holds information about a file or archive
type ScannedObject struct {
	Error        error // If this is set, then there was an error, and the scanner is likely dead
	Path         string
	FileInfo     os.FileInfo
	FileType     FileType         // DIR, FILE, ARCHIVE
	ArchiveError error            // If set, something is wrong with archive
	Children     []*ScannedObject // Holds children of an archive
}

// New allocates a new Scanner that can be used to populate the filesystem
func New(root string, db *db.DB) (*Scanner, error) {
	s := &Scanner{
		RootPath: root,
		db:       db,
	}

	// TODO: Make sure root path exists, and return error if not

	return s, nil
}

// Scan will return a channel to read from to receive files as they're being scanned.
func (s *Scanner) Scan() (FileFeed, error) {

	// Make our channel
	s.fileChan = make(chan *ScannedObject, 1000) // Buffered channel with 1000 slots, so we can burst

	// Kick off our scanner
	go s.realScan()

	// And, return our chan
	return s.fileChan, nil
}

// ScanInsert will kick off a scan, and insert all records into the database
func (s *Scanner) ScanInsert() error {
	ch, err := s.Scan()
	if err != nil {
		return err
	}

	// process our channel
	for {
		so, ok := <-ch
		if !ok {
			return nil
		}
		// do the insert now
		obj := &drivers.DBObj{
			Key: so.Path,
			Obj: so,
		}
		err = s.db.Insert(so.Path, obj, true) // force insert it
		if err != nil {
			fmt.Printf("Error inserting %v:\n     %v\n", so.Path, err)
		}
	}
}

// Dump will dump out the scanned items, demonstrating usage of the scanner
func (s *Scanner) Dump(verbose bool) {
	ch, err := s.Scan()
	if err != nil {
		panic(err)
	}

	// process our channel
	for {
		so, ok := <-ch
		if !ok {
			fmt.Printf("Channel closed - exiting Dump")
			return
		}
		so.Dump(verbose)
	}
}

// realScan will scan our files, create the ScannedObject, and return it
// our our channel. It will close the channel after some errors, or
// when it runs out of files.
func (s *Scanner) realScan() {
	fmt.Printf("Starting Scan: %s\n", s.RootPath)

	err := filepath.Walk(s.RootPath, s.visit)
	if err != nil {
		fmt.Printf("ERROR: file.path.Walk() returned %v\n", err)
	}

	// Now close our channel to let the reader know it's done.
	close(s.fileChan)
}

// visit will be called for every file found in rootpath
func (s *Scanner) visit(path string, f os.FileInfo, err error) error {

	// Check for . and .. or other excludes
	if path == "." || path == ".." {
		fmt.Printf("Ignoring: %v\n", path)
		return nil
	}

	so := &ScannedObject{
		Error:    err,
		Path:     path,
		FileInfo: f}

	// If we got an error, send it!
	if err != nil {
		s.fileChan <- so
		return err
	}

	// Otherwise, set our FileType, check for archive, etc.
	err = so.SetFileType()
	s.fileChan <- so

	return err
}

// Dump displays the scanned object to stdout
func (so *ScannedObject) Dump(verbose bool) {
	if verbose {
		fmt.Printf("+-----------------------------------------------------------------------+\n")
		fmt.Printf("| %-69s |\n", so.Path)
		fmt.Printf("+-----------------------------------------------------------------------+\n")
		fmt.Printf("| %-15s | %51s |\n", "Name", so.FileInfo.Name())
		fmt.Printf("| %-15s | %51v |\n", "Size", so.FileInfo.Size())
		fmt.Printf("| %-15s | %51v |\n", "Mode", so.FileInfo.Mode())
		fmt.Printf("| %-15s | %51v |\n", "ModTime", so.FileInfo.ModTime())
		fmt.Printf("| %-15s | %51v |\n", "IsDir", so.FileInfo.IsDir())
		fmt.Printf("+-----------------------------------------------------------------------+\n")
	} else {
		fmt.Println(so.Path)
	}
}

// SetFileType will set our type based on the file mode
func (so *ScannedObject) SetFileType() error {
	m := so.FileInfo.Mode()

	so.FileType = FileTypeOther

	if m&os.ModeDir != 0 {
		so.FileType = FileTypeDirectory
	} else if m&os.ModeSymlink != 0 {
		so.FileType = FileTypeSymLink
	} else if m&os.ModeDevice != 0 || m&os.ModeCharDevice != 0 {
		so.FileType = FileTypeDevice
	} else if m&os.ModeType == 0 {
		// Regular File
		so.FileType = FileTypeFile

		// TODO: Check for archive, set type and add children
	}
	return nil
}
