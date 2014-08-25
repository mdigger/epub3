package epub

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"github.com/mdigger/commitfile"
	"io"
	"mime"
	"path"
	"path/filepath"
)

var (
	RootPath        = "OEBPS"
	EPUBMimeType    = "application/epub+zip"
	PackageFilename = "package.opf"
)

type Writer struct {
	file      *commitfile.File
	zipWriter *zip.Writer
	opf       *Package
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
	opf := &Package{
		Version:          "3.0",
		UniqueIdentifier: "uid",
		Metadata:         metadata,
		Manifest: &Manifest{
			Items: make([]*Item, 0, 10),
		},
		Spine: &Spine{
			ItemRefs: make([]*ItemRef, 0, 10),
		},
	}
	writer = &Writer{
		file:      file,
		zipWriter: zipWriter,
		opf:       opf,
	}
	return writer, nil
}

func (self *Writer) Add(filename string, spine bool) (io.Writer, error) {
	filename = filepath.ToSlash(filename)
	var mimetype string
	switch ext := path.Ext(filename); ext {
	case ".htm", ".html", ".xhtm", ".xhtml":
		mimetype = "application/xhtml+xml"
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
	self.opf.Manifest.Items = append(self.opf.Manifest.Items, item)
	if spine {
		self.opf.Spine.ItemRefs = append(self.opf.Spine.ItemRefs, &ItemRef{IdRef: id})
	}
	return self.zipWriter.Create(path.Join(RootPath, filename))
}

func (self *Writer) Close() (err error) {
	defer self.file.Close()
	item, err := self.zipWriter.Create(path.Join(RootPath, PackageFilename))
	if err != nil {
		return err
	}
	if _, err := io.WriteString(item, xml.Header); err != nil {
		return err
	}
	enc := xml.NewEncoder(item)
	enc.Indent("", "\t")
	if err := enc.Encode(self.opf); err != nil {
		return err
	}
	if err := self.zipWriter.Close(); err != nil {
		return err
	}
	self.file.Commit()
	return nil
}
