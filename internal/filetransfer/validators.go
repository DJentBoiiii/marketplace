package filetransfer

import (
	"archive/zip"
	"bytes"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
)

func isAudio(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	audioFormats := []string{".mp3", ".wav", ".flac", ".ogg", ".m4a"}
	for _, format := range audioFormats {
		if ext == format {
			return true
		}
	}
	return false
}
func isArchive(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".zip"
}

func hasExecutableFiles(file multipart.File) bool {
	buf := new(bytes.Buffer)
	_, _ = io.Copy(buf, file)
	r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return false
	}
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".exe") || strings.HasSuffix(f.Name, ".bat") || strings.HasSuffix(f.Name, ".sh") {
			return true
		}
	}
	return false
}
