package main

import (
	"encoding/xml"
	"fmt"
	"github.com/kr/pretty"
	"github.com/mdigger/epub3"
	"github.com/mdigger/metadata"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v1"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"
)

// Компилятор публикации
func compiler(sourcePath, outputFilename string) error {
	// Делаем исходный каталог текущим, чтобы не вычислять относительный путь. По окончании
	// обработки восстанавливаем исходный каталог обратно.
	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(sourcePath); err != nil {
		return err
	}
	defer os.Chdir(currentPath)
	// Инициализируем шаблон для преобразования страниц
	tpage, err := template.New("").Parse(pageTemplateText)
	if err != nil {
		return err
	}
	// Инициализируем описание метаданных
	pubmeta := &epub.Metadata{
		DC:   "http://purl.org/dc/elements/1.1/",
		Meta: make([]*epub.Meta, 0),
	}
	publang := "en" // Язык публикации по умолчанию
	// Загружаем описание метаданных публикации
	for _, name := range []string{"metadata.yml", "metadata.yaml", "metadata.json"} {
		fi, err := os.Stat(name)
		if err != nil || fi.IsDir() {
			continue
		}
		data, err := ioutil.ReadFile(name)
		if err != nil {
			return err
		}
		meta := make(metadata.Metadata)
		if err := yaml.Unmarshal(data, meta); err != nil {
			return err
		}
		log.Print(strings.Repeat("—", 80))
		tab := tabwriter.NewWriter(os.Stdout, 8, 1, 1, ' ', 0)
		// Конвертируем описание метаданных в метаданные
		// Добавляем язык
		lang := meta.Lang()
		if lang == "" {
			lang = publang
		} else {
			publang = lang
		}
		pubmeta.Language.Add("", lang)
		fmt.Fprintf(tab, "Lang:\t%s\n", lang)
		// Добавляем заголовок
		if title := meta.Title(); title != "" {
			pubmeta.Title.Add("title", title)
			pubmeta.Meta = append(pubmeta.Meta, &epub.Meta{
				Refines:  "#title",
				Property: "title-type",
				Value:    "main",
			})
			fmt.Fprintf(tab, "Title:\t%s\n", title)
		}
		// Добавляем подзаголовок
		if subtitle := meta.Subtitle(); subtitle != "" {
			pubmeta.Title.Add("subtitle", subtitle)
			pubmeta.Meta = append(pubmeta.Meta, &epub.Meta{
				Refines:  "#subtitle",
				Property: "title-type",
				Value:    "subtitle",
			})
			fmt.Fprintf(tab, "Subtitle:\t%s\n", subtitle)
		}
		// Добавляем название коллекции
		if collection := meta.Get("collection"); collection != "" {
			pubmeta.Title.Add("collection", collection)
			pubmeta.Meta = append(pubmeta.Meta, &epub.Meta{
				Refines:  "#collection",
				Property: "title-type",
				Value:    "collection",
			})
			fmt.Fprintf(tab, "Collection:\t%s\n", collection)
		}
		// Добавляем название редакции
		if edition := meta.Get("edition"); edition != "" {
			pubmeta.Title.Add("edition", edition)
			pubmeta.Meta = append(pubmeta.Meta, &epub.Meta{
				Refines:  "#edition",
				Property: "title-type",
				Value:    "edition",
			})
			fmt.Fprintf(tab, "Edition:\t%s\n", edition)
		}
		// Добавляем авторов
		for i, author := range meta.Authors() {
			pubmeta.Creator.Add("", author)
			if i == 0 {
				fmt.Fprintf(tab, "Author:\t%s\n", author)
			} else {
				fmt.Fprintf(tab, "\t%s\n", author)
			}
		}
		// Добавляем второстепенных авторов
		for i, author := range meta.GetList("contributor") {
			pubmeta.Contributor.Add("", author)
			if i == 0 {
				fmt.Fprintf(tab, "Contributor:\t%s\n", author)
			} else {
				fmt.Fprintf(tab, "\t%s\n", author)
			}
		}
		// Добавляем информацию об издателях
		for i, author := range meta.GetList("publisher") {
			pubmeta.Publisher.Add("", author)
			if i == 0 {
				fmt.Fprintf(tab, "Publisher:\t%s\n", author)
			} else {
				fmt.Fprintf(tab, "\t%s\n", author)
			}
		}
		// Добавляем уникальные идентификаторы
		for _, name := range []string{"uuid", "id", "identifier", "doi", "isbn", "issn"} {
			if value := meta.Get(name); value != "" {
				pubmeta.Identifier.Add(name, value)
				fmt.Fprintf(tab, "%s:\t%s\n", strings.ToUpper(name), value)
			}
		}
		// Добавляем краткое описание
		if description := meta.Description(); description != "" {
			pubmeta.Description.Add("description", description)
		}
		// Добавляем ключевые слова
		if keywords := meta.Keywords(); len(keywords) > 0 {
			for _, keyword := range keywords {
				pubmeta.Subject.Add("", keyword)
			}
			fmt.Fprintf(tab, "Keywords:\t%s\n", strings.Join(keywords, ", "))
		}
		// Добавляем описание сферы действия
		if coverage := meta.Get("coverage"); coverage != "" {
			pubmeta.Coverage.Add("", coverage)
		}
		// Добавляем дату
		if date := meta.Date(); !date.IsZero() {
			pubmeta.Date = &epub.Element{
				Value: date.UTC().Format(time.RFC3339),
			}
		}
		// Добавляем копирайты
		for _, name := range []string{"copyright", "rights"} {
			if rights := meta.Get(name); rights != "" {
				pubmeta.Rights.Add(name, rights)
				fmt.Fprintf(tab, "%s:\t%s\n", strings.Title(name), rights)
			}
		}
		tab.Flush()
		log.Print(strings.Repeat("—", 80))
		break
	}
	// Создаем упаковщик в формат EPUB
	writer, err := epub.Create(filepath.Join(currentPath, outputFilename))
	if err != nil {
		return err
	}
	defer writer.Close()
	// Добавляем метаданные в публикацию
	writer.Metadata = pubmeta
	// Функция для добавления файла в публикацию
	addFile := func(filename string, spine bool, properties ...string) error {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		fileWriter, err := writer.Add(filename, spine, properties...)
		if err != nil {
			return err
		}
		if _, err := io.Copy(fileWriter, file); err != nil {
			return err
		}
		return nil
	}
	// Инициализируем преобразование из формата Markdown
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	markdownRender := blackfriday.HtmlRenderer(htmlFlags, "", "")
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS
	extensions |= blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK
	extensions |= blackfriday.EXTENSION_HEADER_IDS
	// Оглавление
	nav := make(Navigaton, 0)
	// Флаг для избежания двойной обработки обложки: после его установки
	// новые попадающиеся обложки игнорируются.
	var setCover bool
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
		// Проверяем по имени файла
		switch strings.ToLower(filename) {
		// Описание метаданных публикации — уже загружено, если было
		case "metadata.yml", "metadata.yaml", "metadata.json":
			return nil
		// Обложка публикации
		case "cover.gif", "cover.jpg", "cover.jpeg", "cover.png", "cover.svg":
			if setCover {
				log.Println("Ignore duplicate cover image:", filename)
				return nil
			}
			log.Println("Add cover image:", filename)
			if err := addFile(filename, false, "cover-image"); err != nil {
				return err
			}
			setCover = true
		// Другие файлы
		default:
			// В зависимости от расширения имени файла
			switch strings.ToLower(filepath.Ext(filename)) {
			// Статья в формате Markdown: преобразуем и добавляем
			case ".md", ".mdown", ".markdown":
				log.Println("Markdown:", filename)
				// Читаем файл и отделяем метаданные
				meta, data, err := metadata.ReadFile(filename)
				if err != nil {
					return err
				}
				// Преобразуем из Markdown в HTML
				data = blackfriday.Markdown(data, markdownRender, extensions)
				// Сохраняем результат прямо в метаданных под именем content.
				// Предварительно "оборачиваем" в шаблонное представление HTML,
				// чтобы он не декодировался.
				meta["content"] = template.HTML(data)
				// Если не указан язык, то считаем, что он русский.
				if _, ok := meta["lang"]; !ok {
					meta["lang"] = publang
				}
				// Изменяем расширение имени файла на .xhtml
				filename = filename[:len(filename)-len(filepath.Ext(filename))] + ".xhtml"
				// Добавляем в основной список чтения, если имя файла не начинается с подчеркивания
				spine := filepath.Base(filename)[0] != '_'
				fileWriter, err := writer.Add(filename, spine)
				if err != nil {
					return err
				}
				// Добавляем в начало документа XML-заголовок
				if _, err := io.WriteString(fileWriter, xml.Header); err != nil {
					return err
				}
				// Преобразуем по шаблону и записываем в публикацию.
				if err := tpage.Execute(fileWriter, meta); err != nil {
					return err
				}
				// Добавляем информацию о файле в оглавление
				nav = append(nav, &NavigationItem{
					Title:    meta.Title(),
					Subtitle: meta.Subtitle(),
					Filename: filename,
					Spine:    spine,
				})
			// Иллюстрации и другие файлы
			case ".jpg", ".jpe", ".jpeg", ".png", ".gif", ".svg",
				".mp3", ".mp4", ".aac", ".m4a", ".m4v", ".m4b", ".m4p", ".m4r",
				".css", ".js", ".javascript",
				".json",
				".otf", ".woff",
				".pls", ".smil", ".smi", ".sml":
				log.Println("Add file:", filename)
				if err := addFile(filename, false); err != nil {
					return err
				}
			// Другое — игнорируем
			default:
				log.Println("Ignore:", filename)
			}
		}
		return nil
	}
	// Перебираем все файлы и подкаталоги в исходном каталоге
	if err := filepath.Walk(".", walkFn); err != nil {
		return err
	}
	pretty.Println(nav)
	return nil
}

const pageTemplateText = `<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" xml:lang="{{ if .lang }}{{ .lang }}{{ else }}en{{ end }}">
<head>
<meta charset="UTF-8" />
<title>{{ .title }}</title>
</head>
<body>
{{ if .title }}<h1>{{ .title }}</h1>

{{ end }}{{ .content }}
</body>
</html>`
