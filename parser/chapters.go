package parser

import (
	"bytes"
	"io"
	"regexp"
)

type Chapter struct {
	number  string
	content *bytes.Buffer
}

func GetChapters(buf *bytes.Buffer) []Chapter {
	var chapters []Chapter
	contentBuf := &bytes.Buffer{}
	for {
		currentLine, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}
		// Process the line here.
		Matched, _ := regexp.MatchString("CAPÍTULO [MDCLXVI]+", currentLine)
		if Matched {
			if contentBuf.Len() == 0 {
				contentBuf.WriteString(currentLine)
				continue
			} else {
				re := regexp.MustCompile("CAPÍTULO [MDCLXVI]+")
				number := re.FindString(currentLine)
				chapter := Chapter{
					number:  number,
					content: contentBuf,
				}
				chapters = append(chapters, chapter)
				contentBuf = &bytes.Buffer{}
				contentBuf.WriteString(currentLine)
				continue
			}
		}
		if contentBuf.Len() > 0 {
			contentBuf.WriteString(currentLine)
		}
		if err != nil {
			break
		}
	}
	return chapters
}
