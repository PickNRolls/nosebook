package userauth

import "github.com/google/uuid"

type CreateSessionCommand struct {
	UserId uuid.UUID
}

type LoginCommand struct {
	Nick     string `json:"nick" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LogoutCommand struct{}

type RegisterUserCommand struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Nick      string `json:"nick" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type TryGetUserBySessionIdCommand struct {
	SessionId uuid.UUID
}
