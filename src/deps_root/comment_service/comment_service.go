package rootcommentservice

import (
	"nosebook/src/application/services/commenting"
	rootcommentpermissions "nosebook/src/deps_root/comment_permissions"
	rootpostservice "nosebook/src/deps_root/post_service"
	domaincomment "nosebook/src/domain/comment"
	"nosebook/src/errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

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

type permissions struct {
	original rootcommentpermissions.Permissions
}

func (this *permissions) CanRemoveBy(c *domaincomment.Comment, userId uuid.UUID) *errors.Error {
	return this.original.CanRemoveBy(&comment{
		id:       c.Id,
		authorId: c.AuthorId,
	}, userId)
}

func (this *permissions) CanUpdateBy(c *domaincomment.Comment, userId uuid.UUID) *errors.Error {
	return this.original.CanUpdateBy(&comment{
		id:       c.Id,
		authorId: c.AuthorId,
	}, userId)
}

func New(db *sqlx.DB) *commenting.Service {
	commentService := commenting.New(newCommentRepository(db, &permissions{
		original: *rootcommentpermissions.New(db),
	}), rootpostservice.NewRepository(db))

	return commentService
}
