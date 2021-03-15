package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type DivisionType string

const (
	chapterDiv   DivisionType = "chapter"
	sectionDiv                = "section"
	itemDiv                   = "item"
	articleDiv                = "article"
	paragraphDiv              = "paragraph"
	pointDiv                  = "point"
)

type TextDivision interface {
	getType() DivisionType
	getNumber() string
}

type standardDivision struct{}

func GetFontMap(r io.Reader) (*map[string]string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	fontMap := make(map[string]string)
	fontMapResult := dfs(doc, &fontMap)
	return fontMapResult, nil
}

func dfs(n *html.Node, fontMap *map[string]string) *map[string]string {
	if n.Data == "style" {
		fontMap = getFontMapFromCommentNode(n.FirstChild.Data)
		return fontMap
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		fontMap = dfs(c, fontMap)
	}
	return fontMap
}

func getFontMapFromCommentNode(text string) *map[string]string {
	text = strings.ReplaceAll(text, "<!--", "")
	text = strings.ReplaceAll(text, "-->", "")
	s := strings.Split(text, ".")
	s = s[1:]
	fontMap := make(map[string]string)
	for _, css := range s {
		parseCssToFont(css, &fontMap)
	}
	return &fontMap
}

func parseCssToFont(css string, fontMap *map[string]string) *map[string]string {
	fontTitle := css[:4]
	cssClass := css[4:]
	(*fontMap)[cssClass] = fontTitle
	return fontMap
}

func ParseToFont(r io.Reader, fonts DocumentFonts) {
	doc, err := html.Parse(r)
	if err != nil {
		panic(err)
	}
	mainBuf := &bytes.Buffer{}
	writeTextToBuffer(doc, mainBuf)
	chapters := GetChapters(mainBuf)

	// sections := GetSection(chapters[0].content, chapters[0].number)
	// for _, section := range sections {
	// 	fmt.Println(section.content)
	// }

	for chapIdx, chapter := range chapters {
		//fmt.Println(chapter.content)
		hasSection := findSection(chapter.content)
		children := make([]TextDivision, 0)
		if hasSection {
			str1 := regexp.MustCompile("Seção [MDCLXVI]+")
			chapters[chapIdx].Text = str1.Split(chapter.content.String(), 2)[0]
			chapters[chapIdx].Order = chapIdx + 1
			sections := GetSection(chapter.content, chapter)
			for sectIdx, section := range sections {
				articles := GetArticle(section.content, section)
				sections[sectIdx].Articles = articles
				for artIdx, article := range articles {
					hasParagraphs := findParagraph(article.content)
					if hasParagraphs {
						paragraphs := GetParagraph(article.content, article)
						articles[artIdx].Paragraphs = paragraphs
					}
				}
				children = append(children, sections[sectIdx])
			}
			chapters[chapIdx].Children = children
		} else {
			articles := GetArticle(chapter.content, chapter)
			for artIdx, article := range articles {
				hasParagraphs := findParagraph(article.content)
				if hasParagraphs {
					paragraphs := GetParagraph(article.content, article)
					articles[artIdx].Paragraphs = paragraphs
				}
				children = append(children, articles[artIdx])
			}
			chapters[chapIdx].Children = children
		}
	}
	b, err := json.Marshal(chapters)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func GetLinesByFont(n *html.Node, font DocumentFonts, buf *bytes.Buffer) {
	if n.Type == html.ElementNode &&
		n.Data == "p" &&
		(strings.Contains(n.Attr[1].Val, font.ArticleFont) ||
			strings.Contains(n.Attr[1].Val, font.LineBreakFont) ||
			strings.Contains(n.Attr[1].Val, font.ChapterBoldFont)) {

		collectText(n, buf)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		GetLinesByFont(c, font, buf)
	}
}

func collectText(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(strings.ReplaceAll(n.Data, ";", " "))
		buf.WriteString("\n")
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}

func writeTextToBuffer(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.ElementNode &&
		n.Data == "p" {
		collectText(n, buf)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		writeTextToBuffer(c, buf)
	}
}

func getTitleNumberFromLine(regex string, currentLine string) string {
	if strings.Contains(currentLine, "CAPÍTULO") {
		re := regexp.MustCompile(regex)
		chapterName := re.FindString(currentLine)
		splittedChapterName := strings.Split(chapterName, " ")
		reNumber := regexp.MustCompile("([MDCLXVI]+)")
		number := reNumber.FindString(splittedChapterName[1])
		return number
	}
	re := regexp.MustCompile(regex)
	chapterName := re.FindString(currentLine)
	reNumber := regexp.MustCompile("([0-9]+|[MDCLXVI]+)")
	number := reNumber.FindString(chapterName)
	return number
}

func getLabelFromLine(regex string, currentLine string) string {
	re := regexp.MustCompile(regex)
	label := re.FindString(currentLine)
	return label
}
