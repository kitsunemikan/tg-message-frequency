package main

import (
	"flag"
	"fmt"
	"os"
    "path/filepath"
)

var pathToReport = flag.String("report", "", "Path to the JSON report `file` produced by Telegram export feature")

func main() {
    flag.Usage = func() {
        fmt.Fprintf(flag.CommandLine.Output(), "%s [FLAGS]\n\n", filepath.Base(os.Args[0]))
        fmt.Fprintf(flag.CommandLine.Output(), "Flags:\n")

        flag.PrintDefaults()
    }

    flag.Parse()

    if *pathToReport == "" {
        flag.Usage()
        fmt.Print("\nerror: no Telegram JSON message export report specified\n")
        os.Exit(1)
    }

    fmt.Println("OK")
}
