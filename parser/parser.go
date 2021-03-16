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
	chapterDiv    DivisionType = "chapter"
	sectionDiv                 = "section"
	subsectionDiv              = "subssection"
	itemDiv                    = "item"
	articleDiv                 = "article"
	paragraphDiv               = "paragraph"
	pointDiv                   = "point"
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
	//Regexes to find each text division
	chapterRegex := "CAPÍTULO [MDCLXVI]+"
	sectionRegex := "Seção [MDCLXVI]+"
	articleRegex := "(Art.( )*[0-9]+)"
	paragraphRegex := "(Parágrafo( )+único.)|(§( )*[1-9]+)"
	itemRegex := "[MDCLXVI]+( )*-"

	doc, err := html.Parse(r)
	if err != nil {
		panic(err)
	}
	mainBuf := &bytes.Buffer{}
	writeTextToBuffer(doc, mainBuf)
	chapters := GetChapters(mainBuf)

	for chapIdx, chapter := range chapters {
		hasSection := findSection(chapter.content)

		//Chapter's children can be sections or articles
		chapterChildren := make([]TextDivision, 0)
		if hasSection {
			//set text(before the first section) and order of chapter
			chapters[chapIdx].Text = selectTextBetweenTwoRegex(chapterRegex, sectionRegex, chapter.content.String())
			chapters[chapIdx].Order = chapIdx + 1

			sections := GetSection(chapter.content, chapter)
			for sectIdx, section := range sections {
				//set text(before the first section) and order of section
				sections[sectIdx].Text = selectTextBetweenTwoRegex(sectionRegex, articleRegex, section.content.String())
				sections[sectIdx].Order = sectIdx + 1

				articles := GetArticle(section.content, section)
				//section's children are always articles, that is why it's directly assigned
				sections[sectIdx].Children = articles

				//Loop through articles inside sessions
				for artIdx, article := range articles {
					//Article's children can be paragraphs or items
					articleChildren := make([]TextDivision, 0)

					hasParagraphs := findParagraph(article.content)
					if hasParagraphs {
						//set text(before the first paragraph or item) and order of article
						articles[artIdx].Text = selectTextBetweenTwoRegex(articleRegex, paragraphRegex, article.content.String())
						articles[artIdx].Order = artIdx + 1

						paragraphs := GetParagraph(article.content, article)
						for paragIdx, paragraph := range paragraphs {
							paragraphs[paragIdx].Text = selectTextBetweenTwoRegex(paragraphRegex, itemRegex, paragraph.content.String())
							paragraphs[paragIdx].Order = paragIdx + 1
							hasItems := findItems(paragraph.content)
							if hasItems {
								items := GetItems(paragraph.content)
								//Paragraph's children are always items, that is why it's directly assigned
								paragraphs[paragIdx].Children = items
							}
							articleChildren = append(articleChildren, paragraphs[paragIdx])
						}
						articles[artIdx].Children = articleChildren
					} else {
						//set text(before the first paragraph or item) and order of article
						articles[artIdx].Text = selectTextBetweenTwoRegex(articleRegex, itemRegex, article.content.String())
						articles[artIdx].Order = artIdx + 1
						//If articles does not have paragraphs
						hasItems := findItems(article.content)
						if hasItems {
							items := GetItems(article.content)
							for itemIdx, _ := range items {
								items[itemIdx].Order = itemIdx + 1
								articleChildren = append(articleChildren, items[itemIdx])
							}
							articles[artIdx].Children = articleChildren
						}
					}
				}
				chapterChildren = append(chapterChildren, sections[sectIdx])
			}
			chapters[chapIdx].Children = chapterChildren

		} else {
			//If chapter does not have section
			//set text(before the first article) and order of chapter
			chapters[chapIdx].Text = selectTextBetweenTwoRegex(chapterRegex, articleRegex, chapter.content.String())
			chapters[chapIdx].Order = chapIdx + 1

			articles := GetArticle(chapter.content, chapter)
			for artIdx, article := range articles {
				//set text(before the first paragraph or item) and order of article
				articles[artIdx].Text = selectTextBetweenTwoRegex(articleRegex, paragraphRegex, article.content.String())
				articles[artIdx].Order = artIdx + 1

				//Article's children can be paragraphs or items
				articleChildren := make([]TextDivision, 0)
				hasParagraphs := findParagraph(article.content)
				if hasParagraphs {
					paragraphs := GetParagraph(article.content, article)
					for paragIdx, paragraph := range paragraphs {
						//set text(before the first item) and order of paragraph
						paragraphs[paragIdx].Text = selectTextBetweenTwoRegex(paragraphRegex, itemRegex, paragraph.content.String())
						paragraphs[paragIdx].Order = paragIdx + 1

						hasItems := findItems(paragraph.content)
						if hasItems {
							items := GetItems(paragraph.content)
							//Paragraph's children are always items, that is why it's directly assigned
							paragraphs[paragIdx].Children = items
						}
						articleChildren = append(articleChildren, paragraphs[paragIdx])
					}
					articles[artIdx].Children = articleChildren
				} else {
					hasItems := findItems(article.content)
					if hasItems {
						items := GetItems(article.content)
						for itemIdx, _ := range items {
							items[itemIdx].Order = itemIdx + 1
							articleChildren = append(articleChildren, items[itemIdx])
						}
						articles[artIdx].Children = articleChildren
					}
				}
				chapterChildren = append(chapterChildren, articles[artIdx])
			}
			chapters[chapIdx].Children = chapterChildren
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

func selectTextBetweenTwoRegex(startReg string, finishRegex string, text string) string {
	regStart := regexp.MustCompile(startReg)
	regFinish := regexp.MustCompile(finishRegex)
	text1 := regFinish.Split(text, 2)[0]
	text2 := regStart.Split(text1, 2)[1]
	result := strings.ReplaceAll(text2, "\n", " ")
	trimmedResult := strings.TrimSpace(result)
	return trimmedResult
}
