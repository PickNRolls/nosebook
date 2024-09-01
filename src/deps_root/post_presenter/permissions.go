package rootpostpresenter

import (
	"nosebook/src/application/permissions/permissionspost"
	presenterpost "nosebook/src/application/presenters/post"

	"github.com/google/uuid"
)

type permissions struct{}

type post struct {
	authorId uuid.UUID
}

func (this *post) AuthorId() uuid.UUID {
	return this.authorId
}

func (this *permissions) CanRemoveBy(p *presenterpost.Dest, userId uuid.UUID) bool {
	err := permissionspost.CanRemoveBy(&post{
		authorId: p.AuthorId,
	}, userId)
	return err == nil
}

func (this *permissions) CanUpdateBy(p *presenterpost.Dest, userId uuid.UUID) bool {
	err := permissionspost.CanUpdateBy(&post{
		authorId: p.AuthorId,
	}, userId)
	return err == nil
}
