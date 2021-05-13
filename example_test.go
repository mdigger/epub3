package epub_test

import (
	"log"
	"os"

	epub "github.com/mdigger/epub3"
)

func Example() {
	file, err := os.Create("test.epub")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	pub, err := epub.New(file)
	if err != nil {
		log.Fatal(err)
	}
	defer pub.Close()

	pub.Title.Add("", "Test")
	err = pub.AddContentFile("example.html", "example.html", epub.Primary)
	if err != nil {
		log.Fatal(err)
	}
}
