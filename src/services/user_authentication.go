package services

import (
	"nosebook/src/domain/sessions"
	"nosebook/src/domain/users"
	"nosebook/src/errors"
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

func (s *UserAuthenticationService) RegisterUser(c *commands.RegisterUserCommand) (*auth.AuthResult, error) {
	existingUser := s.userRepo.FindByNick(c.Nick)
	if existingUser != nil {
		return nil, errors.New("NickError", "Логин занят")
	}

	passhash, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	user := users.NewUser(c.FirstName, c.LastName, c.Nick, string(passhash))
	user, err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	session, err := s.CreateSession(&commands.CreateSessionCommand{
		UserId: user.ID,
	})
	if err != nil {
		return nil, err
	}

	return &auth.AuthResult{
		User:    user,
		Session: session,
	}, nil
}

func (s *UserAuthenticationService) Login(c *commands.LoginCommand) (*auth.AuthResult, error) {
	existingUser := s.userRepo.FindByNick(c.Nick)
	if existingUser == nil {
		return nil, errors.New("NickError", "Пользователь с таким логином отсутствует")
	}

	err := bcrypt.CompareHashAndPassword([]byte(existingUser.Passhash), []byte(c.Password))
	if err != nil {
		return nil, errors.New("PasswordError", "Некорректный пароль")
	}

	session, err := s.CreateSession(&commands.CreateSessionCommand{
		UserId: existingUser.ID,
	})
	if err != nil {
		return nil, err
	}

	return &auth.AuthResult{
		User:    existingUser,
		Session: session,
	}, nil
}

func (s *UserAuthenticationService) Logout(a *auth.Auth) (*sessions.Session, error) {
	session := s.sessionRepo.FindById(a.SessionId)
	if session == nil {
		return nil, errors.New("LogoutError", "Сессии не существует")
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
		return nil, errors.New("SessionError", "Сессия не существует")
	}

	return s.userRepo.FindById(session.UserId), nil
}

func (s *UserAuthenticationService) MarkSessionActive(sessionId uuid.UUID) error {
	session := s.sessionRepo.FindById(sessionId)
	if session == nil {
		return errors.New("SessionError", "Сессия не существует")
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
