package epub

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Default metadata settings.
var (
	DefaultLang  = Element{Value: "en"}
	DefaultTitle = ElementLang{Value: "Untitled"}
)

// Writer allows you to create publications in epub 3 format.
type Writer struct {
	Metadata
	zipWriter *zip.Writer
	manifest  []Item
	spine     []ItemRef
	counter   uint
}

// New return new epub publication Writer.
func New(w io.Writer) (wr *Writer, err error) {
	zipWriter := zip.NewWriter(w) // create zip-compressor
	// close zip writer on error
	defer func() {
		if err != nil {
			zipWriter.Close()
		}
	}()

	// write mimetype header
	item, err := zipWriter.CreateHeader(&zip.FileHeader{
		Name:   "mimetype",
		Method: zip.Store,
	})
	if err != nil {
		return nil, err
	}
	if _, err = io.WriteString(item, "application/epub+zip"); err != nil {
		return nil, err
	}

	// write container file
	if err = addXMLData(zipWriter,
		"META-INF/container.xml",
		Container{
			Version: "1.0",
			Rootfiles: []RootFile{
				{
					FullPath:  path.Join(RootPath, PackageFilename),
					MediaType: "application/oebps-package+xml",
				},
			},
		}); err != nil {
		return nil, err
	}

	// return initializer Writer
	return &Writer{
		zipWriter: zipWriter,
		manifest:  make([]Item, 0, 20),
		spine:     make([]ItemRef, 0, 20),
	}, nil
}

// ContentType describe type of content file.
type ContentType byte

// Supported types of content file.
const (
	Primary   ContentType = iota // Primary content file
	Auxiliary                    // Auxiliary content file
	Media                        // Media file
)

// AddContent adds data to the publication.
func (w *Writer) AddContent(r io.Reader, name string, ct ContentType, properties ...string) error {
	name = filepath.ToSlash(name) // normalize file name

	// check if already added
	for _, item := range w.manifest {
		if item.Href == name {
			return fmt.Errorf("a file with the name %q has already been added"+
				" to the publication", name)
		}
	}

	// generate file id and add to manifest
	w.counter++
	id := fmt.Sprintf("id%02x", w.counter)
	w.manifest = append(w.manifest, Item{
		ID:         id,
		Href:       name,
		MediaType:  typeByName(name),
		Properties: strings.Join(properties, " "),
	})

	// if it content file than add to spine
	if ct < Media {
		itemref := ItemRef{IDRef: id}
		if ct == Auxiliary {
			itemref.Linear = "no"
		}
		w.spine = append(w.spine, itemref)
	}

	// write file to publication
	file, err := w.zipWriter.Create(path.Join(RootPath, name))
	if err != nil {
		return err
	}
	_, err = io.Copy(file, r)

	return err
}

// Close closes the publication and writes metadata.
func (w *Writer) Close() error {
	metadata := w.Metadata // copy metadata
	// add DC namespace if not defined
	if metadata.DC == "" {
		metadata.DC = "http://purl.org/dc/elements/1.1/"
	}

	// set global publication UID
	var uid string
	for _, item := range metadata.Identifier {
		if item.ID != "" {
			uid = item.ID
			break
		}
	}
	if uid == "" {
		// UID not defined
		metadata.Identifier = append(metadata.Identifier,
			Element{ID: "uuid", Value: NewUUID()})
		uid = "uuid"
	}

	// set modified time
	var setTime bool
	for _, item := range metadata.Meta {
		if item.Property == "dcterms:modified" {
			item.Value = now()
			setTime = true
			break
		}
	}
	if !setTime {
		// modified time not set
		metadata.Meta = append(metadata.Meta, Meta{
			Property: "dcterms:modified",
			Value:    now(),
		})
	}

	// add default language if not defined
	if len(metadata.Language) == 0 {
		metadata.Language = []Element{DefaultLang}
	}

	// add default title if not defined
	if len(metadata.Title) == 0 {
		metadata.Title = []ElementLang{DefaultTitle}
	}

	// create & write publication package file
	if err := addXMLData(w.zipWriter,
		path.Join(RootPath, PackageFilename),
		Package{
			Version:          "3.0",
			UniqueIdentifier: uid,
			Metadata:         metadata,
			Manifest: Manifest{
				Items: w.manifest,
			},
			Spine: Spine{
				ItemRefs: w.spine,
			},
		}); err != nil {
		w.zipWriter.Close() // close zip writer on error
		return err
	}

	// close publication
	return w.zipWriter.Close()
}

// now return string wih current time i RFC 3339 format.
func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// addXMLData serialize & write publication data as XML file.
func addXMLData(w *zip.Writer, name string, data interface{}) error {
	// create new publication file
	item, err := w.Create(name)
	if err != nil {
		return err
	}

	// add XML header
	if _, err := io.WriteString(item, xml.Header); err != nil {
		return err
	}

	// serialize XML data to file
	enc := xml.NewEncoder(item)
	enc.Indent("", "\t")
	return enc.Encode(data)
}
