package rootlikeservice

import (
	"nosebook/src/application/services/like"
	domainlike "nosebook/src/domain/like"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	postId    uuid.UUID
	commentId uuid.UUID
	userId    uuid.UUID
	db        *sqlx.DB
}

func newRepository(db *sqlx.DB) like.Repository {
	return &repository{
		db: db,
	}
}

func resourceTable(resource domainlike.Resource) string {
	resourceTable := ""

	switch resource.Type() {
	case domainlike.COMMENT_RESOURCE:
		resourceTable = "comment_likes"

	case domainlike.POST_RESOURCE:
		resourceTable = "post_likes"
	}

	return resourceTable
}

func resourceIdColumn(resource domainlike.Resource) string {
	resourceIdCol := ""

	switch resource.Type() {
	case domainlike.COMMENT_RESOURCE:
		resourceIdCol = "comment_id"

	case domainlike.POST_RESOURCE:
		resourceIdCol = "post_id"
	}

	return resourceIdCol
}

func ownerIdColumn(owner domainlike.Owner) string {
	ownerIdCol := ""

	switch owner.Type() {
	case domainlike.USER_OWNER:
		ownerIdCol = "user_id"
	}

	return ownerIdCol
}

func (this *repository) Save(like *domainlike.Like) *errors.Error {
	qb := querybuilder.New()

	if like.Event.Type() == domainlike.LIKED {
		sql, args, _ := qb.Insert(resourceTable(like.Resource)).
			Columns(ownerIdColumn(like.Owner), resourceIdColumn(like.Resource)).
			Values(like.Owner.Id(), like.Resource.Id()).
			ToSql()

		_, err := this.db.Exec(sql, args...)
		if err != nil {
			return errors.From(err)
		}

		return nil
	}

	//  TODO: если добавить еще один вариант owner,
	// то возможен баг: название таблицы тоже может быть зависимым от типа овнера, сейчас это не так.
	sql, args, _ := qb.Delete(resourceTable(like.Resource)).
		Where(resourceIdColumn(like.Resource)+" = ?", like.Resource.Id()).
		Where(ownerIdColumn(like.Owner)+" = ?", like.Owner.Id()).
		ToSql()

	_, err := this.db.Exec(sql, args...)
	if err != nil {
		return errors.From(err)
	}

	return nil
}

func (this *repository) WithPostId(id uuid.UUID) like.Repository {
	this.postId = id
	return this
}

func (this *repository) WithCommentId(id uuid.UUID) like.Repository {
	this.commentId = id
	return this
}

func (this *repository) WithUserId(id uuid.UUID) like.Repository {
	this.userId = id
	return this
}

func findByCommentAndUser(db *sqlx.DB, commentId uuid.UUID, userId uuid.UUID) (*domainlike.Like, *errors.Error) {
	qb := querybuilder.New()

	sql, args, _ := qb.Select("id").
		From("comments").
		Where("id = ?", commentId).
		Where("removed_at IS NULL").
		ToSql()

	commentDest := struct {
		Id uuid.UUID `db:"id"`
	}{}
	err := db.Get(&commentDest, sql, args...)
	if err != nil {
		return nil, like.NewCommentNotFoundError()
	}

	sql, args, _ = qb.Select("comment_id", "user_id").
		From("comment_likes").
		Where("comment_id = ? AND user_id = ?", commentId, userId).
		ToSql()

	likeDest := struct {
		CommentId uuid.UUID `db:"comment_id"`
		UserId    uuid.UUID `db:"user_id"`
	}{}
	err = db.Get(&likeDest, sql, args...)

	like, _ := domainlike.New().
		WithOwner(domainlike.NewUserOwner(userId)).
		WithCommentId(commentId)

	return like.WithValue(err == nil), nil
}

func findByPostAndUser(db *sqlx.DB, postId uuid.UUID, userId uuid.UUID) (*domainlike.Like, *errors.Error) {
	qb := querybuilder.New()

	sql, args, _ := qb.Select("id").
		From("posts").
		Where("id = ?", postId).
		Where("removed_at IS NULL").
		ToSql()

	postDest := struct {
		Id uuid.UUID `db:"id"`
	}{}
	err := db.Get(&postDest, sql, args...)
	if err != nil {
		return nil, like.NewPostNotFoundError()
	}

	sql, args, _ = qb.Select("post_id", "user_id").
		From("post_likes").
		Where("post_id = ? AND user_id = ?", postId, userId).
		ToSql()

	likeDest := struct {
		PostId uuid.UUID `db:"post_id"`
		UserId uuid.UUID `db:"user_id"`
	}{}
	err = db.Get(&likeDest, sql, args...)

	like, _ := domainlike.New().
		WithOwner(domainlike.NewUserOwner(userId)).
		WithPostId(postId)

	return like.WithValue(err == nil), nil
}

func (this *repository) FindOne() (*domainlike.Like, *errors.Error) {
	defer func() {
		this.postId = uuid.Nil
		this.commentId = uuid.Nil
		this.userId = uuid.Nil
	}()

	if this.commentId != uuid.Nil {
		return findByCommentAndUser(this.db, this.commentId, this.userId)
	}

	return findByPostAndUser(this.db, this.postId, this.userId)
}
