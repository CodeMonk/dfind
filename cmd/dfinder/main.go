package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/CodeMonk/dfind/cmd/dfinder/scanner"
)

func dumpFile(path string, f os.FileInfo, verbose bool) {
	if f.IsDir() == true {
		return
	}
	if verbose {
		fmt.Printf("+-----------------------------------------------------------------------+\n")
		fmt.Printf("| %-69s |\n", path)
		fmt.Printf("+-----------------------------------------------------------------------+\n")
		fmt.Printf("| %-15s | %51s |\n", "Name", f.Name())
		fmt.Printf("| %-15s | %51v |\n", "Size", f.Size())
		fmt.Printf("| %-15s | %51v |\n", "Mode", f.Mode())
		fmt.Printf("| %-15s | %51v |\n", "ModTime", f.ModTime())
		fmt.Printf("| %-15s | %51v |\n", "IsDir", f.IsDir())
		//    fmt.Printf("| %-15s | %51v |\n","Sys", f.Sys())
		//    fmt.Printf("%+v\n",f)
		fmt.Printf("+-----------------------------------------------------------------------+\n")
	} else {
		fmt.Println(path)
	}
}

func visit(path string, f os.FileInfo, err error) error {
	if f == nil {
		return nil
	}
	dumpFile(path, f, false)
	//fmt.Printf("Visited: %s\n", path)
	return nil
}

func main() {

	// Get our root directory (this needs a lot of work)
	flag.Parse()
	root := flag.Arg(0)

	fmt.Printf("root: %s\n", root)

	s, err := scanner.New(root)
	if err != nil {
		panic(err)
	}

	s.Dump(false)

}
