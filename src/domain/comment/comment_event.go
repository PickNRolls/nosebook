package domaincomment

type CommentEventType string

const (
	CREATED CommentEventType = "CREATED"
	REMOVED CommentEventType = "REMOVED"
)

type CommentEvent interface {
	Type() CommentEventType
}
