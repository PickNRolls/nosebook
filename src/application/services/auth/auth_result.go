package auth

import (
	"nosebook/src/domain/sessions"
	"nosebook/src/domain/user"
)

type AuthResult struct {
	User    *domainuser.User  `json:"user"`
	Session *sessions.Session `json:"session"`
}
