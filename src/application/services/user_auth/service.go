package userauth

import (
	"nosebook/src/application/services/auth"
	"nosebook/src/domain/sessions"
	"nosebook/src/domain/user"
	"nosebook/src/errors"
	"nosebook/src/lib/clock"
	commandresult "nosebook/src/lib/command_result"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
}

func New(userRepo UserRepository, sessionRepo SessionRepository) *Service {
	return &Service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (this *Service) RegisterUser(c *RegisterUserCommand, a *auth.Auth) *commandresult.Result {
	existingUser := this.userRepo.FindByNick(c.Nick)
	if existingUser != nil {
		return commandresult.Fail(errors.New("NickError", "Логин занят"))
	}

	passhash, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.MinCost)
	if err != nil {
		return commandresult.Fail(errors.From(err))
	}

	user := domainuser.New(c.FirstName, c.LastName, c.Nick, string(passhash))
	user, err = this.userRepo.Create(user)
	if err != nil {
		return commandresult.Fail(errors.From(err))
	}

	session, error := this.createSession(&CreateSessionCommand{
		UserId: user.Id,
	})
	if error != nil {
		return commandresult.Fail(errors.From(error))
	}

	return commandresult.Ok().WithData(&auth.AuthResult{
		User:    user,
		Session: session,
	})
}

func (this *Service) Login(c *LoginCommand, a *auth.Auth) *commandresult.Result {
	existingUser := this.userRepo.FindByNick(c.Nick)
	if existingUser == nil {
		return commandresult.Fail(errors.New("NickError", "Пользователь с таким логином отсутствует"))
	}

	err := bcrypt.CompareHashAndPassword([]byte(existingUser.Passhash), []byte(c.Password))
	if err != nil {
		return commandresult.Fail(errors.New("PasswordError", "Некорректный пароль"))
	}

	session, error := this.createSession(&CreateSessionCommand{
		UserId: existingUser.Id,
	})
	if error != nil {
		return commandresult.Fail(errors.From(error))
	}

	return commandresult.Ok().WithData(&auth.AuthResult{
		User:    existingUser,
		Session: session,
	})
}

func (this *Service) Logout(c *LogoutCommand, a *auth.Auth) *commandresult.Result {
	session := this.sessionRepo.FindById(a.SessionId)
	if session == nil {
		return commandresult.Fail(errors.New("LogoutError", "Сессии не существует"))
	}

	session, err := this.sessionRepo.Remove(a.SessionId)
	if err != nil {
		return commandresult.Fail(errors.From(err))
	}

	return commandresult.Ok().WithData(session)
}

func (this *Service) createSession(c *CreateSessionCommand) (*sessions.Session, *errors.Error) {
	session := sessions.NewSession(c.UserId)
	session, err := this.sessionRepo.Create(session)
	if err != nil {
		return nil, errors.From(err)
	}

	return session, nil
}

func (this *Service) TryGetUserBySessionId(c *TryGetUserBySessionIdCommand) (*domainuser.User, error) {
	session := this.sessionRepo.FindById(c.SessionId)
	if session == nil {
		return nil, errors.New("SessionError", "Сессия не существует")
	}

	return this.userRepo.FindById(session.UserId), nil
}

func (this *Service) MarkSessionActive(sessionId uuid.UUID) error {
	session := this.sessionRepo.FindById(sessionId)
	if session == nil {
		return errors.New("SessionError", "Сессия не существует")
	}

	err := session.Refresh()
	if err != nil {
		return err
	}

	_, err = this.sessionRepo.Update(session)
	if err != nil {
		return nil
	}

	err = this.userRepo.UpdateActivity(session.UserId, clock.Now())
	return err
}
