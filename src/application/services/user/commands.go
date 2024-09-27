package user

import (
	"nosebook/src/lib/image"
)

type ChangeAvatarCommand struct {
	Image *image.Image
}

func (this *ChangeAvatarCommand) Write(filename string, data []byte) error {
	this.Image = image.New(filename, data)
	return nil
}
