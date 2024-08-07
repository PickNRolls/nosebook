package postgres

import (
	"nosebook/src/domain/comments"
	"nosebook/src/services/commenting/interfaces"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CommentRepository struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) interfaces.CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (repo *CommentRepository) FindById(id uuid.UUID) *comments.Comment {
	comment := comments.Comment{}
	err := repo.db.Get(&comment, `SELECT
		id,
		author_id,
		message,
		created_at
			FROM comments WHERE
		id = $1
	`, id)

	if err != nil {
		return nil
	}

	var likes []struct {
		UserId uuid.UUID `db:"user_id"`
	}
	err = repo.db.Select(&likes, `SELECT
		user_id
			FROM comment_likes WHERE
		comment_id = $1
	`, id)

	if err != nil {
		return nil
	}

	for _, like := range likes {
		comment.LikedBy = append(comment.LikedBy, like.UserId)
	}

	return &comment
}

func (repo *CommentRepository) Save(comment *comments.Comment) (*comments.Comment, error) {
	tx, err := repo.db.Beginx()
	if err != nil {
		return nil, err
	}

	for _, event := range comment.Events {
		switch event.Type() {
		case comments.LIKED:
			likeEvent := event.(*comments.CommentLikeEvent)
			_, err := tx.Exec(`INSERT INTO comment_likes (
				user_id,
				comment_id
			) VALUES (
				$1,
				$2
			)`, likeEvent.UserId, comment.Id)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

		case comments.UNLIKED:
			unlikeEvent := event.(*comments.CommentUnlikeEvent)
			_, err := tx.Exec(`DELETE FROM comment_likes WHERE
				user_id = $1 AND comment_id = $2
			`, unlikeEvent.UserId, comment.Id)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (repo *CommentRepository) Create(comment *comments.Comment) (*comments.Comment, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`INSERT INTO comments (
    id,
    author_id,
    message,
    created_at
  ) VALUES (
    $1,
    $2,
    $3,
    $4
  )`, comment.Id, comment.AuthorId, comment.Message, comment.CreatedAt)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if comment.PostId != uuid.Nil {
		_, err = tx.Exec(`INSERT INTO post_comments (
    	post_id,
    	comment_id
		) VALUES (
    	$1,
    	$2
		)`, comment.PostId, comment.Id)

		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return comment, nil
}

func (repo *CommentRepository) Remove(id uuid.UUID) (*comments.Comment, error) {
	_, err := repo.db.Exec(`UPDATE comments SET
		removed_at = NOW()
			WHERE
		id = $1
	`, id)

	if err != nil {
		return nil, err
	}

	var comment comments.Comment
	err = repo.db.Get(&comment, `SELECT
		id,
		author_id,
		message,
		created_at,
		removed_at
	FROM comments WHERE
		id = $1`, id)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}
