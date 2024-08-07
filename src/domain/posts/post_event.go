package posts

type PostEventType string

const (
	LIKED   PostEventType = "LIKED"
	UNLIKED PostEventType = "UNLIKED"
)

type PostEvent interface {
	Type() PostEventType
}
