package comments

type CommentEventType string

const (
	LIKED   CommentEventType = "LIKED"
	UNLIKED CommentEventType = "UNLIKED"
)

type CommentEvent interface {
	Type() CommentEventType
}
