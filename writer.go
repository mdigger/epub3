package epub

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// ContentType describe type of content file.
type ContentType byte

// Supported types of content file.
const (
	Primary   ContentType = iota // Primary content file
	Auxiliary                    // Auxiliary content file
	Media                        // Media file
)

// Writer allows you to create publications in epub 3 format.
type Writer struct {
	*Metadata
	zipWriter *zip.Writer
	manifest  []*Item
	spine     []*ItemRef
	counter   uint
}

// New return new epub publication Writer.
func New(w io.Writer) (*Writer, error) {
	zipWriter := zip.NewWriter(w) // create zip-compressor

	// write mimetype header
	item, err := zipWriter.CreateHeader(&zip.FileHeader{
		Name:   "mimetype",
		Method: zip.Store,
	})
	if err != nil {
		return nil, err
	}
	if _, err = io.WriteString(item, EPUBMimeType); err != nil {
		return nil, err
	}

	// write container file
	item, err = zipWriter.Create(path.Join(METAINF, CONTAINER))
	if err != nil {
		return nil, err
	}
	if _, err = io.WriteString(item, xml.Header); err != nil {
		return nil, err
	}
	enc := xml.NewEncoder(item)
	enc.Indent("", "\t")
	err = enc.Encode(DefaultContainer)
	if err != nil {
		return nil, err
	}

	// return initializer Writer
	writer := &Writer{
		zipWriter: zipWriter,
		manifest:  make([]*Item, 0, 10),
		spine:     make([]*ItemRef, 0, 10),
		Metadata:  new(Metadata),
	}
	return writer, nil
}

// AddContentFile adds the source file to the publication.
func (w *Writer) AddContentFile(sourceFilename, filename string, ct ContentType,
	properties ...string) error {
	file, err := os.Open(sourceFilename)
	if err != nil {
		return err
	}
	defer file.Close()
	return w.AddContent(filename, ct, file, properties...)
}

// AddContent adds data to the publication.
func (w *Writer) AddContent(filename string, ct ContentType, r io.Reader,
	properties ...string) error {
	filename = filepath.ToSlash(filename)
	for _, item := range w.manifest {
		if item.Href == filename {
			return fmt.Errorf("a file with the name %q has already been added"+
				" to the publication", filename)
		}
	}
	w.counter++
	id := fmt.Sprintf("id%02x", w.counter)
	item := &Item{
		ID:         id,
		Href:       filename,
		MediaType:  TypeByFilename(filename),
		Properties: strings.Join(properties, " "),
	}
	w.manifest = append(w.manifest, item)
	if ct < Media {
		itemref := &ItemRef{IDRef: id}
		if ct == Auxiliary {
			itemref.Linear = "no"
		}
		w.spine = append(w.spine, itemref)
	}
	fileWriter, err := w.zipWriter.Create(path.Join(RootPath, filename))
	if err != nil {
		return err
	}
	_, err = io.Copy(fileWriter, r)
	return err
}

// Close closes the publication and writes metadata.
func (w *Writer) Close() (err error) {
	metadata := w.Metadata
	if metadata == nil {
		metadata = new(Metadata)
	}
	if metadata.DC == "" {
		metadata.DC = "http://purl.org/dc/elements/1.1/"
	}

	var uid string
	for _, item := range metadata.Identifier {
		if item.ID != "" {
			uid = item.ID
			break
		}
	}
	if uid == "" {
		metadata.AddMeta("uuid", "uuid", "urn:uuid:"+NewUUID())
		uid = "uuid"
	}

	var setTime bool
	for _, item := range metadata.Meta {
		if item.Property == "dcterms:modified" {
			item.Value = time.Now().UTC().Format(time.RFC3339)
			setTime = true
			break
		}
	}
	if !setTime {
		if metadata.Meta == nil {
			metadata.Meta = make([]*Meta, 0, 1)
		}
		metadata.Meta = append(metadata.Meta, &Meta{
			Property: "dcterms:modified",
			Value:    time.Now().UTC().Format(time.RFC3339),
		})
	}

	if len(metadata.Language) == 0 {
		metadata.Language.Add("", "en")
	}

	if len(metadata.Title) == 0 {
		metadata.Title.Add("", "Untitled")
	}

	item, err := w.zipWriter.Create(path.Join(RootPath, PackageFilename))
	if err != nil {
		return err
	}
	if _, err := io.WriteString(item, xml.Header); err != nil {
		return err
	}
	opf := &Package{
		Version:          "3.0",
		UniqueIdentifier: uid,
		Metadata:         metadata,
		Manifest: &Manifest{
			Items: w.manifest,
		},
		Spine: &Spine{
			ItemRefs: w.spine,
		},
	}
	enc := xml.NewEncoder(item)
	enc.Indent("", "\t")
	if err = enc.Encode(opf); err != nil {
		return err
	}

	return w.zipWriter.Close()
}
