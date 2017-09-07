package epub_test

import (
	"log"

	epub "github.com/mdigger/epub3"
)

func Example() {
	file, err := epub.Create("test.epub")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.Metadata.Title.Add("", "Test")
	err = file.AddFile("example.html", "example.html", epub.CTPrimary)
	if err != nil {
		log.Fatal(err)
	}
}
