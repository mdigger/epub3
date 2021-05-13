// Package epub describes the format of the data used in the epub, and offers easy writer to create
// publications in this format.
package epub

import (
	"encoding/xml"
)

// Package element is the root container of the Package Document and encapsulates Publication
// metadata and resource information.
type Package struct {
	XMLName          xml.Name    `xml:"http://www.idpf.org/2007/opf package"`
	Version          string      `xml:"version,attr"`                 // Specifies the EPUB specification version to which the Publication conforms
	UniqueIdentifier string      `xml:"unique-identifier,attr"`       // An IDREF that identifies the dc:identifier element that provides the package's preferred, or primary, identifier
	Prefix           string      `xml:"prefix,attr,omitempty"`        // Declaration mechanism for prefixes not reserved by this specification.
	Lang             string      `xml:"xml:lang,attr,omitempty"`      // Specifies the language used in the contents and attribute values of the carrying element and its descendants
	Dir              string      `xml:"dir,attr,omitempty"`           // Specifies the base text direction of the content and attribute values of the carrying element and its descendants.
	ID               string      `xml:"id,attr,omitempty"`            // The ID of this element, which must be unique within the document scope
	Metadata         Metadata    `xml:"metadata"`                     // The metadata element encapsulates Publication meta information
	Manifest         Manifest    `xml:"manifest"`                     // The manifest element provides an exhaustive list of the Publication Resources that constitute the EPUB Publication, each represented by an item element.
	Spine            Spine       `xml:"spine"`                        // The spine element defines the default reading order of the EPUB Publication content
	Bindings         []MediaType `xml:"mediaType>bindings,omitempty"` // The bindings element defines a set of custom handlers for media types not supported by this specification.
	Collection       *Collection `xml:"collection,omitempty"`         // The collection element defines a related group of resources. (Added in EPUB 301.)
}

// Manifest element provides an exhaustive list of the Publication Resources that constitute
// the EPUB Publication, each represented by an item element.
type Manifest struct {
	ID    string `xml:"id,attr,omitempty"` // The ID [XML] of this element, which must be unique within the document scope.
	Items []Item `xml:"item"`              // List of the Publication Resources
}

// Item element represents a Publication Resource.
type Item struct {
	ID           string `xml:"id,attr"`                      // The ID [XML] of this element, which must be unique within the document scope.
	Href         string `xml:"href,attr"`                    // An IRI [RFC3987] specifying the location of the Publication Resource described by this item.
	MediaType    string `xml:"media-type,attr"`              // A media type [RFC2046] that specifies the type and format of the Publication Resource described by this item.
	Fallback     string `xml:"fallback,attr,omitempty"`      // An IDREF [XML] that identifies the fallback for a non-Core Media Type.
	Properties   string `xml:"properties,attr,omitempty"`    // A space-separated list of property values.
	MediaOverlay string `xml:"media-overlay,attr,omitempty"` // An IDREF [XML] that identifies the Media Overlay Document for the resource described by this item.
}

// Spine element defines the default reading order of the EPUB Publication content by defining
// an ordered list of manifest item references.
type Spine struct {
	ID            string    `xml:"id,attr,omitempty"`                         // The ID [XML] of this element, which must be unique within the document scope.
	Toc           string    `xml:"toc,attr,omitempty"`                        // An IDREF [XML] that identifies the manifest item that represents the superseded NCX.
	PageDirection string    `xml:"page-progression-direction,attr,omitempty"` // The global direction in which the Publication content flows. Allowed values are ltr (left-to-right), rtl (right-to-left) and default.
	ItemRefs      []ItemRef `xml:"itemref"`                                   // Ordered subset of the Publication Resources listed in the manifest
}

// ItemRef elements of the spine represent a sequential list of Publication Resources
// (typically EPUB Content Documents). The order of the itemref elements defines the default
// reading order of the Publication.
type ItemRef struct {
	IDRef      string `xml:"idref,attr"`                // An IDREF [XML] that identifies a manifest item.
	Linear     string `xml:"linear,attr,omitempty"`     // Specifies whether the referenced content is primary. The value of the attribute must be yes or no. The default value is yes.
	ID         string `xml:"id,attr,omitempty"`         // The ID [XML] of this element, which must be unique within the document scope.
	Properties string `xml:"properties,attr,omitempty"` // A space-separated list of property values.
}

// MediaType element associates a Foreign Resource media type with a handler XHTML Content Document.
type MediaType struct {
	MediaType string `xml:"media-type,attr"` // A media type [RFC2046] that specifies the type and format of the resource to be handled.
	Handler   string `xml:"handler,attr"`    // An IDREF [XML] that identifies the manifest XHTML Content Document to be invoked to handle content of the type specified in this element
}

// Collection element defines a related group of resources.
type Collection struct {
	Lang        string       `xml:"xml:lang,attr,omitempty"` // Specifies the language used in the contents and attribute values of the carrying element and its descendants
	Dir         string       `xml:"dir,attr,omitempty"`      // Specifies the base text direction of the content and attribute values of the carrying element and its descendants.
	ID          string       `xml:"id,attr,omitempty"`       // The ID [XML] of this element, which must be unique within the document scope.
	Role        string       `xml:"role,attr"`               // Specifies the nature of the collection
	Metadata    *Metadata    `xml:"metadata,omitempty"`      // The optional metadata element child of collection is an adaptation of the package metadata element.
	Collections []Collection `xml:"collection,omitempty"`    // A collection may define sub-collections through the inclusion of one or more child collection elements.
	Links       []Link       `xml:"link,omitempty"`          // The link element child of collection is an adaptation of the metadata link element.
}
