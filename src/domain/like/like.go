package domainlike

import (
	"fmt"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type Like struct {
	Resource Resource
	Owner    Owner
	Event    Event
	value    bool
}

func New() *Like {
	return &Like{}
}

func (this *Like) WithValue(value bool) *Like {
	this.value = value
	return this
}

func (this *Like) WithOwner(owner Owner) *Like {
	this.Owner = owner
	return this
}

func (this *Like) setResource(resource Resource) (*Like, *errors.Error) {
	if this.Resource != nil {
		return this, errors.New(
			"Like Error",
			fmt.Sprintf("Попытка присвоить ресурс типа '%s' к лайку с уже присовенным ресурсом типа '%s'",
				resource.Type(),
				this.Resource.Type(),
			),
		)
	}

	this.Resource = resource
	return this, nil
}

func (this *Like) WithPostId(id uuid.UUID) (*Like, *errors.Error) {
	return this.setResource(NewPostResource(id))
}

func (this *Like) WithCommentId(id uuid.UUID) (*Like, *errors.Error) {
	return this.setResource(NewCommentResource(id))
}

func (this *Like) Toggle() *errors.Error {
	if this.Resource == nil {
		return errors.New("Like Error", "У лайка отсутствует ресурс")
	}

	this.value = !this.value

	if this.Event == nil {
		if this.value {
			this.Event = NewLikeEvent()
		} else {
			this.Event = NewUnlikeEvent()
		}
	} else {
		this.Event = nil
	}

	return nil
}
