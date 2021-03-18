package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	parser "html-parse/parser"
)

func main() {
	f, err := os.Open("bigin41.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	fontMap, err := parser.GetFontMap(tee)
	if err != nil {
		panic(err)
	}
	font := parser.MapDocumentFonts(fontMap)
	fmt.Println(font)
	parser.ParseToFont(&buf, *font)
}
