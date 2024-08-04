package services

import (
	"nosebook/src/domain/friendship"
	"nosebook/src/services/auth"
	"nosebook/src/services/friendship/commands"
	"nosebook/src/services/friendship/interfaces"
)

type FriendshipService struct {
	userFriendsRepo interfaces.UserFriendsRepository
}

func NewFriendshipService(userFriendsRepo interfaces.UserFriendsRepository) *FriendshipService {
	return &FriendshipService{
		userFriendsRepo: userFriendsRepo,
	}
}

func (s *FriendshipService) SendFriendRequest(c *commands.SendFriendRequestCommand, a *auth.Auth) (*friendship.FriendRequest, error) {
	requesterId := a.UserId
	friendRequest := s.userFriendsRepo.FindByBoth(requesterId, c.ResponderId)
	if friendRequest == nil {
		friendRequest = s.userFriendsRepo.FindByBoth(c.ResponderId, requesterId)
	}

	if friendRequest == nil {
		newFriendRequest := friendship.NewFriendRequest(requesterId, c.ResponderId, c.Message)
		return s.userFriendsRepo.Create(newFriendRequest)
	}

	return nil, nil
}

func (s *FriendshipService) AcceptFriendRequest(c *commands.SendFriendRequestCommand, a auth.Auth) (interface{}, error) {
	return nil, nil
}

func (s *FriendshipService) DenyFriendRequest(c *commands.SendFriendRequestCommand, a auth.Auth) (interface{}, error) {
	return nil, nil
}

func (s *FriendshipService) RemoveFriend(c *commands.SendFriendRequestCommand, a auth.Auth) (interface{}, error) {
	return nil, nil
}
