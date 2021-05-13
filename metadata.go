package epub

// Metadata element encapsulates Publication meta information.
type Metadata struct {
	DC string `xml:"xmlns:dc,attr"` // “http://purl.org/dc/elements/1.1/”
	// Required Elements
	Identifier []Element     `xml:"dc:identifier"` // The [DCMES] identifier element contains a single identifier associated with the EPUB Publication, such as a UUID, DOI, ISBN or ISSN.
	Title      []ElementLang `xml:"dc:title"`      // The [DCMES] title element represents an instance of a name given to the EPUB Publication.
	Language   []Element     `xml:"dc:language"`   // The [DCMES] language element specifies the language of the Publication content.
	// DCMES Optional Elements
	Creator     []ElementLang `xml:"dc:creator,omitempty"`     // The creator element represents the name of a person, organization, etc. responsible for the creation of the content of a Publication. The role property can be attached to the element to indicate the function the creator played in the creation of the content.
	Contributor []ElementLang `xml:"dc:contributor,omitempty"` // The contributor element is used to represent the name of a person, organization, etc. that played a secondary role in the creation of the content of a Publication.
	Date        *Element      `xml:"dc:date,omitempty"`        // The date element must only be used to define the publication date of the EPUB Publication. The publication date is not the same as the last modified date (the last time the content was changed), which must be included using the [DCTERMS] modified property.
	Coverage    []ElementLang `xml:"dc:coverage,omitempty"`
	Description []ElementLang `xml:"dc:description,omitempty"`
	Format      []Element     `xml:"dc:format,omitempty"`
	Publisher   []ElementLang `xml:"dc:publisher,omitempty"`
	Relation    []ElementLang `xml:"dc:relation,omitempty"`
	Rights      []ElementLang `xml:"dc:rights,omitempty"`
	Subject     []ElementLang `xml:"dc:subject,omitempty"`
	// Meta
	Meta []Meta `xml:"meta,omitempty"` // The meta element provides a generic means of including package metadata, allowing the expression of primary metadata about the package or content and refinement of that metadata.
	Link []Link `xml:"link,omitempty"` // The link element is used to associate resources with a Publication, such as metadata records.
}

// Element with optional ID.
type Element struct {
	ID    string `xml:"id,attr,omitempty"` // The ID of this element, which must be unique within the document scope.
	Value string `xml:",chardata"`
}

// ElementLang with optional ID, xml:lang & dir.
type ElementLang struct {
	ID    string `xml:"id,attr,omitempty"` // The ID of this element, which must be unique within the document scope.
	Value string `xml:",chardata"`
	Lang  string `xml:"xml:lang,attr,omitempty"` // Specifies the language used in the contents and attribute values of the carrying element and its descendants
	Dir   string `xml:"dir,attr,omitempty"`      // Specifies the base text direction of the content and attribute values of the carrying element and its descendants.
}
