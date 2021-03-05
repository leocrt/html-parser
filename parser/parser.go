package parser

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

type Font struct {
	fontSize   string
	lineHeight string
	fontFamily []string
	color      string
}

//Parse will take an HTML document and return
//a slice of links parsed from it.
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

func GetLinesByFont(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "p" {
		text := &bytes.Buffer{}
		collectText(n, text)
		fmt.Println(text)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		GetLinesByFont(c)
	}
}

func ParseToFont(r io.Reader) {
	doc, err := html.Parse(r)
	if err != nil {
		panic(err)
	}
	GetLinesByFont(doc)
}

func collectText(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}
