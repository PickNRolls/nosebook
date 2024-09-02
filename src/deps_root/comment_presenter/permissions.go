package rootcommentpresenter

import (
	presentercomment "nosebook/src/application/presenters/comment"
	rootcommentpermissions "nosebook/src/deps_root/comment_permissions"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type permissions struct {
	original *rootcommentpermissions.Permissions
}

func newPermissions(db *sqlx.DB) *permissions {
	return &permissions{
		original: rootcommentpermissions.New(db),
	}
}

type comment struct {
	id       uuid.UUID
	authorId uuid.UUID
}

func (this *comment) Id() uuid.UUID {
	return this.id
}

func (this *comment) AuthorId() uuid.UUID {
	return this.authorId
}

func (this *permissions) CanRemoveBy(dest *presentercomment.Dest, userId uuid.UUID) bool {
	error := this.original.CanRemoveBy(&comment{
		id:       dest.Id,
		authorId: dest.AuthorId,
	}, userId)
	return error == nil
}

func (this *permissions) CanUpdateBy(dest *presentercomment.Dest, userId uuid.UUID) bool {
	err := this.original.CanUpdateBy(&comment{
		authorId: dest.AuthorId,
	}, userId)
	return err == nil
}
