package parser

import (
	"bytes"
	"io"
	"regexp"
)

type Chapter struct {
	number   string
	content  *bytes.Buffer
	Label    string
	Text     string
	Order    int
	Children []TextDivision
}

func (c Chapter) getType() DivisionType {
	return chapterDiv
}

func (c Chapter) getNumber() string {
	return c.number
}

func GetChapters(buf *bytes.Buffer) []Chapter {
	var chapters []Chapter
	var chapterNumber string
	var chapterLabel string
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
				chapterNumber = getTitleNumberFromLine("CAPÍTULO [MDCLXVI]+", currentLine)
				chapterLabel = getLabelFromLine("CAPÍTULO [MDCLXVI]+", currentLine)
				continue
			} else {
				chapter := Chapter{
					Label:   chapterLabel,
					number:  chapterNumber,
					content: contentBuf,
				}
				chapters = append(chapters, chapter)
				chapterNumber = getTitleNumberFromLine("CAPÍTULO [MDCLXVI]+", currentLine)
				chapterLabel = getLabelFromLine("CAPÍTULO [MDCLXVI]+", currentLine)
				contentBuf = &bytes.Buffer{}
				contentBuf.WriteString(currentLine)
				continue
			}
		}
		if contentBuf.Len() > 0 {
			contentBuf.WriteString(currentLine)
		}
		if err == io.EOF {
			chapter := Chapter{
				Label:   chapterLabel,
				number:  chapterNumber,
				content: contentBuf,
			}
			chapters = append(chapters, chapter)
			break
		}
		if err != nil {
			break
		}
	}
	return chapters
}
