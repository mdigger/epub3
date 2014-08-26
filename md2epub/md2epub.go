package main

import (
	"encoding/xml"
	"flag"
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
)

func main() {
	log.SetFlags(0)
	// Разбираем входящие параметры
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(2)
	}
	sourcePath := flag.Arg(0)
	var outputFilename string // Имя результирующего файла с публикацией
	if flag.NArg() > 1 {
		outputFilename = flag.Arg(1)
	} else {
		outputFilename = filepath.Base(sourcePath) + ".epub"
	}
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
	// Инициализируем шаблон для преобразования страниц
	tpage, err := template.New("").Parse(pageTemplateText)
	if err != nil {
		log.Fatal(err)
	}
	// Создаем упаковщик в формат EPUB
	writer, err := epub.Create(filepath.Join(currentPath, outputFilename))
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()
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
	// Флаги для избежания двойной обработки метаданных и обложки
	var setCover, setMetada bool
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
		switch filename {
		// Описание метаданных публикации
		case "metadata.yml", "metadata.yaml", "metadata.json":
			if setMetada {
				log.Println("Ignore duplicate metadata:", filename)
				return nil
			}
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Fatal(err)
			}
			meta := make(metadata.Metadata)
			if err := yaml.Unmarshal(data, meta); err != nil {
				log.Fatal(err)
			}
			// TODO: заполнить метаданные
			setMetada = true
		// Обложка публикации
		case "cover.gif", "cover.jpg", "cover.jpeg", "cover.png", "cover.svg":
			if setCover {
				log.Println("Ignore duplicate cover image:", filename)
				return nil
			}
			log.Println("Add cover image:", filename)
			file, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			// Добавляем свойство, что это именно обложка
			fileWriter, err := writer.Add(filename, false, "cover-image")
			if err != nil {
				log.Fatal(err)
			}
			if _, err := io.Copy(fileWriter, file); err != nil {
				log.Fatal(err)
			}
			setCover = true
		// Другие файлы
		default:
			// В зависимости от расширения имени файла
			switch filepath.Ext(filename) {
			// Статья в формате Markdown: преобразуем и добавляем
			case ".md", ".mdown", ".markdown":
				log.Println("Markdown:", filename)
				// Читаем файл и отделяем метаданные
				meta, data, err := metadata.ReadFile(filename)
				if err != nil {
					log.Fatal(err)
				}
				// Преобразуем из Markdown в HTML
				data = blackfriday.Markdown(data, markdownRender, extensions)
				// Сохраняем результат прямо в метаданных под именем content.
				// Предварительно "оборачиваем" в шаблонное представление HTML,
				// чтобы он не декодировался.
				meta["content"] = template.HTML(data)
				// Если не указан язык, то считаем, что он русский.
				if _, ok := meta["lang"]; !ok {
					meta["lang"] = "ru"
				}
				// Изменяем расширение имени файла на .xhtml
				filename = filename[:len(filename)-len(filepath.Ext(filename))] + ".xhtml"
				// Добавляем в основной список чтения, если имя файла не начинается с подчеркивания
				fileWriter, err := writer.Add(filename, filepath.Base(filename)[0] != '_')
				if err != nil {
					log.Fatal(err)
				}
				// Добавляем в начало документа XML-заголовок
				io.WriteString(fileWriter, xml.Header)
				// Преобразуем по шаблону и записываем в публикацию.
				if err := tpage.Execute(fileWriter, meta); err != nil {
					log.Fatal(err)
				}
			// Иллюстрация — добавляем в публикацию как есть
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
			// Другое — игнорируем
			default:
				log.Println("Ignore:", filename)
			}
		}
		return nil
	}
	// Перебираем все файлы и подкаталоги в исходном каталоге
	if err := filepath.Walk(".", walkFn); err != nil {
		log.Fatal(err)
	}
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
