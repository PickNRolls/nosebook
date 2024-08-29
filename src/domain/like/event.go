package domainlike

type EventType string

const (
	LIKED   EventType = "LIKED"
	UNLIKED EventType = "UNLIKED"
)

type Event interface {
	Type() EventType
}
