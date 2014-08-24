package main

import (
	"encoding/xml"
	"flag"
	"github.com/mdigger/epub3"
	"github.com/mdigger/metadata"
	"github.com/russross/blackfriday"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	log.SetFlags(0)
	var (
		sourcePath     string // Путь к файлам проекта
		outputFilename string // Имя результирующего файла с публикацией
	)
	flag.StringVar(&sourcePath, "source", "", "path to source files")
	flag.StringVar(&outputFilename, "out", "output.epub", "publication output filename")
	flag.Parse()
	// Делаем исходный каталог текущим, чтобы не вычислять относительный путь. По окончании
	// обработки восстанавливаем исходный каталог обратно.
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if err := os.Chdir(sourcePath); err != nil {
		log.Fatal(err)
	}
	defer os.Chdir(currentPath)
	// Создаем описание публикации
	meta := &epub.Metadata{
		DC:       "http://purl.org/dc/elements/1.1/",
		Title:    epub.SimpleProperty("Тестовая публикация"),
		Language: epub.SimpleProperty("ru"),
		Identifier: []*epub.MetaProperty{
			&epub.MetaProperty{
				Id:    "uid",
				Value: "test",
			},
		},
		Creator: epub.SimpleProperty("Дмитрий Седых"),
	}
	// Инициализируем шаблон для преобразования страниц
	tpage, err := template.New("").Parse(pageTemplateText)
	if err != nil {
		log.Fatal(err)
	}
	// Создаем упаковщик в формат EPUB
	writer, err := epub.Create(filepath.Join(currentPath, outputFilename), meta)
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()
	// Определяем функция для обработки перебора файлов и каталогов
	walkFn := func(filename string, finfo os.FileInfo, err error) error {
		// Игнорируем, если открытие файла произошло с ошибкой
		if err != nil {
			return nil
		}
		// Не обрабатываем отдельно каталоги
		if finfo.IsDir() {
			return nil
		}
		switch filepath.Ext(filename) {
		case ".md", ".mdown", ".markdown":
			log.Println("Markdown:", filename)
			// Читаем файл и отделяем метаданные
			meta, data, err := metadata.ReadFile(filename)
			if err != nil {
				log.Fatal(err)
			}
			// Преобразуем из Markdown в HTML
			data = MarkdownCommon(data)
			// Сохраняем результат прямо в метаданных под именем content.
			// Предварительно "оборачиваем" в шаблонное представление HTML,
			// чтобы он не декодировался.
			meta["content"] = template.HTML(data)
			// Если не указан язык, то считаем, что он русский.
			if _, ok := meta["lang"]; !ok {
				meta["lang"] = "ru"
			}
			// Изменяем расширение имени файла на .html
			filename = filename[:len(filename)-len(filepath.Ext(filename))] + ".xhtml"
			// TODO: добавлять в spine или нет, в зависимости от имени.
			fileWriter, err := writer.Add(filename, true)
			if err != nil {
				log.Fatal(err)
			}
			// Добавляем в начало документа XML-заголовок
			io.WriteString(fileWriter, xml.Header)
			// Преобразуем по шаблону и записываем в публикацию.
			if err := tpage.Execute(fileWriter, meta); err != nil {
				log.Fatal(err)
			}
		case ".jpg", ".jpe", ".jpeg", ".png", ".gif", ".svg":
			log.Println("Add image:", filename)
			file, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			fileWriter, err := writer.Add(filename, false)
			if err != nil {
				log.Fatal(err)
			}
			if _, err := io.Copy(fileWriter, file); err != nil {
				log.Fatal(err)
			}
		default:
			log.Println("Ignore:", filename)
		}
		return nil
	}
	// Перебираем все файлы и подкаталоги в исходном каталоге
	if err = filepath.Walk(".", walkFn); err != nil {
		log.Fatal(err)
	}
}

func MarkdownCommon(input []byte) []byte {
	// set up the HTML renderer
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	// htmlFlags |= blackfriday.HTML_SANITIZE_OUTPUT // Error in img tag
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

	// set up the parser
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS
	extensions |= blackfriday.EXTENSION_HEADER_IDS

	return blackfriday.Markdown(input, renderer, extensions)
}

const pageTemplateText = `<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" xml:lang="{{ if .lang }}{{ .lang }}{{ else }}en{{ end }}">
<head>
<meta charset="utf-8" />
<title>{{ .title }}</title>
</head>
<body>
{{ if .title }}<h1>{{ .title }}</h1>

{{ end }}{{ .content }}
</body>
</html>`
