package interfaces

type EventDispatcher[E any] interface {
	Events() []E
}
