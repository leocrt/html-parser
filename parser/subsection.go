package parser

import (
	"bytes"
	"io"
	"regexp"
)

type Subsection struct {
	number       string
	content      *bytes.Buffer
	Label        string
	Order        int
	Text         string
	parentType   DivisionType
	parentNumber string
	Children     []Article
}

func (s Subsection) getType() DivisionType {
	return sectionDiv
}

func (s Subsection) getNumber() string {
	return s.number
}

func findSubsection(b *bytes.Buffer) bool {
	Matched, err := regexp.MatchString("Subseção [MDCLXVI]+", b.String())
	if err != nil {
		panic(err)
	}
	return Matched
}

func GetSubsection(buf *bytes.Buffer, parent TextDivision) []Subsection {
	var subsections []Subsection
	var subsectionNumber string
	var subsectionLabel string
	subsectionRegex := "Subseção [MDCLXVI]+"
	contentBuf := &bytes.Buffer{}
	for {
		currentLine, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}
		// Process the line here.
		Matched, _ := regexp.MatchString(subsectionRegex, currentLine)
		if Matched {
			if contentBuf.Len() == 0 {
				contentBuf.WriteString(currentLine)
				subsectionNumber = getTitleNumberFromLine(subsectionRegex, currentLine)
				subsectionLabel = getLabelFromLine(subsectionRegex, currentLine)
				continue
			} else {
				subsection := Subsection{
					Label:        subsectionLabel,
					number:       subsectionNumber,
					content:      contentBuf,
					parentType:   parent.getType(),
					parentNumber: parent.getNumber(),
				}
				subsections = append(subsections, subsection)
				subsectionNumber = getTitleNumberFromLine(subsectionRegex, currentLine)
				subsectionLabel = getLabelFromLine(subsectionRegex, currentLine)
				contentBuf = &bytes.Buffer{}
				contentBuf.WriteString(currentLine)
				continue
			}
		}
		if contentBuf.Len() > 0 {
			contentBuf.WriteString(currentLine)
		}
		if err == io.EOF {
			subsection := Subsection{
				Label:        subsectionLabel,
				number:       subsectionNumber,
				content:      contentBuf,
				parentType:   parent.getType(),
				parentNumber: parent.getNumber(),
			}
			subsections = append(subsections, subsection)
			break
		}
		if err != nil {
			break
		}
	}
	return subsections
}
