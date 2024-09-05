package domainchat

import "github.com/google/uuid"

type Permissions interface {
	CanJoinBy(chat *Chat, userId uuid.UUID) *Error
	CanSendMessageBy(chat *Chat, userId uuid.UUID) *Error
}

type defaultPermissions struct{}

func (this *defaultPermissions) CanJoinBy(chat *Chat, userId uuid.UUID) *Error {
	if chat.Private {
		if len(chat.MemberIds) < 2 {
			return nil
		}
		return newError("Вы не можете присоединиться к чату")
	}

	return nil
}

func (this *defaultPermissions) CanSendMessageBy(chat *Chat, userId uuid.UUID) *Error {
	isMember := false
	for _, memberId := range chat.MemberIds {
		if memberId == userId {
			isMember = true
			break
		}
	}

	if !isMember {
		return newError("Только участники чата могут отправлять в него сообщения")
	}

	return nil
}
