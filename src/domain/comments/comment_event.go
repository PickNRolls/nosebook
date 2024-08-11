package comments

type CommentEventType string

const (
	LIKED   CommentEventType = "LIKED"
	UNLIKED CommentEventType = "UNLIKED"
	REMOVED CommentEventType = "REMOVED"
)

type CommentEvent interface {
	Type() CommentEventType
}
