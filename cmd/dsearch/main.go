// dsearch is the search side of dfind

package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/CodeMonk/dfind/db"
)

var (
	// Verbose holds our verbose output file
	Verbose = false
	// DataDir is the location of our database files
	DataDir = "/tmp"
	// Insensitive searches
	Insensitive = false
)

func init() {
	flag.BoolVar(&Verbose, "verbose", Verbose, "Turns on verbose output")
	flag.BoolVar(&Insensitive, "insensitive", Insensitive, "Perform search ignoring case")
	flag.StringVar(&DataDir, "data_dir", DataDir, "Where to store databases")
}

func main() {

	// Get our root directory (this needs a lot of work)
	flag.Parse()

	pattern := strings.Join(flag.Args(), " ")

	if Verbose {
		fmt.Println("Arguments:")
		fmt.Printf("        Verbose: %v\n", Verbose)
		fmt.Printf("    Insensitive: %v\n", Insensitive)
		fmt.Printf("        Pattern: %v\n", pattern)
	}

	// Open database read only
	db, err := db.New(true, DataDir)
	if err != nil {
		panic(err)
	}

	ch, err := db.Search(pattern, Insensitive, false)
	if err != nil {
		fmt.Printf("Error starting search: %v\n", err)
		return
	}

	// process our channel
	for {
		obj, ok := <-ch
		if !ok {
			fmt.Printf("Done")
			return
		}
		fmt.Println(obj.Key)
	}
}
