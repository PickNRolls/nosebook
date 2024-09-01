package friendship

import (
	"nosebook/src/application/services/auth"
	"nosebook/src/domain/friendship"
	"nosebook/src/errors"
	"nosebook/src/lib/clock"
)

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (this *Service) SendRequest(c *SendRequestCommand, a *auth.Auth) (*domainfriendship.FriendRequest, *errors.Error) {
	friendRequest := this.repository.
		RequesterId(a.UserId).
		ResponderId(c.ResponderId).
		FindOne()

	if friendRequest == nil {
		friendRequest = this.repository.
			RequesterId(c.ResponderId).
			ResponderId(a.UserId).
			FindOne()
	}

	if friendRequest == nil {
		newFriendRequest := domainfriendship.NewBuilder().
			RequesterId(a.UserId).
			ResponderId(c.ResponderId).
			Message(c.Message).
			CreatedAt(clock.Now()).
			RaiseCreatedEvent().
			Build()

		err := this.repository.Save(newFriendRequest)
		if err != nil {
			return nil, err
		}

		return newFriendRequest, nil
	}

	return nil, errors.New("Friendship Error", "Заявка в друзья уже отправлена")
}

func (this *Service) AcceptRequest(
	c *AcceptRequestCommand, a *auth.Auth,
) (*domainfriendship.FriendRequest, *errors.Error) {
	friendRequest := this.repository.
		RequesterId(c.RequesterId).
		ResponderId(a.UserId).
		OnlyNotAccepted().
		FindOne()

	if friendRequest == nil {
		return nil, errors.New("Friendship Error", "Заявка не найдена")
	}

	err := friendRequest.AcceptBy(a.UserId)
	if err != nil {
		return nil, err
	}

	err = this.repository.Save(friendRequest)
	if err != nil {
		return nil, err
	}

	return friendRequest, nil
}

func (this *Service) DenyRequest(c *DenyRequestCommand, a *auth.Auth) (*domainfriendship.FriendRequest, *errors.Error) {
	friendRequest := this.repository.
		RequesterId(c.RequesterId).
		ResponderId(a.UserId).
		OnlyNotAccepted().
		FindOne()

	if friendRequest == nil {
		return nil, errors.New("Friendship Error", "Заявка не найдена")
	}

	err := friendRequest.DenyBy(a.UserId)
	if err != nil {
		return nil, err
	}

	err = this.repository.Save(friendRequest)
	if err != nil {
		return nil, err
	}

	return friendRequest, nil
}

func (this *Service) RemoveFriend(c *RemoveFriendCommand, a *auth.Auth) (*domainfriendship.FriendRequest, *errors.Error) {
	request := this.repository.
		RequesterId(c.FriendId).
		ResponderId(a.UserId).
		OnlyAccepted().
		FindOne()

	if request == nil {
		request = this.repository.
			ResponderId(c.FriendId).
			RequesterId(a.UserId).
			OnlyAccepted().
			FindOne()
	}

	if request == nil {
		return nil, errors.New("Friendship Error", "Заявка не найдена")
	}

	err := request.RemoveBy(a.UserId)
	if err != nil {
		return nil, err
	}

	err = this.repository.Save(request)
	if err != nil {
		return nil, err
	}

	return request, nil
}
