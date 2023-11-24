package mailx

import (
	"mime"
	"path/filepath"
)

// @author valor.

type file struct {
	// It is the name of file.
	// It is used for 'Content-ID', if the file is embedded.
	filename string
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
	return mediaType
}

func (f *file) disposition() string {
	disp := ""
	if f.attachment {
		disp = "attachment"
	} else {
		disp = "inline"
	}
	return disp + `; filename="` + f.filename + `"`
}
