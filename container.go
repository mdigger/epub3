package epub

import (
	"encoding/xml"
)

const (
	METAINF   = "META-INF"
	CONTAINER = "container.xml"
)

type Container struct {
	XMLName   xml.Name    `xml:"urn:oasis:names:tc:opendocument:xmlns:container container"`
	Version   string      `xml:"version,attr"`
	Rootfiles []*RootFile `xml:"rootfiles>rootfile"`
}

type RootFile struct {
	FullPath  string `xml:"full-path,attr"`
	MediaType string `xml:"media-type,attr"`
}
