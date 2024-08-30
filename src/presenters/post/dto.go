package presenterpost

import (
	presenterdto "nosebook/src/presenters/dto"
	"time"

	"github.com/google/uuid"
)

type FindByFilterInput struct {
	OwnerId  string
	AuthorId string
	Cursor   string
}

type FindByFilterOutput presenterdto.FindOut[*post]

type user = presenterdto.User
type comments = presenterdto.FindOut[*presenterdto.Comment]
type likes = presenterdto.Likes

type post struct {
	Id             uuid.UUID `json:"id"`
	Author         *user     `json:"author"`
	Owner          *user     `json:"owner"`
	Message        string    `json:"message"`
	Likes          *likes    `json:"likes"`
	RecentComments *comments `json:"recentComments"`
	CreatedAt      time.Time `json:"createdAt"`
}
