package parser

import (
	"bytes"
	"io"
	"regexp"
)

type Paragraph struct {
	Label        string
	Order        int
	Text         string
	Children     []Item
	number       string
	content      *bytes.Buffer
	parentType   DivisionType
	parentNumber string
}

func findParagraph(b *bytes.Buffer) bool {
	Matched, err := regexp.MatchString("(Parágrafo( )+único.)|(§( )*[1-9]+)", b.String())
	if err != nil {
		panic(err)
	}
	return Matched
}

func (p Paragraph) getType() DivisionType {
	return paragraphDiv
}

func (p Paragraph) getNumber() string {
	return p.number
}

func GetParagraph(buf *bytes.Buffer, parent TextDivision) []Paragraph {
	var paragraphs []Paragraph
	var paragraphNumber string
	var paragraphLabel string
	paragraphRegex := "(Parágrafo( )+único.)|(§( )*[1-9]+)"
	contentBuf := &bytes.Buffer{}
	for {
		currentLine, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}

		// Process the line here.
		Matched, _ := regexp.MatchString(paragraphRegex, currentLine)
		if Matched {
			if contentBuf.Len() == 0 {
				contentBuf.WriteString(currentLine)
				paragraphNumber = getTitleNumberFromLine(paragraphRegex, currentLine)
				paragraphLabel = getLabelFromLine(paragraphRegex, currentLine)
				continue
			} else {
				paragraph := Paragraph{
					Label:        paragraphLabel,
					number:       paragraphNumber,
					content:      contentBuf,
					parentType:   parent.getType(),
					parentNumber: parent.getNumber(),
					//Description:  contentBuf.String(),
				}
				paragraphs = append(paragraphs, paragraph)
				paragraphNumber = getTitleNumberFromLine(paragraphRegex, currentLine)
				paragraphLabel = getLabelFromLine(paragraphRegex, currentLine)
				contentBuf = &bytes.Buffer{}
				contentBuf.WriteString(currentLine)
				continue
			}
		}
		if contentBuf.Len() > 0 {
			contentBuf.WriteString(currentLine)
		}
		if err == io.EOF {
			paragraph := Paragraph{
				Label:        paragraphLabel,
				number:       paragraphNumber,
				content:      contentBuf,
				parentType:   parent.getType(),
				parentNumber: parent.getNumber(),
				//Description:  contentBuf.String(),
			}
			paragraphs = append(paragraphs, paragraph)
			break
		}
		if err != nil {
			break
		}
	}
	return paragraphs
}
