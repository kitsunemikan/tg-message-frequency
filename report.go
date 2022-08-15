package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"
)

type TextStyle int

const (
    Unknown TextStyle = iota
    Bold
    Italic
    Underline
    Strikethrough
    Monospace
    Spoiler
    Link
    TextLink // In this case Href field is populated
    Hashtag
    Mention
    Phone
)

func (s *TextStyle) UnmarshalJSON(b []byte) error {
    var str string
    if err := json.Unmarshal(b, &str); err != nil {
        return err
    }

    switch strings.ToLower(str) {
    default:
        *s = Unknown
        return errors.New("unknown text style '" + str + "'")
    case "bold":
        *s = Bold
    case "italic":
        *s = Italic
    case "underline":
        *s = Underline
    case "strikethrough":
        *s = Strikethrough
    case "monospace":
        *s = Monospace
    case "spoiler":
        *s = Spoiler
    case "link":
        *s = Link
    case "text_link":
        *s = TextLink
    case "hashtag":
        *s = Hashtag
    case "mention":
        *s = Mention
    case "phone":
        *s = Phone
    }

    return nil
}

type RichPassage struct {
    Styles []TextStyle `json:"type"`
    Text string
    Href string
}

func (p *RichPassage) UnmarshalJSON(b []byte) error {
    var singleText string
    if err := json.Unmarshal(b, &singleText); err == nil {
        p.Text = singleText
        return nil
    } 

    // An unknown formatting option may be in either case with
    // a single string, or an array...
    var possibleErrors error

    var singleStyle struct {
        Style TextStyle `json:"type"`
        Text string
    }
    if err := json.Unmarshal(b, &singleStyle); err == nil {
        p.Styles = []TextStyle{singleStyle.Style}
        p.Text = singleStyle.Text
        return nil
    } else {
        err = fmt.Errorf("rich passage with complex style: %v", err)
        possibleErrors = multierror.Append(possibleErrors, err)
    }

    // Yep, we'll have a copy of RichPassage here, since
    // calling json.Unmarshall on RichPassage itself will cause
    // an infinite recursion
    var multiStyle struct {
        Styles []TextStyle `json:"type"`
        Text string
    }
    if err := json.Unmarshal(b, &multiStyle); err != nil {
        err = fmt.Errorf("rich passage with complex style: %v", err)
        possibleErrors = multierror.Append(possibleErrors, err)
        return possibleErrors
    }

    p.Styles = multiStyle.Styles
    p.Text = multiStyle.Text
    return nil
}

func (p RichPassage) String() string {
    var decor strings.Builder

    styleToStr := func(style TextStyle) string {
        switch style {
        default:
            panic(fmt.Sprintf("RichPassage.String() unhandled style %d", int(style)))
        case Bold:
            return "**"
        case Italic:
            return "_"
        case Underline:
            return ""   // TODO: maybe...
        case Strikethrough:
            return "~~"
        case Monospace:
            return "`"
        case Spoiler:
            return "[!"
        case Link, TextLink, Hashtag, Mention, Phone:
            return ""
        }
    }

    for _, style := range p.Styles {
        decor.WriteString(styleToStr(style))
    }

    content := p.Text
    if len(p.Styles) == 1 && p.Styles[0] == TextLink {
        content = fmt.Sprintf("[%s](%s)", p.Text, p.Href)
    }

    var final strings.Builder
    final.WriteString(decor.String())
    final.WriteString(content)

    for i, decorStr := len(decor.String())-1, decor.String(); i >= 0; i-- {
        final.WriteByte(decorStr[i])
    }

    return final.String()
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

func (t RichText) String() string {
    var str strings.Builder

    for _, passage := range t {
        str.WriteString(passage.String())
    }

    return str.String()
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
