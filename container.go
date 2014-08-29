package epub

import (
	"encoding/xml"
)

// Predefined names of folder and container file.
const (
	METAINF   = "META-INF"
	CONTAINER = "container.xml"
)

// Container describes the contents of the container
type Container struct {
	XMLName   xml.Name    `xml:"urn:oasis:names:tc:opendocument:xmlns:container container"`
	Version   string      `xml:"version,attr"`
	Rootfiles []*RootFile `xml:"rootfiles>rootfile"`
}

// RootFile describes the path to description of publication.
type RootFile struct {
	FullPath  string `xml:"full-path,attr"`
	MediaType string `xml:"media-type,attr"`
}
