package services

import (
	"errors"
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

func (s *FriendshipService) AcceptFriendRequest(c *commands.AcceptFriendRequestCommand, a *auth.Auth) (*friendship.FriendRequest, error) {
	friendRequest := s.userFriendsRepo.FindByBoth(c.RequesterId, a.UserId)
	if friendRequest == nil {
		return nil, errors.New("No such friend request.")
	}

	friendRequest.Accepted = true
	friendRequest.Viewed = true

	_, err := s.userFriendsRepo.Update(friendRequest)
	if err != nil {
		return nil, err
	}

	return friendRequest, nil
}

func (s *FriendshipService) DenyFriendRequest(c *commands.SendFriendRequestCommand, a *auth.Auth) (interface{}, error) {
	return nil, nil
}

func (s *FriendshipService) RemoveFriend(c *commands.SendFriendRequestCommand, a *auth.Auth) (interface{}, error) {
	return nil, nil
}
