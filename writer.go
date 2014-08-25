package epub

import (
	"archive/zip"
	"code.google.com/p/go-uuid/uuid"
	"encoding/xml"
	"fmt"
	"github.com/mdigger/commitfile"
	"io"
	"mime"
	"path"
	"path/filepath"
	"time"
)

var (
	RootPath        = "OEBPS"
	EPUBMimeType    = "application/epub+zip"
	PackageFilename = "package.opf"
)

type Writer struct {
	file      *commitfile.File
	zipWriter *zip.Writer
	metadata  *Metadata
	manifest  []*Item
	spine     []*ItemRef
	counter   uint
}

func Create(filename string, metadata *Metadata) (writer *Writer, err error) {
	file, err := commitfile.Create(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			file.Close()
		}
	}()
	zipWriter := zip.NewWriter(file)
	var item io.Writer
	item, err = zipWriter.CreateHeader(&zip.FileHeader{
		Name:   "mimetype",
		Method: zip.Store,
	})
	if err != nil {
		return nil, err
	}
	if _, err = io.WriteString(item, EPUBMimeType); err != nil {
		return nil, err
	}
	item, err = zipWriter.Create(path.Join(METAINF, CONTAINER))
	if err != nil {
		return nil, err
	}
	if _, err = io.WriteString(item, xml.Header); err != nil {
		return nil, err
	}
	enc := xml.NewEncoder(item)
	enc.Indent("", "\t")
	err = enc.Encode(&Container{
		Version: "1.0",
		Rootfiles: []*RootFile{
			&RootFile{
				FullPath:  path.Join(RootPath, PackageFilename),
				MediaType: "application/oebps-package+xml",
			},
		},
	})
	if err != nil {
		return nil, err
	}
	writer = &Writer{
		file:      file,
		zipWriter: zipWriter,
		metadata:  metadata,
		manifest:  make([]*Item, 0, 10),
		spine:     make([]*ItemRef, 0, 10),
	}
	return writer, nil
}

func (self *Writer) Add(filename string, spine bool) (io.Writer, error) {
	filename = filepath.ToSlash(filename)
	var mimetype string
	switch ext := path.Ext(filename); ext {
	case ".gif":
		mimetype = "image/gif"
	case ".jpg", ".jpeg", ".jpe":
		mimetype = "image/jpeg"
	case ".png":
		mimetype = "image/png"
	case ".svg":
		mimetype = "image/svg+xml"
	case ".htm", ".html", ".xhtm", ".xhtml":
		mimetype = "application/xhtml+xml"
	case ".ncx":
		mimetype = "application/x-dtbncx+xml"
	case ".otf":
		mimetype = "application/vnd.ms-opentype"
	case ".woff":
		mimetype = "application/application/font-woff"
	case ".smil", ".smi", ".sml":
		mimetype = "application/smil+xml"
	case ".pls":
		mimetype = "application/pls+xml"
	case ".mp3":
		mimetype = "audio/mpeg"
	case ".mp4", ".aac", ".m4a", ".m4v", ".m4b", ".m4p", ".m4r":
		mimetype = "audio/mp4"
	case ".css":
		mimetype = "text/css"
	case ".js", ".javascript":
		mimetype = "text/javascript"
	default:
		if mimetype = mime.TypeByExtension(ext); mimetype == "" {
			mimetype = "application/octet-stream"
		}
	}
	self.counter++
	id := fmt.Sprintf("id%02x", self.counter)
	item := &Item{
		Id:        id,
		Href:      filename,
		MediaType: mimetype,
	}
	self.manifest = append(self.manifest, item)
	if spine {
		self.spine = append(self.spine, &ItemRef{IdRef: id})
	}
	return self.zipWriter.Create(path.Join(RootPath, filename))
}

func (self *Writer) Close() (err error) {
	defer self.file.Close()
	if self.metadata == nil {
		self.metadata = CreateMetadata(nil)
	}
	var uid string
	for _, item := range self.metadata.Identifier {
		if item.Id != "" {
			uid = item.Id
			break
		}
	}
	if uid == "" {
		self.metadata.Set("uid", "uid", uuid.New())
		uid = "uid"
	}
	var setTime bool
	for _, item := range self.metadata.Meta {
		if item.Property == "dcterms:modified" {
			item.Value = time.Now().UTC().Format(time.RFC3339)
			setTime = true
			break
		}
	}
	if !setTime {
		if self.metadata.Meta == nil {
			self.metadata.Meta = make([]*Meta, 0, 1)
		}
		self.metadata.Meta = append(self.metadata.Meta, &Meta{
			Property: "dcterms:modified",
			Value:    time.Now().UTC().Format(time.RFC3339),
		})
	}
	if len(self.metadata.Language) == 0 {
		self.metadata.Language.Set("", "en")
	}
	if len(self.metadata.Title) == 0 {
		self.metadata.Title.Set("", "Untitled")
	}
	item, err := self.zipWriter.Create(path.Join(RootPath, PackageFilename))
	if err != nil {
		return err
	}
	if _, err := io.WriteString(item, xml.Header); err != nil {
		return err
	}
	enc := xml.NewEncoder(item)
	enc.Indent("", "\t")
	opf := &Package{
		Version:          "3.0",
		UniqueIdentifier: uid,
		Metadata:         self.metadata,
		Manifest: &Manifest{
			Items: self.manifest,
		},
		Spine: &Spine{
			ItemRefs: self.spine,
		},
	}

	if err := enc.Encode(opf); err != nil {
		return err
	}
	if err := self.zipWriter.Close(); err != nil {
		return err
	}
	self.file.Commit()
	return nil
}
