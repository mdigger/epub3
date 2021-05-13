package epub

import (
	"archive/zip"
	"crypto/rand"
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

// newUUID returns the canonical namespaced string representation of a UUID:
//  urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func newUUID() string {
	var uuid [16]byte
	if _, err := io.ReadFull(rand.Reader, uuid[:]); err != nil {
		panic(err)
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // set version byte
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // set high order byte 0b10{8,9,a,b}
	return fmt.Sprintf("urn:uuid:%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
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
