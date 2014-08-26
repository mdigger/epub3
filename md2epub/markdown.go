package main

import (
	"bytes"
	"fmt"
	"github.com/russross/blackfriday"
	"html"
	"strconv"
	"strings"
)

type Html struct {
	css   string
	title string
	lang  string
}

func (self *Html) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	count := 0
	for _, elt := range strings.Fields(lang) {
		if elt[0] == '.' {
			elt = elt[1:]
		}
		if len(elt) == 0 {
			continue
		}
		out.WriteString("<pre lang=\"")
		out.WriteString(html.EscapeString(elt))
		out.WriteString("\"><code>")
		count++
		break
	}
	if count == 0 {
		out.WriteString("<pre><code>")
	}
	out.WriteString(html.EscapeString(string(text)))
	out.WriteString("</code></pre>\n")
}

func (self *Html) BlockQuote(out *bytes.Buffer, text []byte) {
	out.WriteString("<blockquote>\n")
	out.Write(text)
	out.WriteString("</blockquote>\n")
}

func (self *Html) BlockHtml(out *bytes.Buffer, text []byte) {
	out.Write(text)
	out.WriteByte('\n')
}

func (self *Html) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	marker := out.Len()
	out.WriteString(fmt.Sprintf("<h%d>", level))
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString(fmt.Sprintf("</h%d>\n", level))
}

func (self *Html) HRule(out *bytes.Buffer) {
	out.WriteString("<hr />")
}

func (self *Html) List(out *bytes.Buffer, text func() bool, flags int) {
	marker := out.Len()
	if flags&blackfriday.LIST_TYPE_ORDERED != 0 {
		out.WriteString("<ol>")
	} else {
		out.WriteString("<ul>")
	}
	if !text() {
		out.Truncate(marker)
		return
	}
	if flags&blackfriday.LIST_TYPE_ORDERED != 0 {
		out.WriteString("</ol>\n")
	} else {
		out.WriteString("</ul>\n")
	}
}

func (self *Html) ListItem(out *bytes.Buffer, text []byte, flags int) {
	out.WriteString("<li>")
	out.Write(text)
	out.WriteString("</li>\n")
}

func (self *Html) Paragraph(out *bytes.Buffer, text func() bool) {
	marker := out.Len()
	out.WriteString("<p>")
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("</p>\n")
}

func (self *Html) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
	out.WriteString("<table>\n<thead>\n")
	out.Write(header)
	out.WriteString("</thead>\n\n<tbody>\n")
	out.Write(body)
	out.WriteString("</tbody>\n</table>\n")
}

func (self *Html) TableRow(out *bytes.Buffer, text []byte) {
	out.WriteString("<tr>\n")
	out.Write(text)
	out.WriteString("\n</tr>\n")
}

func (self *Html) TableHeaderCell(out *bytes.Buffer, text []byte, flags int) {

}

func (self *Html) TableCell(out *bytes.Buffer, text []byte, align int) {
	switch align {
	case blackfriday.TABLE_ALIGNMENT_LEFT:
		out.WriteString("<td align=\"left\">")
	case blackfriday.TABLE_ALIGNMENT_RIGHT:
		out.WriteString("<td align=\"right\">")
	case blackfriday.TABLE_ALIGNMENT_CENTER:
		out.WriteString("<td align=\"center\">")
	default:
		out.WriteString("<td>")
	}
	out.Write(text)
	out.WriteString("</td>")
}

func (self *Html) Footnotes(out *bytes.Buffer, text func() bool) {
	out.WriteString("<div class=\"footnotes\">\n")
	self.HRule(out)
	self.List(out, text, blackfriday.LIST_TYPE_ORDERED)
	out.WriteString("</div>\n")
}

func (self *Html) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
	out.WriteString(`<li id="fn:`)
	out.Write(slugify(name))
	out.WriteString(`">`)
	out.Write(text)
	out.WriteString("</li>\n")
}

func (self *Html) TitleBlock(out *bytes.Buffer, text []byte) {

}

func (self *Html) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	out.WriteString("<a href=\"")
	if kind == blackfriday.LINK_TYPE_EMAIL {
		out.WriteString("mailto:")
	}
	out.WriteString(html.EscapeString(string(link)))
	out.WriteString("\">")
	switch {
	case bytes.HasPrefix(link, []byte("mailto://")):
		out.WriteString(html.EscapeString(string(link[len("mailto://"):])))
	case bytes.HasPrefix(link, []byte("mailto:")):
		out.WriteString(html.EscapeString(string(link[len("mailto:"):])))
	default:
		out.WriteString(html.EscapeString(string(link)))
	}
	out.WriteString("</a>")
}

func (self *Html) CodeSpan(out *bytes.Buffer, text []byte) {
	out.WriteString("<code>")
	out.WriteString(html.EscapeString(string(text)))
	out.WriteString("</code>")
}

func (self *Html) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	out.WriteString("<strong>")
	out.Write(text)
	out.WriteString("</strong>")
}

func (self *Html) Emphasis(out *bytes.Buffer, text []byte) {
	if len(text) == 0 {
		return
	}
	out.WriteString("<em>")
	out.Write(text)
	out.WriteString("</em>")
}

func (self *Html) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	out.WriteString("<img src=\"")
	out.WriteString(html.EscapeString(string(link)))
	out.WriteString("\" alt=\"")
	if len(alt) > 0 {
		out.WriteString(html.EscapeString(string(alt)))
	}
	if len(title) > 0 {
		out.WriteString("\" title=\"")
		out.WriteString(html.EscapeString(string(title)))
	}
	out.WriteString("\" />")
}

func (self *Html) LineBreak(out *bytes.Buffer) {
	out.WriteString("<br />")
}

func (self *Html) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	out.WriteString("<a href=\"")
	out.WriteString(html.EscapeString(string(link)))
	if len(title) > 0 {
		out.WriteString("\" title=\"")
		out.WriteString(html.EscapeString(string(title)))
	}
	out.WriteString("\">")
	out.Write(content)
	out.WriteString("</a>")
}

func (self *Html) RawHtmlTag(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (self *Html) TripleEmphasis(out *bytes.Buffer, text []byte) {
	out.WriteString("<strong><em>")
	out.Write(text)
	out.WriteString("</em></strong>")
}

func (self *Html) StrikeThrough(out *bytes.Buffer, text []byte) {
	out.WriteString("<s>")
	out.Write(text)
	out.WriteString("</s>")
}

func (self *Html) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
	slug := slugify(ref)
	out.WriteString(`<sup class="footnote-ref" id="fnref:`)
	out.Write(slug)
	out.WriteString(`"><a rel="footnote" href="#fn:`)
	out.Write(slug)
	out.WriteString(`">`)
	out.WriteString(strconv.Itoa(id))
	out.WriteString(`</a></sup>`)
}

func (self *Html) Entity(out *bytes.Buffer, entity []byte) {
	out.Write(entity)
}

func (self *Html) NormalText(out *bytes.Buffer, text []byte) {
	out.WriteString(html.EscapeString(string(text)))
}

func (self *Html) DocumentHeader(out *bytes.Buffer) {
	out.WriteString("<!DOCTYPE html>\n")
	out.WriteString("<html xmlns=\"http://www.w3.org/1999/xhtml\"")
	if self.lang != "" {
		out.WriteString(" xml:lang=\"")
		out.WriteString(self.lang)
		out.WriteRune('"')
	}
	out.WriteString(">\n")
	out.WriteString("<head>\n")
	out.WriteString("  <title>")
	self.NormalText(out, []byte(self.title))
	out.WriteString("</title>\n")
	out.WriteString("  <meta charset=\"utf-8\" />")
	out.WriteString(">\n")
	if self.css != "" {
		out.WriteString("  <link rel=\"stylesheet\" type=\"text/css\" href=\"")
		out.WriteString(html.EscapeString(self.css))
		out.WriteString("\" />")
	}
	out.WriteString("</head>\n")
	out.WriteString("<body>\n")
}

func (self *Html) DocumentFooter(out *bytes.Buffer) {
	out.WriteString("</body>\n")
	out.WriteString("</html>\n")
}

func (self *Html) GetFlags() int {
	return 0
}
