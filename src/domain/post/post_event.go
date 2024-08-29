package domainpost

type PostEventType string

const (
	CREATED PostEventType = "CREATED"
	REMOVED PostEventType = "REMOVED"
	EDITED  PostEventType = "EDITED"
)

type PostEvent interface {
	Type() PostEventType
}
