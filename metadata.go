package epub

import (
	"fmt"
	"time"
)

// Element with optional ID.
type Element struct {
	Value string `xml:",chardata"`
	ID    string `xml:"id,attr,omitempty"` // The ID of this element, which must be unique within the document scope.
}

// ElementLang with optional ID, xml:lang & dir.
type ElementLang struct {
	Value string `xml:",chardata"`
	ID    string `xml:"id,attr,omitempty"`       // The ID of this element, which must be unique within the document scope.
	Dir   string `xml:"dir,attr,omitempty"`      // Specifies the base text direction of the content and attribute values of the carrying element and its descendants.
	Lang  string `xml:"xml:lang,attr,omitempty"` // Specifies the language used in the contents and attribute values of the carrying element and its descendants
}

// Metadata element encapsulates Publication meta information.
type Metadata struct {
	DC string `xml:"xmlns:dc,attr"` // “http://purl.org/dc/elements/1.1/”
	// Required Elements
	Identifier []Element     `xml:"dc:identifier"` // The [DCMES] identifier element contains a single identifier associated with the EPUB Publication, such as a UUID, DOI, ISBN or ISSN.
	Title      []ElementLang `xml:"dc:title"`      // The [DCMES] title element represents an instance of a name given to the EPUB Publication.
	Language   []Element     `xml:"dc:language"`   // The [DCMES] language element specifies the language of the Publication content.
	// DCMES Optional Elements
	Date        *Element      `xml:"dc:date,omitempty"`        // The date element must only be used to define the publication date of the EPUB Publication. The publication date is not the same as the last modified date (the last time the content was changed), which must be included using the [DCTERMS] modified property.
	Creator     []ElementLang `xml:"dc:creator,omitempty"`     // The creator element represents the name of a person, organization, etc. responsible for the creation of the content of a Publication. The role property can be attached to the element to indicate the function the creator played in the creation of the content.
	Contributor []ElementLang `xml:"dc:contributor,omitempty"` // The contributor element is used to represent the name of a person, organization, etc. that played a secondary role in the creation of the content of a Publication.
	Subject     []ElementLang `xml:"dc:subject,omitempty"`
	Description []ElementLang `xml:"dc:description,omitempty"`
	Type        []Element     `xml:"dc:type,omitempty"`
	Format      []Element     `xml:"dc:format,omitempty"`
	Publisher   []ElementLang `xml:"dc:publisher,omitempty"`
	Source      []Element     `xml:"dc:source,omitempty"`
	Relation    []ElementLang `xml:"dc:relation,omitempty"`
	Coverage    []ElementLang `xml:"dc:coverage,omitempty"`
	Rights      []ElementLang `xml:"dc:rights,omitempty"`
	// Meta
	Meta []Meta `xml:"meta,omitempty"` // The meta element provides a generic means of including package metadata, allowing the expression of primary metadata about the package or content and refinement of that metadata.
	Link []Link `xml:"link,omitempty"` // The link element is used to associate resources with a Publication, such as metadata records.
}

// AddTitle new publication title.
func (m *Metadata) AddTitle(name string) {
	m.Title = append(m.Title, ElementLang{Value: name})
}

// AddAuthor new publication author.
func (m *Metadata) AddAuthors(names ...string) {
	for _, name := range names {
		m.Creator = append(m.Creator, ElementLang{Value: name})
	}
}

// AddSubjects add publication subject.
func (m *Metadata) AddSubjects(names ...string) {
	for _, name := range names {
		m.Subject = append(m.Subject, ElementLang{Value: name})
	}
}

// SetDescription set publication description.
func (m *Metadata) SetDescription(description string) {
	m.Description = []ElementLang{{Value: description}}
}

// SetRights set publication rights description.
func (m *Metadata) SetRights(rights string) {
	m.Rights = []ElementLang{{Value: rights}}
}

// SetDate set publication date (not last modified).
func (m *Metadata) SetDate(date string) (err error) {
	// check data format
	var dateTime time.Time
	for _, layout := range []string{"2006-01-02", "2006-01", "2006", time.RFC3339} {
		if dateTime, err = time.Parse(layout, date); err == nil {
			break
		}
	}
	if dateTime.IsZero() {
		return fmt.Errorf("bad date %v", date)
	}

	// set publication date
	m.Date = &Element{Value: date}

	return nil
}

// SetUUID set publication identifier as UUID.
func (m *Metadata) SetUUID(id string) {
	if id == "" {
		id = NewUUID() // generate random UUID if not defined
	}
	m.Identifier = []Element{{Value: id, ID: "uuid"}}
}

// SetPublisher set publication publisher.
func (m *Metadata) SetPublisher(name string) {
	m.Publisher = []ElementLang{{Value: name}}
}

// SetLang set publication language.
func (m *Metadata) SetLang(lang string) {
	m.Language = []Element{{Value: lang}}
}

// Meta element provides a generic means of including package metadata, allowing the expression
// of primary metadata about the package or content and refinement of that metadata.
type Meta struct {
	Refines  string `xml:"refines,attr,omitempty"`  // Identifies the expression or resource augmented by this element. The value of the attribute must be a relative IRI [RFC3987] pointing to the resource or element it describes.
	Property string `xml:"property,attr"`           // A property. Refer to Vocabulary Association Mechanisms for more information.
	Scheme   string `xml:"scheme,attr,omitempty"`   // A property data type value indicating the source the value of the element is drawn from.
	ID       string `xml:"id,attr,omitempty"`       // The ID of this element, which must be unique within the document scope.
	Dir      string `xml:"dir,attr,omitempty"`      // Specifies the base text direction of the content and attribute values of the carrying element and its descendants.
	Lang     string `xml:"xml:lang,attr,omitempty"` // Specifies the language used in the contents and attribute values of the carrying element and its descendants
	Value    string `xml:",chardata"`
}

// Link element is used to associate resources with a Publication, such as metadata records.
type Link struct {
	Refines   string `xml:"refines,attr,omitempty"`    // Identifies the expression or resource augmented by this element. The value of the attribute must be a relative IRI [RFC3987] pointing to the resource or element it describes.
	Rel       string `xml:"rel,attr"`                  // A space-separated list of property values.
	Href      string `xml:"href,attr"`                 // An absolute or relative IRI reference [RFC3987] to a resource.
	ID        string `xml:"id,attr,omitempty"`         // The ID [XML] of this element, which must be unique within the document scope.
	MediaType string `xml:"media-type,attr,omitempty"` // A media type [RFC2046] that specifies the type and format of the resource referenced by this link.
}
