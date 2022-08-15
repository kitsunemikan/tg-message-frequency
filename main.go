package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
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

    fmt.Printf("# Chat with %s\n\nChat:\n", report.Name)
    fromMsg := len(report.Messages) - 100
    if fromMsg < 0 {
        fromMsg = 0
    }

    var lastDate time.Time
    for i := 0; i < len(report.Messages); i++ {
        if i == len(report.Messages) {
            break
        }

        date := time.Unix(report.Messages[i].DateUnix, 0)
        if date.Day() != lastDate.Day() {
            fmt.Printf("## %s\n\n", date.Format("January _2, 2006"))
        }

        lastDate = date
        timestamp := date.Format("15:04")
        datestamp := date.Format("01/02")
        fmt.Printf("**%s %s %s:** %s\n\n", timestamp, datestamp, report.Messages[i].From, report.Messages[i].Text.String())
    }
}
