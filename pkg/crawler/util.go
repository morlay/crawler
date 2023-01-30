package crawler

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func fixURI(path string, origin string) string {
	if strings.HasPrefix(path, "http:") || strings.HasPrefix(path, "https:") {
		return path
	}
	return origin + path
}

func DocumentFromBytes(data []byte) *goquery.Document {
	d, _ := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	return d
}
