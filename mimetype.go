package epub

import (
	"mime"
	"path/filepath"
	"strings"
)

// MimeTypes list supported MIME types
var MimeTypes = map[string]string{
	".gif":   "image/gif",
	".jpg":   "image/jpeg",
	".jpeg":  "image/jpeg",
	".jpe":   "image/jpeg",
	".png":   "image/png",
	".svg":   "image/svg+xml",
	".xhtm":  "application/xhtml+xml",
	".xhtml": "application/xhtml+xml",
	".ncx":   "application/x-dtbncx+xml",
	".otf":   "application/vnd.ms-opentype",
	".woff":  "application/application/font-woff",
	".smil":  "application/smil+xml",
	".smi":   "application/smil+xml",
	".sml":   "application/smil+xml",
	".pls":   "application/pls+xml",
	".mp3":   "audio/mpeg",
	".mp4":   "audio/mp4",
	".aac":   "audio/mp4",
	".m4a":   "audio/mp4",
	".m4v":   "audio/mp4",
	".m4b":   "audio/mp4",
	".m4p":   "audio/mp4",
	".m4r":   "audio/mp4",
	".css":   "text/css",
	".js":    "text/javascript",
}

// TypeByFilename returns the MIME type associated with the file name.
func TypeByFilename(filename string) (mimetype string) {
	ext := strings.ToLower(filepath.Ext(filename))
	if mimetype = MimeTypes[ext]; mimetype != "" {
		return mimetype
	}
	if mimetype = mime.TypeByExtension(ext); mimetype != "" {
		return mimetype
	}
	return "application/octet-stream"
}
