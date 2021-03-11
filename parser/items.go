package parser

import (
	"bytes"
	"io"
	"regexp"
)

type Item struct {
	number string
	text   string
}

func extractItemsFromBuffer(buf *bytes.Buffer) []Item {
	var items []Item
	var item Item
	for {
		currentLine, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}
		// Process the line here.
		Matched, _ := regexp.MatchString("^[MDCLXVI]+", currentLine)
		if Matched {
			items = append(items, item)
			re := regexp.MustCompile("^[MDCLXVI]+")
			number := re.FindString(currentLine)
			item = Item{
				number: number,
				text:   currentLine,
			}
		} else {
			item.text = item.text + " " + currentLine
		}
		if err != nil {
			break
		}
	}
	return items
}
