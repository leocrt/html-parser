package parser

import (
	"bytes"
	"io"
	"regexp"
)

type Item struct {
	Label   string
	Order   int
	number  string
	Text    string
	content *bytes.Buffer
}

func (i Item) getType() DivisionType {
	return itemDiv
}

func (i Item) getNumber() string {
	return i.number
}

func findItems(b *bytes.Buffer) bool {
	Matched, err := regexp.MatchString("[MDCLXVI]+( )*-", b.String())
	if err != nil {
		panic(err)
	}
	return Matched
}

func GetItems(buf *bytes.Buffer) []Item {
	var items []Item
	var itemNumber string
	var itemLabel string
	contentBuf := &bytes.Buffer{}
	itemRegex := "[MDCLXVI]+( )*-"
	for {
		currentLine, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}
		// Process the line here.
		Matched, _ := regexp.MatchString(itemRegex, currentLine)
		if Matched {
			if contentBuf.Len() == 0 {
				contentBuf.WriteString(currentLine)
				itemNumber = getTitleNumberFromLine(itemRegex, currentLine)
				itemLabel = getLabelFromLine(itemRegex, currentLine)
				continue
			} else {
				item := Item{
					Label:   itemLabel,
					number:  itemNumber,
					content: contentBuf,
					Text:    contentBuf.String(),
				}
				items = append(items, item)
				itemNumber = getTitleNumberFromLine(itemRegex, currentLine)
				itemLabel = getLabelFromLine(itemRegex, currentLine)
				contentBuf = &bytes.Buffer{}
				contentBuf.WriteString(currentLine)
				continue
			}
		}
		if contentBuf.Len() > 0 {
			contentBuf.WriteString(currentLine)
		}
		if err == io.EOF {
			item := Item{
				Label:   itemLabel,
				number:  itemNumber,
				content: contentBuf,
				Text:    contentBuf.String(),
			}
			items = append(items, item)
			break
		}
		if err != nil {
			break
		}
	}
	return items
}
