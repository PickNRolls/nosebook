package domainlike

import "github.com/google/uuid"

type ResourceType string

const (
	POST_RESOURCE    ResourceType = "post"
	COMMENT_RESOURCE ResourceType = "comment"
)

type Resource interface {
	Type() ResourceType
	Id() uuid.UUID
	Owner() Owner
}
