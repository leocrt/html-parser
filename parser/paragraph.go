package parser

import (
	"bytes"
	"io"
	"regexp"
)

type Paragraph struct {
	Label        string
	number       string
	content      *bytes.Buffer
	ParentType   DivisionType
	ParentNumber string
	Description  string
}

func findParagraph(b *bytes.Buffer) bool {
	Matched, err := regexp.MatchString("(Parágrafo( )+único.)|(§( )*[1-9]+)", b.String())
	if err != nil {
		panic(err)
	}
	return Matched
}

// func checkSingleParagraph(line string) bool {
// 	Matched, err := regexp.MatchString("(Parágrafo único.))", b.String())
// 	if err != nil {
// 		panic(err)
// 	}
// 	return Matched
// }

func GetParagraph(buf *bytes.Buffer, parent TextDivision) []Paragraph {
	var paragraphs []Paragraph
	var paragraphNumber string
	var paragraphLabel string
	contentBuf := &bytes.Buffer{}
	for {
		currentLine, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}

		// Process the line here.
		Matched, _ := regexp.MatchString("(Parágrafo( )+único.)|(§( )*[1-9]+)", currentLine)
		if Matched {
			if contentBuf.Len() == 0 {
				contentBuf.WriteString(currentLine)
				paragraphNumber = getTitleNumberFromLine("(Parágrafo( )+único.)|(§( )*[1-9]+)", currentLine)
				paragraphLabel = getLabelFromLine("(Parágrafo( )+único.)|(§( )*[1-9]+)", currentLine)
				continue
			} else {
				paragraph := Paragraph{
					Label:        paragraphLabel,
					number:       paragraphNumber,
					content:      contentBuf,
					ParentType:   parent.getType(),
					ParentNumber: parent.getNumber(),
					Description:  contentBuf.String(),
				}
				paragraphs = append(paragraphs, paragraph)
				paragraphNumber = getTitleNumberFromLine("(Parágrafo( )+único.)|(§( )*[1-9]+)", currentLine)
				paragraphLabel = getLabelFromLine("(Parágrafo( )+único.)|(§( )*[1-9]+)", currentLine)
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
				ParentType:   parent.getType(),
				ParentNumber: parent.getNumber(),
				Description:  contentBuf.String(),
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
