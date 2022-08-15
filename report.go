package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type TextStyle int

const (
    Unknown = iota
    Italic
    Bold
    Strikethrough
    Monospace
)

func (s *TextStyle) UnmarshalJSON(b []byte) error {
    var str string
    if err := json.Unmarshal(b, &str); err != nil {
        return nil
    }

    switch strings.ToLower(str) {
    default:
        *s = Unknown
        return errors.New("unknown text style '" + str + "'")
    case "italic":
        *s = Italic
    case "bold":
        *s = Bold
    case "monospace":
        *s = Monospace
    case "Strikethrough":
        *s = Strikethrough
    }

    return nil
}

type RichPassage struct {
    Styles []TextStyle
    Text string
}

func (p *RichPassage) UnmarshalJSON(b []byte) error {
    var singleText string
    if err := json.Unmarshal(b, &singleText); err == nil {
        p.Text = singleText
        return nil
    } 

    var singleStyle struct {
        Style TextStyle
        Text string
    }
    if err := json.Unmarshal(b, &singleStyle); err == nil {
        p.Styles = []TextStyle{singleStyle.Style}
        p.Text = singleStyle.Text
        return nil
    }

    // Yep, we'll have a copy of RichPassage here, since
    // calling json.Unmarshall on RichPassage itself will cause
    // an infinite recursion
    var multiStyle struct {
        Styles []TextStyle
        Text string
    }
    if err := json.Unmarshal(b, &multiStyle); err != nil {
        return fmt.Errorf("rich passage with complex style: %v", err)
    }

    p.Styles = multiStyle.Styles
    p.Text = multiStyle.Text
    return nil
}

type RichText []RichPassage

func (t *RichText) UnmarshalJSON(b []byte) error {
    var s string
    if err := json.Unmarshal(b, &s); err == nil {
        *t = []RichPassage{
            {
                Styles: nil,
                Text: s,
            },
        }
        return nil
    }

    var passages []RichPassage
    if err := json.Unmarshal(b, &passages); err != nil {
        return fmt.Errorf("formatted text: %v", err)
    }

    *t = passages
    return nil
}

type Report struct {
    Name string
    Type string
    ID int
    Messages []Message
}

type Message struct {
    ID int
    Type string
    DateUnix int64 `json:"date_unixtime,string"`
    From string
    FromID string `json:"from_id"`
    Text RichText
}

func loadReport(pathToReport string) (Report, error) {
    reportFile, err := os.Open(pathToReport)
    if err != nil {
        return Report{}, err
    }
    defer reportFile.Close()

    reportRaw, err := io.ReadAll(reportFile)
    if err != nil {
        return Report{}, fmt.Errorf("read content: %v", err)
    }

    var report Report
    err = json.Unmarshal(reportRaw, &report)
    if err != nil {
        return Report{}, fmt.Errorf("unmarshal: %v", err)
    }

    return report, nil
}
