package parser

import (
	"bytes"
	"io"
	"regexp"
)

type Section struct {
	number       string
	parentNumber string
	content      *bytes.Buffer
}

func GetSection(buf *bytes.Buffer) []Section {
	var sections []Section
	var number string
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
				continue
			} else {
				re := regexp.MustCompile("Seção [MDCLXVI]+")
				number = re.FindString(currentLine)
				section := Section{
					number:  number,
					content: contentBuf,
				}
				sections = append(sections, section)
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
				number:  number,
				content: contentBuf,
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
