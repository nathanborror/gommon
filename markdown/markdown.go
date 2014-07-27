package markdown

import (
	"fmt"

	"github.com/russross/blackfriday"
)

// Markdown returns HTML given a string of markdown text
func Markdown(args ...interface{}) string {
	s := blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...)))
	return string(s)
}
