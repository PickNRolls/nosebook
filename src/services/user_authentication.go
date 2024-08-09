package services

import (
	"errors"
	"nosebook/src/domain/sessions"
	"nosebook/src/domain/users"
	"nosebook/src/services/auth"
	common_interfaces "nosebook/src/services/common/interfaces"
	"nosebook/src/services/user_authentication/commands"
	"nosebook/src/services/user_authentication/interfaces"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserAuthenticationService struct {
	userRepo    common_interfaces.UserRepository
	sessionRepo interfaces.SessionRepository
}

func NewUserAuthenticationService(userRepo common_interfaces.UserRepository, sessionRepo interfaces.SessionRepository) *UserAuthenticationService {
	return &UserAuthenticationService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *UserAuthenticationService) RegisterUser(c *commands.RegisterUserCommand) (*users.User, error) {
	existingUser := s.userRepo.FindByNick(c.Nick)
	if existingUser != nil {
		return nil, errors.New("Can't register user with such nickname.")
	}

	passhash, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	user := users.NewUser(c.FirstName, c.LastName, c.Nick, string(passhash))
	return s.userRepo.Create(user)
}

func (s *UserAuthenticationService) Login() (*users.User, error) {
	return users.NewUser("", "", "", ""), nil
}

func (s *UserAuthenticationService) Logout(a *auth.Auth) (*sessions.Session, error) {
	session := s.sessionRepo.FindById(a.SessionId)
	if session == nil {
		return nil, errors.New("Can't logout, session does not exist.")
	}

	session, err := s.sessionRepo.Remove(a.SessionId)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *UserAuthenticationService) CreateSession(c *commands.CreateSessionCommand) (*sessions.Session, error) {
	session := sessions.NewSession(c.UserId)
	session, err := s.sessionRepo.Create(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *UserAuthenticationService) TryGetUserBySessionId(c *commands.TryGetUserBySessionIdCommand) (*users.User, error) {
	session := s.sessionRepo.FindById(c.SessionId)
	if session == nil {
		return nil, errors.New("Invalid session id.")
	}

	return s.userRepo.FindById(session.UserId), nil
}

func (s *UserAuthenticationService) MarkSessionActive(sessionId uuid.UUID) error {
	session := s.sessionRepo.FindById(sessionId)
	if session == nil {
		return errors.New("Invalid session id.")
	}

	err := session.Refresh()
	if err != nil {
		return err
	}

	_, err = s.sessionRepo.Update(session)
	if err != nil {
		return nil
	}

	err = s.userRepo.UpdateActivity(session.UserId, time.Now())
	return err
}
