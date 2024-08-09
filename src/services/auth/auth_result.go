package auth

import (
	"nosebook/src/domain/sessions"
	"nosebook/src/domain/users"
)

type AuthResult struct {
	User    *users.User       `json:"user"`
	Session *sessions.Session `json:"session"`
}
