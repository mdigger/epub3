package main

import (
	"bytes"
)

type Html struct{}

func (self *Html) BlockCode(out *bytes.Buffer, text []byte, lang string)                 {}
func (self *Html) BlockQuote(out *bytes.Buffer, text []byte)                             {}
func (self *Html) BlockHtml(out *bytes.Buffer, text []byte)                              {}
func (self *Html) Header(out *bytes.Buffer, text func() bool, level int, id string)      {}
func (self *Html) HRule(out *bytes.Buffer)                                               {}
func (self *Html) List(out *bytes.Buffer, text func() bool, flags int)                   {}
func (self *Html) ListItem(out *bytes.Buffer, text []byte, flags int)                    {}
func (self *Html) Paragraph(out *bytes.Buffer, text func() bool)                         {}
func (self *Html) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {}
func (self *Html) TableRow(out *bytes.Buffer, text []byte)                               {}
func (self *Html) TableHeaderCell(out *bytes.Buffer, text []byte, align int)             {}
func (self *Html) TableCell(out *bytes.Buffer, text []byte, align int)                   {}
func (self *Html) Footnotes(out *bytes.Buffer, text func() bool)                         {}
func (self *Html) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int)          {}
func (self *Html) TitleBlock(out *bytes.Buffer, text []byte)                             {}

func (self *Html) AutoLink(out *bytes.Buffer, link []byte, kind int)                 {}
func (self *Html) CodeSpan(out *bytes.Buffer, text []byte)                           {}
func (self *Html) DoubleEmphasis(out *bytes.Buffer, text []byte)                     {}
func (self *Html) Emphasis(out *bytes.Buffer, text []byte)                           {}
func (self *Html) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte)    {}
func (self *Html) LineBreak(out *bytes.Buffer)                                       {}
func (self *Html) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {}
func (self *Html) RawHtmlTag(out *bytes.Buffer, text []byte)                         {}
func (self *Html) TripleEmphasis(out *bytes.Buffer, text []byte)                     {}
func (self *Html) StrikeThrough(out *bytes.Buffer, text []byte)                      {}
func (self *Html) FootnoteRef(out *bytes.Buffer, ref []byte, id int)                 {}

func (self *Html) Entity(out *bytes.Buffer, entity []byte)   {}
func (self *Html) NormalText(out *bytes.Buffer, text []byte) {}

func (self *Html) DocumentHeader(out *bytes.Buffer) {}
func (self *Html) DocumentFooter(out *bytes.Buffer) {}

func (self *Html) GetFlags() int { return 0 }

func (self *Html) BlockCodeGithub(out *bytes.Buffer, text []byte, lang string) {}
func (self *Html) BlockCodeNormal(out *bytes.Buffer, text []byte, lang string) {}
func (self *Html) Smartypants(out *bytes.Buffer, text []byte)                  {}
func (self *Html) TocFinalize()                                                {}
func (self *Html) TocHeader(text []byte, level int)                            {}
