package mailx

import (
	"mime"
	"path/filepath"
)

// @author valor.

type file struct {
	// It is the name of file.
	// It is used for 'Content-ID', when the file is embedded
	filename string
	// It is the name of file when downloading, if the file is attachment
	// If empty, use filename.
	saved string

	// If true, the file is attachment.
	// If false, the file is embedded.
	attachment bool

	copier CopyFunc
}

func (f *file) contentType() string {
	mediaType := mime.TypeByExtension(filepath.Ext(f.filename))
	if mediaType == "" {
		mediaType = "application/octet-stream"
	}
	return mediaType + `; name="` + f.filename + `"`
}

func (f *file) disposition() string {
	filename := f.saved
	if filename == "" {
		filename = f.filename
	}

	disp := ""
	if f.attachment {
		disp = "attachment"
	} else {
		disp = "inline"
	}
	return disp + `; filename="` + filename + `"`
}
