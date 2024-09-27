package image

import (
	"bytes"
	"path/filepath"
)

type Image struct {
	data     []byte
	filename string
}

func New(filename string, bytes []byte) *Image {
	return &Image{
		data:     bytes,
		filename: filename,
	}
}

func (this *Image) Extension() string {
	return filepath.Ext(this.filename)
}

func (this *Image) NewReader() *bytes.Reader {
	return bytes.NewReader(this.data)
}
