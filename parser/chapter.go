package parser

import (
	"bytes"
	"io"
	"regexp"
	"strings"
)

type Chapter struct {
	number  string
	content *bytes.Buffer
}

func GetChapters(buf *bytes.Buffer) []Chapter {
	var chapters []Chapter
	var number string
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
				re := regexp.MustCompile("CAPÍTULO [MDCLXVI]+")
				chapterName := re.FindString(currentLine)
				splittedChapterName := strings.Split(chapterName, " ")
				number = splittedChapterName[1]
				continue
			} else {
				chapter := Chapter{
					number:  number,
					content: contentBuf,
				}
				chapters = append(chapters, chapter)
				re := regexp.MustCompile("CAPÍTULO [MDCLXVI]+")
				chapterName := re.FindString(currentLine)
				splittedChapterName := strings.Split(chapterName, " ")
				number = splittedChapterName[1]

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
				number:  number,
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
