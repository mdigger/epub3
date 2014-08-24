package epub

import (
	"encoding/xml"
)

// An EPUB provides almost all of this fundamental information in an XML file called the package
// document. This contains that invaluable packing list and those indispensable assembly
// instructions that enable a reading system to know what it has and what to do with it.
//
// The root element of the package document is the package element. This, in turn, contains the
// metadata and resource information in its child elements, in this order:
//  • metadata (required)
//  • manifest (required)
//  • spine (required)
//  • guide (optional and deprecated; a carryover from EPUB 2)
//  • bindings (optional)
type Package struct {
	XMLName    xml.Name   `xml:"http://www.idpf.org/2007/opf package"`
	Version    string     `xml:"version,attr"` // EPUB 3s must declare "3.0"
	LangAndDir            // Language & reading direction
	UID        string     `xml:"unique-identifier,attr"`  // Unique identifier ID
	Metadata   *Metadata  `xml:"metadata"`                // Most of the metadata in a typical EPUB is associated with the publication as a whole.
	Manifest   []*Item    `xml:"manifest>item"`           // The manifest documents all of the individual resources that together constitute the EPUB
	Spine      []*ItemRef `xml:"spine>itemref"`           // The spine provides a default reading order by which those resources may be presented to a user
	Bindings   []*Link    `xml:"bindings>link,omitempty"` // Contain fallbacks that are more sophisticated than those provided by the HTML5 object element’s fallback mechanisms
}

// Most of the metadata in a typical EPUB is associated with the publication as a whole.
//
// This is intended to tell a reading system, when it opens up the EPUB, everything it needs to
// know about what’s inside. Which EPUB is this (identifiers)? What names is it known by (titles)?
// Does it use any vocabularies I don’t necessarily understand (prefixes)? What language does it
// use? What are all the things in the box (manifest)? Which one is the cover image, and do any of
// them contain MathML or SVG or scripting (spine itemref properties)? In what order should I
// present the content (spine), and how can a user navigate this EPUB (the nav document)? Are there
// resources I need to link to (link)? Are there any media objects I’m not designed by default to
// handle (bindings)?
//
// Having all of this information up-front in the EPUB makes things much easier for a reading
// system, rather than requiring it to simply discover that unrecognized vocabulary, or that
// MathML buried deep in a content document, only when it comes across it, as a browser does with
// a normal website.
//
// The metadata element contains the same three required elements as it did in EPUB 2, one new
// required element, and a number of optional elements, including that all-powerful meta element
// described previously.
//
// As mentioned earlier, EPUB continues to use the Dublin Core Metadata Element Set (DCMES) for
// most of its required and optional metadata.
//
// XML rules require that you declare the Dublin Core namespace in order to use the elements. This
// declaration is typically added to the metadata element, but can also be added to the root
// package element. For example:
//  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
//    <dc:identifier id="pub-identifier">urn:isbn:9781449325299</dc:identifier>
//    <dc:title id="pub-title">EPUB 3 Best Practices</dc:title>
//    <dc:language id="pub-language">en</dc:language>
//  </metadata>
type Metadata struct {
	DC          string          `xml:"xmlns:dc,attr"`            // Dublin Core namespace - "http://purl.org/dc/elements/1.1/"
	Identifier  []*MetaProperty `xml:"dc:identifier"`            // Identifier for the publication
	Title       []*MetaProperty `xml:"dc:title"`                 // Title for the publication
	Language    []*MetaProperty `xml:"dc:language"`              // Specifies the language of the publication’s content
	Creator     []*MetaProperty `xml:"dc:creator,omitempty"`     // The name of a person or organization with primary responsibility for creating the content, such as an author
	Contributor []*MetaProperty `xml:"dc:contributor,omitempty"` // Is used in the same way, but indicates a secondary level of involvement (for example, a translator or an illustrator)
	Date        string          `xml:"dc:date,omitempty"`        // Provide the date of the EPUB publication in the standard W3C date and time format (2000-01-01T00:00:00Z)
	Source      *MetaProperty   `xml:"dc:source,omitempty"`      // Identifier of the source publication from which the EPUB was derived, such as the print version
	Type        *MetaProperty   `xml:"dc:type,omitempty"`        // Presents a bit of a curveball at the moment, because the IDPF has not yet defined values for it
	Coverage    []*MetaProperty `xml:"dc:coverage,omitempty"`
	Description []*MetaProperty `xml:"dc:description,omitempty"`
	Format      []*MetaProperty `xml:"dc:format,omitempty"`
	Publisher   []*MetaProperty `xml:"dc:publisher,omitempty"`
	Relation    []*MetaProperty `xml:"dc:relation,omitempty"`
	Rights      []*MetaProperty `xml:"dc:rights,omitempty"`
	Subject     []*MetaProperty `xml:"dc:Subject,omitempty"`
	Meta        []*Meta         `xml:"meta"` // The metadata element can also contain any number of those useful meta elements
}

// Contains an identifier and value element for Metadata.
type MetaProperty struct {
	LangAndDir        // Language & reading direction
	Id         string `xml:"id,attr,omitempty"`
	Value      string `xml:",chardata"`
}

// The workhorse of EPUB 3 metadata is the meta element, which provides a simple, generic, and yet
// surprisingly flexible and powerful mechanism for associating metadata of virtually unlimited
// richness with the EPUB package and its contents. An EPUB can have any number of meta elements.
// They’re contained in the metadata element, the first child of the package element, and from that
// central location they serve as a hub for metadata about the EPUB, its resources, its content
// documents, and even locations within the content documents.
//
// The meta element uses the refines attribute to specify what it applies to, using an ID in the
// form of a relative IRI.
//
// When the refines attribute is not provided, it is assumed that the meta element applies to the
// package as a whole; this is referred to as a primary expression. When the meta element does have
// a refines attribute, it is called a subexpression.
//
// Each meta has a property attribute that defines what kind of statement is being made in the
// text of the meta element. The values of property can be the default vocabulary, a term from one
// of the reserved vocabularies, or a term from one of the vocabularies defined via the prefix
// mechanism.
//
// The default vocabulary for meta consists of the following property values:
//
// • "alternate-script"
// — typically used to provide versions of titles and the names of authors or contributors in a
// language and script identified by the xml:lang attribute:
//   <meta refines="#creator" property="alternate-script" xml:lang="ja">村上 春樹</meta>
//
// • "display-seq"
// — used to specify the sequence in which multiple versions of the same thing—for example, multiple
// forms of the title—should be displayed:
//  <meta refines="#title2" property="display-seq">1</meta>
//
// • "file-as"
// — provides an alternate version—again, typically of a title or the name of an author or other
// contributor—in a form that will alphabetize properly, e.g., last-name-first for an author’s name
// or putting “The” at the end of a title that begins with it:
//  <meta refines="#creator" property="file-as">Murakami, Haruki</meta>
//
// • "group-position"
// — specifies the position of the referenced item in relation to others that it is grouped with.
// This is useful, for example, so that all the titles in a series are displayed in proper order in
// a reader’s bookshelf:
//  <meta refines="#title3" property="group-position">2</meta>
//
// • "identifier-type"
// — provides a way to distinguish between different types of identifiers (e.g., ISBN versus DOI).
// Its values can be drawn from an authority like the ONIX Code List 5, which is specified with the
// scheme attribute:
//  <meta refines="#src-id" property="identifier-type" scheme="onix:codelist5">15</meta>
//
// • "meta-auth"
// — documents the “metadata authority” responsible for a given instance of metadata:
//  <meta refines="isbn-id" property="meta-auth">isbn-international.org</meta>
//
// • "role"
// — most often used to specify the exact role performed by a contributor—for example, a translator
// or illustrator:
//  <meta refines="#creator" property="role" scheme="marc:relators">ill</meta>
//
// • "title-type"
// — distinguishes six specific forms of titles:
//  <meta refines="#title" property="title-type">subtitle</meta>
//
// A meta element may also have an ID of its own, as the value of the id attribute:
//  <meta refines="isbn-id" property="meta-auth" id="meta-auth">isbn-international.org</meta>
// This ID can be used to make metadata chains, where one meta refines another. The element may
// also have a formal identifier of the scheme used for the value of the property (using the scheme
// attribute).
//
// You can also use property values, which must include the proper prefix, from any of the reserved
// vocabularies or any vocabulary for which you’ve declared the prefix:
//  <meta property="dcterms:dateCopyrighted">2012</meta>
// You’ll notice that the previous example did not include a refines attribute. This was
// intentional, as the other use for the meta element is to define metadata for the publication as
// a whole. We’ll look at the Dublin Core elements for publication metadata shortly, but you are
// not limited to using them. If another vocabulary provides richer metadata, you can use the meta
// element to express it.
//
// There is one very specific use of the meta element that is quite important; in fact, it is
// a requirement for EPUB 3. The meta element is used to provide a timestamp that records the
// modification date on which the EPUB was created. It uses the dcterms:modified property and
// requires a value conforming to the W3C dateTime form, like this:
//  <meta property="dcterms:modified">2011-01-01T12:00:00Z</meta>
type Meta struct {
	MetaProperty
	Refines  string `xml:"refines,attr,omitempty"`
	Property string `xml:"property,attr,omitempty"`
}

// Additional XML attributes: language & reading direction
type LangAndDir struct {
	Lang string `xml:"xml:lang,attr,omitempty"` // The language of the package document
	Dir  string `xml:"dir,attr,omitempty"`      // The text directionality of the package document: left-to-right (ltr) or right-to-left (rtl)
}

// Each and every resource that is part of the EPUB — every content document, every image, every
// video and audio file, every font, every style sheet: every individual resource — is documented
// by an item element in the manifest. The purpose is to alert a reading system, up front, about
// everything it must find in the publication, what kind of media each thing is, and where it can
// find it. They can be in any order, but they all have to be in the manifest.
//
// Each item contains three required attributes: id, href & media-type.
//  <item id="chapter01" href="xhtml/c01.xhtml" media-type="application/xhtml+xml"/>
//
// An item may also have a properties attribute with one or more space-separated property values
// that alert the reading system to some useful information about the item. The values for manifest
// property values in EPUB 3 are:
//
// • "cover-image"
// — clearly documenting which item should be displayed as the cover.
//
// • "mathml", "scripted", "svg", "remote-resources", and "switch"
// — alerting the reading system to where it will have to deal with MathML, JavaScript, SVG,
// remote resources, or non-EPUB XML fragments (see the other chapters for more detail on these).
//
// • "nav"
// — the XHTML5 document that is the required Navigation Document
//
// And finally, an item may have a media-overlay attribute, the value of which is an IDREF of the
// Media Overlay Document for the item. Media Overlays are the EPUB 3 mechanism for synchronizing
// text and recorded audio.
type Item struct {
	Id           string `xml:"id,attr"`
	Href         string `xml:"href,attr"`
	MediaType    string `xml:"media-type,attr"`
	Properties   string `xml:"properties,attr,omitempty"`
	MediaOverlay string `xml:"media-overlay,attr,omitempty"`
	Fallback     string `xml:"fallback,attr,omitempty"`
}

// Whereas the manifest documents each and every item in the EPUB, in no particular order, the
// spine provides a default reading order, and it is required to list only those components that
// are not referenced by other components (primary content). The point is to provide at least one
// path by which everything in the EPUB will be presented to the reader, in at least one logical order.
type ItemRef struct {
	IdRef      string `xml:"idref,attr"`
	Linear     string `xml:"linear,attr,omitempty"`
	Properties string `xml:"properties,attr,omitempty"`
}

// The link element requires an href attribute to provide either an absolute or relative IRI to a
// resource, and a rel attribute to provide the property value—i.e., what kind of resource is being
// linked to. The values defined for the rel attribute in EPUB 3.0 are:
//
// • "marc21xml-record"
// — for a MARC21 record providing bibliographic metadata for the publication
//
// • "mods-record"
// — for a MODS record of the publication conforming to the Library of Congress’s Metadata Object
// Description Schema
//
// • "onix-record"
// — for an ONIX record providing book supply chain metadata for the publication conforming to
// EDItEUR’s ONIX for Books specification
//
// • "xml-signature"
// — for an XML Signature applying to the publication or an associated property conforming to the
// W3C’s XML Signature specification
//
// • "xmp-record"
// — for an XMP record conforming to the ISO Extensible Metadata Platform that applies to the
// publication (not just a component, like an image, for which the prefix mechanism and meta
// element should be used to provide metadata using the xmp: prefix)
//
// Here is how you might reference an external ONIX record:
//  <link rel="onix-record" href="http://example.org/meta/records/onix/121099"/>
type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}
