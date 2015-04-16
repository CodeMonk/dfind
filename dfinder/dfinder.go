package main

import (
    "path/filepath"
    "os"
    "flag"
    "fmt"
)

func dumpFile(path string, f os.FileInfo) {
    if f.IsDir() == true {
        return
    }
    fmt.Printf("+-----------------------------------------------------------------------+\n")
    fmt.Printf("| %-69s |\n", path)
    fmt.Printf("+-----------------------------------------------------------------------+\n")
    fmt.Printf("| %-15s | %51s |\n","Name", f.Name())
    fmt.Printf("| %-15s | %51v |\n","Size", f.Size())
    fmt.Printf("| %-15s | %51v |\n","Mode", f.Mode())
    fmt.Printf("| %-15s | %51v |\n","ModTime", f.ModTime())
    fmt.Printf("| %-15s | %51v |\n","IsDir", f.IsDir())
//    fmt.Printf("| %-15s | %51v |\n","Sys", f.Sys())
//    fmt.Printf("%+v\n",f)
    fmt.Printf("+-----------------------------------------------------------------------+\n")
}

func visit(path string, f os.FileInfo, err error) error {
    if f == nil {
        return nil
    }
    dumpFile(path, f)
    //fmt.Printf("Visited: %s\n", path)
    return nil
}

func main() {

    flag.Parse()

    root := flag.Arg(0)

    fmt.Printf("root: %s\n", root)

    err := filepath.Walk(root, visit)
    fmt.Printf("file.path.Walk() returned %v\n", err)
}