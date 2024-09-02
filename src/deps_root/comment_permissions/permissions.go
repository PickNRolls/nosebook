package rootcommentpermissions

import (
	permissionscomment "nosebook/src/application/permissions/comment"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// TODO: загружать parent resource для доменной сущности сразу же

type Permissions struct {
	db *sqlx.DB
}

type CommentToRemove interface {
	Id() uuid.UUID
	AuthorId() uuid.UUID
}

func New(db *sqlx.DB) *Permissions {
	return &Permissions{
		db: db,
	}
}

type comment struct {
	authorId        uuid.UUID
	resourceOwnerId uuid.UUID
}

func (this *comment) AuthorId() uuid.UUID {
	return this.authorId
}

func (this *comment) ResourceOwnerId() uuid.UUID {
	return this.resourceOwnerId
}

func (this *Permissions) CanRemoveBy(commentToRemove CommentToRemove, userId uuid.UUID) *errors.Error {
	qb := querybuilder.New()
	sql, args, _ := qb.
		Select("owner_id as id").
		From("post_comments as pc").
		Join("posts as p on pc.post_id = p.id").
		Where("pc.comment_id = ?", commentToRemove.Id()).
		ToSql()

	var resourceOwner struct {
		Id uuid.UUID `db:"id"`
	}
	err := this.db.Get(&resourceOwner, sql, args...)
	if err != nil {
		return errors.New("Comment Permissions Error", "Вы не можете удалить комментарий")
	}

	c := &comment{
		authorId:        commentToRemove.AuthorId(),
		resourceOwnerId: resourceOwner.Id,
	}
	return permissionscomment.CanRemoveBy(c, userId)
}

func (this *Permissions) CanUpdateBy(comment permissionscomment.CommentToUpdate, userId uuid.UUID) *errors.Error {
	return permissionscomment.CanUpdateBy(comment, userId)
}
