package epub

import (
	"encoding/xml"
	"path"
)

// Default file names and mime-type.
var (
	RootPath        = "OEBPS"                // Folder with content of publication
	PackageFilename = "package.opf"          // Package description file name
	EPUBMimeType    = "application/epub+zip" // EPUB mime-type
)

// DefaultContainer is initialized Container with default properties.
var DefaultContainer = &Container{
	Version: "1.0",
	Rootfiles: []*RootFile{
		&RootFile{
			FullPath:  path.Join(RootPath, PackageFilename),
			MediaType: "application/oebps-package+xml",
		},
	},
}

// Predefined names of folder and container file.
const (
	METAINF   = "META-INF"      // Predefined metadata folder
	CONTAINER = "container.xml" // Predefined container file name
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
