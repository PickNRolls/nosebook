package rootcommentpresenter

import (
	permissionscomment "nosebook/src/application/permissions/comment"
	presentercomment "nosebook/src/application/presenters/comment"

	"github.com/google/uuid"
)

type permissions struct{}

type comment struct {
	authorId uuid.UUID
}

func (this *comment) AuthorId() uuid.UUID {
	return this.authorId
}

func (this *permissions) CanRemoveBy(dest *presentercomment.Dest, userId uuid.UUID) bool {
	comment := &comment{
		authorId: dest.AuthorId,
	}
	err := permissionscomment.CanRemoveBy(comment, userId)
	return err == nil
}

func (this *permissions) CanUpdateBy(dest *presentercomment.Dest, userId uuid.UUID) bool {
	comment := &comment{
		authorId: dest.AuthorId,
	}
	err := permissionscomment.CanUpdateBy(comment, userId)
	return err == nil
}
