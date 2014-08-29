package epub

import (
	"archive/zip"
	"code.google.com/p/go-uuid/uuid"
	"encoding/xml"
	"fmt"
	"github.com/mdigger/commitfile"
	"io"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Content Type
type ContentType byte

const (
	ContentTypeMedia     ContentType = iota // Media file
	ContentTypeAuxiliary                    // Auxiliary content file
	ContentTypePrimary                      // Primary content file
)

// Writer allows you to create publications in epub 3 format.
type Writer struct {
	file      *commitfile.File
	zipWriter *zip.Writer
	Metadata  *Metadata
	manifest  []*Item
	spine     []*ItemRef
	counter   uint
}

// Create new epub publication.
func Create(filename string) (writer *Writer, err error) {
	// Создаем временный файл с публикацией
	file, err := commitfile.Create(filename)
	if err != nil {
		return nil, err
	}
	// В случае ошибки закрываем и удаляем его при выходе из функции
	defer func() {
		if err != nil {
			file.Close()
		}
	}()
	// Инициализируем упаковку в архив
	zipWriter := zip.NewWriter(file)
	var item io.Writer
	// Добавляем информацию о mimetype
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
	// Добавляем описание контейнера
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
	// Инициализируем объект с описанием публикации
	writer = &Writer{
		file:      file,
		zipWriter: zipWriter,
		manifest:  make([]*Item, 0, 10),
		spine:     make([]*ItemRef, 0, 10),
	}
	return writer, nil
}

// AddFile adds the source file to the publication with name filename.
func (w *Writer) AddFile(sourceFilename, filename string, ct ContentType, properties ...string) error {
	file, err := os.Open(sourceFilename)
	if err != nil {
		return err
	}
	defer file.Close()
	fileWriter, err := w.Add(filename, ct, properties...)
	if err != nil {
		return err
	}
	if _, err := io.Copy(fileWriter, file); err != nil {
		return err
	}
	return nil
}

// Add returns the io.writer to write the data to the publication.
func (w *Writer) Add(filename string, ct ContentType, properties ...string) (io.Writer, error) {
	filename = filepath.ToSlash(filename) // Нормализуем имя файла
	// Проверяем, что файла с таким именем еще нет в публикации.
	// Иначе возвращаем ошибку.
	for _, item := range w.manifest {
		if item.Href == filename {
			return nil,
				fmt.Errorf("a file with the name %q has already been added to the publication", filename)
		}
	}
	// Вычисляем mimetype по расширению файла
	var mimetype string
	switch ext := strings.ToLower(path.Ext(filename)); ext {
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
	// case ".json":
	// 	mimetype = "application/json"
	default:
		if mimetype = mime.TypeByExtension(ext); mimetype == "" {
			mimetype = "application/octet-stream"
		}
	}
	w.counter++ // Увеличиваем счетчик добавленных файлов
	id := fmt.Sprintf("id%02x", w.counter)
	// Создаем описание добавляемого файла
	item := &Item{
		ID:         id,
		Href:       filename,
		MediaType:  mimetype,
		Properties: strings.Join(properties, " "),
	}
	// Добавляем описание в список
	w.manifest = append(w.manifest, item)
	// Если необходимо, то добавляем идентификатор файла в список чтения
	if ct > ContentTypeMedia {
		itemref := &ItemRef{IDRef: id}
		if ct == ContentTypeAuxiliary {
			itemref.Linear = "no"
		}
		w.spine = append(w.spine, itemref)
	}
	// Возвращаем writer для записи содержимого файла
	return w.zipWriter.Create(path.Join(RootPath, filename))
}

// Close closes the publication and writes metadata.
func (w *Writer) Close() (err error) {
	// Закрываем файл по окончании
	defer w.file.Close()
	// Инициализируем метаданные, если они не были инициализированы раньше
	metadata := w.Metadata
	if metadata == nil {
		metadata = new(Metadata)
	}
	// Проверяем и добавляем по необходимости обязательные элементы
	if metadata.DC == "" {
		metadata.DC = "http://purl.org/dc/elements/1.1/"
	}
	// Получаем идентификатор уникального идентификатора публикации
	var uid string
	for _, item := range metadata.Identifier {
		if item.ID != "" {
			uid = item.ID
			break
		}
	}
	if uid == "" {
		metadata.Add("uuid", "uuid", "urn:uuid:"+uuid.New())
		uid = "uuid"
	}
	// Добавляем дату модификации
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
	// Устанавливаем язык, если его нет
	if len(metadata.Language) == 0 {
		metadata.Language.Add("", "en")
	}
	// Добавляем заголовок, если его нет
	if len(metadata.Title) == 0 {
		metadata.Title.Add("", "Untitled")
	}
	// Сериализуем описание публикации
	item, err := w.zipWriter.Create(path.Join(RootPath, PackageFilename))
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
		Metadata:         metadata,
		Manifest: &Manifest{
			Items: w.manifest,
		},
		Spine: &Spine{
			ItemRefs: w.spine,
		},
	}

	if err := enc.Encode(opf); err != nil {
		return err
	}
	// Закрываем упаковку
	if err := w.zipWriter.Close(); err != nil {
		return err
	}
	// Отменяем автоудаление файла
	w.file.Commit()
	return nil
}
