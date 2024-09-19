package userauth

import (
	"context"
	"nosebook/src/application/services/auth"
	"nosebook/src/domain/sessions"
	"nosebook/src/domain/user"
	"nosebook/src/errors"
	"nosebook/src/lib/clock"
	commandresult "nosebook/src/lib/command_result"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
  tracer trace.Tracer
}

func New(userRepo UserRepository, sessionRepo SessionRepository, tracer trace.Tracer) *Service {
	return &Service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
    tracer: tracer,
	}
}

func (this *Service) RegisterUser(c RegisterUserCommand, a *auth.Auth) *commandresult.Result {
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

func (this *Service) Login(c LoginCommand, a *auth.Auth) *commandresult.Result {
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

func (this *Service) Logout(c LogoutCommand, a *auth.Auth) *commandresult.Result {
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

func (this *Service) TryGetUserBySessionId(c TryGetUserBySessionIdCommand) (*domainuser.User, error) {
	session := this.sessionRepo.FindById(c.SessionId)
	if session == nil {
		return nil, errors.New("SessionError", "Сессия не существует")
	}

	return this.userRepo.FindById(session.UserId), nil
}

func (this *Service) MarkSessionActive(parent context.Context, sessionId uuid.UUID) error {
  ctx, span := this.tracer.Start(parent, "user_auth_service.mark_session_active")
  defer span.End()
  
  _, span = this.tracer.Start(ctx, "session_repo.find_by_id")
	session := this.sessionRepo.FindById(sessionId)
  span.End()

	if session == nil {
		return errors.New("SessionError", "Сессия не существует")
	}

	err := session.Refresh()
	if err != nil {
		return err
	}

  _, span = this.tracer.Start(ctx, "session_repo.update")
	_, err = this.sessionRepo.Update(session)
  span.End()
	if err != nil {
		return nil
	}

  _, span = this.tracer.Start(ctx, "user_repo.update_activity")
	err = this.userRepo.UpdateActivity(session.UserId, clock.Now())
  span.End()
	return err
}
