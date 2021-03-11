package parser

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

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

	for i, chapter := range chapters {
		if i == 0 || i == 5 {
			fmt.Println(chapter.content)
		}
	}
	// GetLinesByFont(doc, fonts, buf)

	// // items := extractItemsFromBuffer(buf)
	// // fmt.Println(items)
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
