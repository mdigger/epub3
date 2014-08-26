package main

import (
	"encoding/xml"
	"github.com/kr/pretty"
	"github.com/mdigger/epub3"
	"github.com/mdigger/metadata"
	"github.com/russross/blackfriday"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	metadataFiles = []string{"metadata.yaml", "metadata.yml", "metadata.json"}
	coverFiles    = []string{"cover.png", "cover.svg", "cover.jpeg", "cover.jpg", "cover.gif"}
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
	// Загружаем описание метаданных публикации
	var pubMetadata *epub.Metadata
	for _, name := range metadataFiles {
		fi, err := os.Stat(name)
		if err != nil || fi.IsDir() {
			continue
		}
		if pubMetadata, err = loadMetadata(name); err != nil {
			return err
		}
		break
	}
	// Если описания не найдено, то инициализируем пустое описание
	if pubMetadata == nil {
		pubMetadata = defaultMetada()
	}
	// Добавляем язык, если он не определен
	if len(pubMetadata.Language) == 0 {
		pubMetadata.Language.Add("", defaultLang)
	}
	// Вытаскиваем язык публикации
	publang := pubMetadata.Language[0].Value
	// Создаем упаковщик в формат EPUB
	writer, err := epub.Create(filepath.Join(currentPath, outputFilename))
	if err != nil {
		return err
	}
	defer writer.Close()
	// Добавляем метаданные в публикацию
	writer.Metadata = pubMetadata
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
		switch strings.ToLower(filepath.Ext(filename)) {
		case ".md", ".mdown", ".markdown": // Статья в формате Markdown: преобразуем и добавляем
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
			fileWriter, err := writer.Writer(filename, spine)
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
		case ".jpg", ".jpe", ".jpeg", ".png", ".gif", ".svg",
			".mp3", ".mp4", ".aac", ".m4a", ".m4v", ".m4b", ".m4p", ".m4r",
			".css", ".js", ".javascript",
			".json",
			".otf", ".woff",
			".pls", ".smil", ".smi", ".sml": // Иллюстрации и другие файлы
			var properties []string
			// Специальная обработка обложки
			if !setCover && isFilename(filename, coverFiles) {
				properties = []string{"cover-image"}
				setCover = true
			}
			log.Printf("Add: %s %s", filename, strings.Join(properties, ", "))
			if err := writer.AddFile(filename, filename, false, properties...); err != nil {
				return err
			}
		default: // Другое — игнорируем
			log.Println("Ignore:", filename)
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
