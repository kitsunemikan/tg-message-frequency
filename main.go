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

    report, err := loadReport(*pathToReport)
    if err != nil {
        fmt.Printf("error: load report %q: %v\n", *pathToReport, err)
        os.Exit(1)
    }

    fmt.Printf("Person: %s\n\nChat:\n", report.Name)
    for i := 0; i < 20; i++ {
        if i == len(report.Messages) {
            break
        }

        fmt.Printf("%s: %s\n", report.Messages[i].From, report.Messages[i].Text[0].Text)
    }
}
