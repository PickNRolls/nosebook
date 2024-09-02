package presenterpost

import (
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/application/services/auth"
	"nosebook/src/errors"
	"time"

	"github.com/google/uuid"
)

type FindByFilterInput struct {
	Ids      []string
	OwnerId  string
	AuthorId string
	Cursor   string
}

type FindByFilterOutput presenterdto.FindOut[*Post]

type user = presenterdto.User
type comments = presenterdto.FindOut[*presenterdto.Comment]
type likes = presenterdto.Likes

type Post struct {
	Id             uuid.UUID                 `json:"id"`
	Author         *user                     `json:"author"`
	Owner          *user                     `json:"owner"`
	Message        string                    `json:"message"`
	Likes          *likes                    `json:"likes"`
	RecentComments *comments                 `json:"recentComments"`
	Permissions    *presenterdto.Permissions `json:"permissions"`
	CreatedAt      time.Time                 `json:"createdAt"`
}

type UserPresenter interface {
	FindByIds(ids uuid.UUIDs) ([]*user, *errors.Error)
}

type CommentPresenter interface {
	FindByPostId(id uuid.UUID, auth *auth.Auth) *comments
}

type LikePresenter interface {
	FindByPostIds(ids uuid.UUIDs, auth *auth.Auth) (map[uuid.UUID]*likes, *errors.Error)
}

type Permissions interface {
	CanRemoveBy(post *Dest, userId uuid.UUID) bool
	CanUpdateBy(post *Dest, userId uuid.UUID) bool
}
