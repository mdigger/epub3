package epub

import (
	"encoding/xml"
)

// Default file names and mime-type.
var (
	RootPath        = "OEBPS"       // Folder with content of publication
	PackageFilename = "package.opf" // Package description file name
)

// RootFile describes the path to description of publication.
type RootFile struct {
	FullPath  string `xml:"full-path,attr"`
	MediaType string `xml:"media-type,attr"`
}

// Container describes the contents of the container.
type Container struct {
	XMLName   xml.Name   `xml:"urn:oasis:names:tc:opendocument:xmlns:container container"`
	Version   string     `xml:"version,attr"`
	Rootfiles []RootFile `xml:"rootfiles>rootfile"`
}
