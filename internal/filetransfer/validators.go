package filetransfer

import (
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
