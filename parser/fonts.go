package parser

import "strings"

const CHAPTER_FONT = "{font-size:18px;font-family:Times New Roman,Bold;color:#000000;}"
const ARTICLE_FONT = "{font-size:18px;font-family:Times New Roman;color:#000000;}"
const LINE_BREAK_FONT = "{font-size:18px;line-height:20px;font-family:Times New Roman;color:#000000;}"

type Font struct {
	fontSize   string
	lineHeight string
	fontFamily []string
	color      string
}

type DocumentFonts struct {
	ChapterBoldFont string
	ArticleFont     string
	LineBreakFont   string
}

func MapDocumentFonts(fontMap *map[string]string) *DocumentFonts {
	var font DocumentFonts
	for k, v := range *fontMap {
		if strings.TrimSpace(k) == CHAPTER_FONT {
			font.ChapterBoldFont = v
		}
		if strings.TrimSpace(k) == ARTICLE_FONT {
			font.ArticleFont = v
		}
		if strings.TrimSpace(k) == LINE_BREAK_FONT {
			font.LineBreakFont = v
		}
	}
	return &font
}
