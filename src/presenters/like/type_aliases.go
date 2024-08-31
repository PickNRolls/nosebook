package presenterlike

import (
	presenterdto "nosebook/src/presenters/dto"

	"github.com/google/uuid"
)

type likesMap = map[uuid.UUID]*presenterdto.Likes
type usersMap = map[uuid.UUID]*presenterdto.User
