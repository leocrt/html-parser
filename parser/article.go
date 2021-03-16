package parser

import (
	"bytes"
	"io"
	"regexp"
)

type Article struct {
	Label        string
	Order        int
	Text         string
	number       string
	content      *bytes.Buffer
	parentType   DivisionType
	parentNumber string
	Children     []TextDivision
}

func (a Article) getType() DivisionType {
	return articleDiv
}

func (a Article) getNumber() string {
	return a.number
}

func GetArticle(buf *bytes.Buffer, parent TextDivision) []Article {
	var articles []Article
	var articleNumber string
	var articleLabel string
	contentBuf := &bytes.Buffer{}
	for {
		currentLine, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}
		// Process the line here.
		Matched, _ := regexp.MatchString("(Art.( )*[0-9]+)", currentLine)
		if Matched {
			if contentBuf.Len() == 0 {
				contentBuf.WriteString(currentLine)
				articleNumber = getTitleNumberFromLine("(Art.( )*[0-9]+)", currentLine)
				articleLabel = getLabelFromLine("(Art.( )*[0-9]+)", currentLine)
				continue
			} else {
				article := Article{
					Label:        articleLabel,
					number:       articleNumber,
					content:      contentBuf,
					parentType:   parent.getType(),
					parentNumber: parent.getNumber(),
				}
				articles = append(articles, article)
				articleNumber = getTitleNumberFromLine("(Art.( )*[0-9]+)", currentLine)
				articleLabel = getLabelFromLine("(Art.( )*[0-9]+)", currentLine)
				contentBuf = &bytes.Buffer{}
				contentBuf.WriteString(currentLine)
				continue
			}
		}
		if contentBuf.Len() > 0 {
			contentBuf.WriteString(currentLine)
		}
		if err == io.EOF {
			article := Article{
				Label:        articleLabel,
				number:       articleNumber,
				content:      contentBuf,
				parentType:   parent.getType(),
				parentNumber: parent.getNumber(),
			}
			articles = append(articles, article)
			break
		}
		if err != nil {
			break
		}
	}
	return articles
}
