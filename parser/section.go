package parser

import (
	"bytes"
	"io"
	"regexp"
)

type Section struct {
	number       string
	content      *bytes.Buffer
	Label        string
	ParentType   DivisionType
	ParentNumber string
	Articles     []Article
}

func (s Section) getType() DivisionType {
	return sectionDiv
}

func (s Section) getNumber() string {
	return s.number
}

func findSection(b *bytes.Buffer) bool {
	Matched, err := regexp.MatchString("Seção [MDCLXVI]+", b.String())
	if err != nil {
		panic(err)
	}
	return Matched
}

func GetSection(buf *bytes.Buffer, parent TextDivision) []Section {
	var sections []Section
	var sectionNumber string
	var sectionLabel string
	contentBuf := &bytes.Buffer{}
	for {
		currentLine, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}
		// Process the line here.
		Matched, _ := regexp.MatchString("Seção [MDCLXVI]+", currentLine)
		if Matched {
			if contentBuf.Len() == 0 {
				contentBuf.WriteString(currentLine)
				sectionNumber = getTitleNumberFromLine("Seção [MDCLXVI]+", currentLine)
				sectionLabel = getLabelFromLine("Seção [MDCLXVI]+", currentLine)
				continue
			} else {
				section := Section{
					Label:        sectionLabel,
					number:       sectionNumber,
					content:      contentBuf,
					ParentType:   parent.getType(),
					ParentNumber: parent.getNumber(),
				}
				sections = append(sections, section)
				sectionNumber = getTitleNumberFromLine("Seção [MDCLXVI]+", currentLine)
				sectionLabel = getLabelFromLine("Seção [MDCLXVI]+", currentLine)
				contentBuf = &bytes.Buffer{}
				contentBuf.WriteString(currentLine)
				continue
			}
		}
		if contentBuf.Len() > 0 {
			contentBuf.WriteString(currentLine)
		}
		if err == io.EOF {
			section := Section{
				Label:        sectionLabel,
				number:       sectionNumber,
				content:      contentBuf,
				ParentType:   parent.getType(),
				ParentNumber: parent.getNumber(),
			}
			sections = append(sections, section)
			break
		}
		if err != nil {
			break
		}
	}
	return sections
}
