package rootcommentpresenter

import (
	permissionscomment "nosebook/src/application/permissions/comment"
	presentercomment "nosebook/src/application/presenters/comment"
	"nosebook/src/infra/postgres"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type permissions struct {
	db *sqlx.DB
}

func newPermissions(db *sqlx.DB) *permissions {
	return &permissions{
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

func (this *permissions) CanRemoveBy(dest *presentercomment.Dest, userId uuid.UUID) bool {
	qb := postgres.NewSquirrel()
	sql, args, _ := qb.
		Select("owner_id as id").
		From("post_comments as pc").
		Join("posts as p on pc.post_id = p.id").
		Where("pc.comment_id = ?", dest.Id).
		ToSql()

	var resourceOwner struct {
		Id uuid.UUID `db:"id"`
	}
	err := this.db.Get(&resourceOwner, sql, args...)
	if err != nil {
		return false
	}

	comment := &comment{
		authorId:        dest.AuthorId,
		resourceOwnerId: resourceOwner.Id,
	}
	error := permissionscomment.CanRemoveBy(comment, userId)
	return error == nil
}

func (this *permissions) CanUpdateBy(dest *presentercomment.Dest, userId uuid.UUID) bool {
	comment := &comment{
		authorId: dest.AuthorId,
	}
	err := permissionscomment.CanUpdateBy(comment, userId)
	return err == nil
}
