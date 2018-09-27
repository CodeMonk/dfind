package main

import (
	"flag"
	"fmt"

	"github.com/CodeMonk/dfind/cmd/dfinder/scanner"
	"github.com/CodeMonk/dfind/db"
)

var (
	// Verbose holds our verbose output file
	Verbose = false
	// DataDir is the location of our database files
	DataDir = "/tmp"
	// Root to search from
	Root = "/"
)

func init() {
	flag.BoolVar(&Verbose, "verbose", Verbose, "Turns on verbose output")
	flag.StringVar(&Root, "root", Root, "Root directory to scan from ")
	flag.StringVar(&DataDir, "data_dir", DataDir, "Where to store databases")
}

func main() {

	// Get our root directory (this needs a lot of work)
	flag.Parse()

	if Verbose {
		fmt.Println("Arguments:")
		fmt.Printf("    Verbose: %v", Verbose)
		fmt.Printf("       Root: %v", Root)
	}

	db, err := db.New(false, DataDir)
	if err != nil {
		panic(err)
	}

	s, err := scanner.New(Root, db)
	if err != nil {
		panic(err)
	}

	err = s.ScanInsert()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
