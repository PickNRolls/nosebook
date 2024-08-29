package reqcontext

import (
	"net/http"
	"nosebook/src/domain/users"
	"nosebook/src/errors"
	infraerrors "nosebook/src/infra/errors"
	"nosebook/src/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReqContext struct {
	ctx      *gin.Context
	errs     []*errors.Error
	errsRead bool
}

func From(ctx *gin.Context) *ReqContext {
	return &ReqContext{
		ctx:  ctx,
		errs: make([]*errors.Error, 0),
	}
}

func (this *ReqContext) SetUser(user *users.User) {
	this.ctx.Set("user", user)
}

func (this *ReqContext) User() *users.User {
	userAny, ok := this.ctx.Get("user")
	if !ok {
		return nil
	}

	user, ok := userAny.(*users.User)
	if !ok {
		return nil
	}

	return user
}

func (this *ReqContext) UserOrForbidden() *users.User {
	user := this.User()

	if user == nil {
		this.ctx.Status(http.StatusForbidden)
		this.ctx.Error(infraerrors.NewNotAuthenticatedError())
		this.ctx.Abort()
	}

	return user
}

func (this *ReqContext) SetSessionId(id uuid.UUID) {
	this.ctx.Set("sessionId", id)
}

func (this *ReqContext) SessionId() uuid.UUID {
	unknown, ok := this.ctx.Get("sessionId")
	if !ok {
		return uuid.Nil
	}

	sessionId, ok := unknown.(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return sessionId
}

func (this *ReqContext) Auth() *auth.Auth {
	return auth.From(this.User(), this.SessionId())
}

func (this *ReqContext) SetResponseData(data any) {
	this.ctx.Set("data", data)
}

func (this *ReqContext) ResponseData() any {
	data, exists := this.ctx.Get("data")
	if !exists {
		return nil
	}

	return data
}

func (this *ReqContext) Errors() []*errors.Error {
	if this.errsRead {
		return this.errs
	}

	for _, ginErr := range this.ctx.Errors {
		err := errors.From(ginErr.Err)
		this.errs = append(this.errs, err)
	}

	this.errsRead = true
	return this.errs
}
