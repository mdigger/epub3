package epub

import (
	"encoding/xml"
)

// The package element is the root container of the Package Document and encapsulates Publication
// metadata and resource information.
type Package struct {
	XMLName          xml.Name    `xml:"http://www.idpf.org/2007/opf package"`
	Version          string      `xml:"version,attr"`            // Specifies the EPUB specification version to which the Publication conforms
	UniqueIdentifier string      `xml:"unique-identifier,attr"`  // An IDREF that identifies the dc:identifier element that provides the package's preferred, or primary, identifier
	Prefix           string      `xml:"prefix,attr,omitempty"`   // Declaration mechanism for prefixes not reserved by this specification.
	Lang             string      `xml:"xml:lang,attr,omitempty"` // Specifies the language used in the contents and attribute values of the carrying element and its descendants
	Dir              string      `xml:"dir,attr,omitempty"`      // Specifies the base text direction of the content and attribute values of the carrying element and its descendants.
	Id               string      `xml:"id,attr,omitempty"`       // The ID of this element, which must be unique within the document scope
	Metadata         *Metadata   `xml:"metadata"`                // The metadata element encapsulates Publication meta information
	Manifest         *Manifest   `xml:"manifest"`                // The manifest element provides an exhaustive list of the Publication Resources that constitute the EPUB Publication, each represented by an item element.
	Spine            *Spine      `xml:"spine"`                   // The spine element defines the default reading order of the EPUB Publication content
	Bindings         *Bindings   `xml:"bindings,omitempty"`      // The bindings element defines a set of custom handlers for media types not supported by this specification.
	Collection       *Collection `xml:"collection,omitempty"`    // The collection element defines a related group of resources.
}

// The metadata element encapsulates Publication meta information.
type Metadata struct {
	DC string `xml:"xmlns:dc,attr"` // “http://purl.org/dc/elements/1.1/”
	// Requered Elements
	Identifier Elements     `xml:"dc:identifier"` // The [DCMES] identifier element contains a single identifier associated with the EPUB Publication, such as a UUID, DOI, ISBN or ISSN.
	Title      LangElements `xml:"dc:title"`      // The [DCMES] title element represents an instance of a name given to the EPUB Publication.
	Language   Elements     `xml:"dc:language"`   // The [DCMES] language element specifies the language of the Publication content.
	// DCMES Optional Elements
	Creator     LangElements `xml:"dc:creator"`     // The creator element represents the name of a person, organization, etc. responsible for the creation of the content of a Publication. The role property can be attached to the element to indicate the function the creator played in the creation of the content.
	Contributor LangElements `xml:"dc:contributor"` // The contributor element is used to represent the name of a person, organization, etc. that played a secondary role in the creation of the content of a Publication.
	Date        *Element     `xml:"dc:date"`        // The date element must only be used to define the publication date of the EPUB Publication. The publication date is not the same as the last modified date (the last time the content was changed), which must be included using the [DCTERMS] modified property.
	Source      *Element     `xml:"dc:source"`      // The source element must only be used to specify the identifier of the source publication from which this EPUB Publication is derived.
	Type        *Element     `xml:"dc:type"`        // The type element is used to indicate that the given Publication is of a specialized type (e.g., annotations packaged in EPUB format or a dictionary).
	Coverage    LangElements `xml:"dc:coverage"`
	Description LangElements `xml:"dc:description"`
	Format      Elements     `xml:"dc:format"`
	Publisher   LangElements `xml:"dc:publisher"`
	Relation    LangElements `xml:"dc:relation"`
	Rights      LangElements `xml:"dc:rights"`
	Subject     LangElements `xml:"dc:subject"`
	// Meta
	Meta []*Meta `xml:"meta"` // The meta element provides a generic means of including package metadata, allowing the expression of primary metadata about the package or content and refinement of that metadata.
	Link []*Link `xml:"link"` // The link element is used to associate resources with a Publication, such as metadata records.
}

func (self *Metadata) Add(name, id, value string) {
	switch name {
	case "identifier", "id", "uid", "pub-id", "UUID", "DOI", "ISBN", "ISSN":
		self.Identifier.Add(id, value)
	case "title":
		self.Title.Add(id, value)
	case "language", "lang":
		self.Language.Add(id, value)
	case "creator", "author":
		self.Creator.Add(id, value)
	case "contributor":
		self.Contributor.Add(id, value)
	case "date", "created":
		self.Date = &Element{Id: id, Value: value}
	case "source":
		self.Source = &Element{Id: id, Value: value}
	case "type":
		self.Type = &Element{Id: id, Value: value}
	case "coverage":
		self.Coverage.Add(id, value)
	case "description":
		self.Description.Add(id, value)
	case "format":
		self.Format.Add(id, value)
	case "publisher":
		self.Publisher.Add(id, value)
	case "relation":
		self.Relation.Add(id, value)
	case "rights":
		self.Rights.Add(id, value)
	case "subject":
		self.Subject.Add(id, value)
	}
}

func CreateMetadata(metadata map[string]string) *Metadata {
	data := &Metadata{
		DC: "http://purl.org/dc/elements/1.1/",
	}
	for item, value := range metadata {
		data.Add(item, item, value)
	}
	return data
}

// Element with optional ID
type Element struct {
	Id    string `xml:"id,attr,omitempty"` // The ID of this element, which must be unique within the document scope.
	Value string `xml:",chardata"`
}

type Elements []*Element

func (self *Elements) Add(id, value string) {
	if self == nil {
		elements := make(Elements, 0)
		self = &elements
	}
	*self = append(*self, &Element{Id: id, Value: value})
}

// Element with optional ID, xml:lang & dir
type LangElement struct {
	Id    string `xml:"id,attr,omitempty"`       // The ID of this element, which must be unique within the document scope.
	Lang  string `xml:"xml:lang,attr,omitempty"` // Specifies the language used in the contents and attribute values of the carrying element and its descendants
	Dir   string `xml:"dir,attr,omitempty"`      // Specifies the base text direction of the content and attribute values of the carrying element and its descendants.
	Value string `xml:",chardata"`
}

type LangElements []*LangElement

func (self *LangElements) Add(id, value string) {
	if self == nil {
		elements := make(LangElements, 0)
		self = &elements
	}
	*self = append(*self, &LangElement{Id: id, Value: value})
}

// The meta element provides a generic means of including package metadata, allowing the expression
// of primary metadata about the package or content and refinement of that metadata.
type Meta struct {
	Property string `xml:"property,attr"`           // A property. Refer to Vocabulary Association Mechanisms for more information.
	Refines  string `xml:"refines,attr,omitempty"`  // Identifies the expression or resource augmented by this element. The value of the attribute must be a relative IRI [RFC3987] pointing to the resource or element it describes.
	Id       string `xml:"id,attr,omitempty"`       // The ID of this element, which must be unique within the document scope.
	Scheme   string `xml:"scheme,attr,omitempty"`   // A property data type value indicating the source the value of the element is drawn from.
	Lang     string `xml:"xml:lang,attr,omitempty"` // Specifies the language used in the contents and attribute values of the carrying element and its descendants
	Dir      string `xml:"dir,attr,omitempty"`      // Specifies the base text direction of the content and attribute values of the carrying element and its descendants.
	Value    string `xml:",chardata"`
}

// The link element is used to associate resources with a Publication, such as metadata records.
type Link struct {
	Href      string `xml:"href,attr"`                 // An absolute or relative IRI reference [RFC3987] to a resource.
	Rel       string `xml:"rel,attr"`                  // A space-separated list of property values.
	Id        string `xml:"id,attr,omitempty"`         // The ID [XML] of this element, which must be unique within the document scope.
	Refines   string `xml:"refines,attr,omitempty"`    // Identifies the expression or resource augmented by this element. The value of the attribute must be a relative IRI [RFC3987] pointing to the resource or element it describes.
	MediaType string `xml:"media-type,attr,omitempty"` // A media type [RFC2046] that specifies the type and format of the resource referenced by this link.
}

// The manifest element provides an exhaustive list of the Publication Resources that constitute
// the EPUB Publication, each represented by an item element.
type Manifest struct {
	Id    string  `xml:"id,attr,omitempty"` // The ID [XML] of this element, which must be unique within the document scope.
	Items []*Item `xml:"item"`              // List of the Publication Resources
}

// The item element represents a Publication Resource.
type Item struct {
	Id           string `xml:"id,attr"`                      // The ID [XML] of this element, which must be unique within the document scope.
	Href         string `xml:"href,attr"`                    // An IRI [RFC3987] specifying the location of the Publication Resource described by this item.
	MediaType    string `xml:"media-type,attr"`              // A media type [RFC2046] that specifies the type and format of the Publication Resource described by this item.
	Fallback     string `xml:"fallback,attr,omitempty"`      // An IDREF [XML] that identifies the fallback for a non-Core Media Type.
	Properties   string `xml:"properties,attr,omitempty"`    // A space-separated list of property values.
	MediaOverlay string `xml:"media-overlay,attr,omitempty"` // An IDREF [XML] that identifies the Media Overlay Document for the resource described by this item.
}

// The spine element defines the default reading order of the EPUB Publication content by defining
// an ordered list of manifest item references.
type Spine struct {
	Id            string     `xml:"id,attr,omitempty"`                         // The ID [XML] of this element, which must be unique within the document scope.
	Toc           string     `xml:"toc,attr,omitempty"`                        // An IDREF [XML] that identifies the manifest item that represents the superseded NCX.
	PageDirection string     `xml:"page-progression-direction,attr,omitempty"` // The global direction in which the Publication content flows. Allowed values are ltr (left-to-right), rtl (right-to-left) and default.
	ItemRefs      []*ItemRef `xml:"itemref"`                                   // Ordered subset of the Publication Resources listed in the manifest
}

// The child itemref elements of the spine represent a sequential list of Publication Resources
// (typically EPUB Content Documents). The order of the itemref elements defines the default
// reading order of the Publication.
type ItemRef struct {
	IdRef      string `xml:"idref,attr"`                // An IDREF [XML] that identifies a manifest item.
	Linear     string `xml:"linear,attr,omitempty"`     // Specifies whether the referenced content is primary. The value of the attribute must be yes or no. The default value is yes.
	Id         string `xml:"id,attr,omitempty"`         // The ID [XML] of this element, which must be unique within the document scope.
	Properties string `xml:"properties,attr,omitempty"` // A space-separated list of property values.
}

// The bindings element defines a set of custom handlers for media types not supported by this
// specification.
type Bindings struct {
	MediaTypes []*MediaType `xml:"mediaType"`
}

// The mediaType element associates a Foreign Resource media type with a handler XHTML Content Document.
type MediaType struct {
	MediaType string `xml:"media-type,attr"` // A media type [RFC2046] that specifies the type and format of the resource to be handled.
	Handler   string `xml:"handler,attr"`    // An IDREF [XML] that identifies the manifest XHTML Content Document to be invoked to handle content of the type specified in this element
}

// The collection element defines a related group of resources.
type Collection struct {
	Lang        string        `xml:"xml:lang,attr,omitempty"` // Specifies the language used in the contents and attribute values of the carrying element and its descendants
	Dir         string        `xml:"dir,attr,omitempty"`      // Specifies the base text direction of the content and attribute values of the carrying element and its descendants.
	Id          string        `xml:"id,attr,omitempty"`       // The ID [XML] of this element, which must be unique within the document scope.
	Role        string        `xml:"role,attr"`               // Specifies the nature of the collection
	Metadata    *Metadata     `xml:"metadata,omitempty"`      // The optional metadata element child of collection is an adaptation of the package metadata element.
	Collections []*Collection `xml:"collection,omitempty"`    // A collection may define sub-collections through the inclusion of one or more child collection elements.
	Links       []*Link       `xml:"link,omitempty"`          // The link element child of collection is an adaptation of the metadata link element.
}
