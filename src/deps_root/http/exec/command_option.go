package exec

import (
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type CommandOption[T any] interface {
	AvoidBinding() bool
	Tracer() trace.Tracer
	OutputMapper() func(out T) any
}

type avoidBindingOption[T any] struct{}

func (this *avoidBindingOption[T]) AvoidBinding() bool            { return true }
func (this *avoidBindingOption[T]) Tracer() trace.Tracer          { return nil }
func (this *avoidBindingOption[T]) OutputMapper() func(out T) any { return nil }

func WithAvoidBinding[T any]() CommandOption[T] {
	return &avoidBindingOption[T]{}
}

type tracerOption[T any] struct {
	tracer trace.Tracer
}

func (this *tracerOption[T]) AvoidBinding() bool            { return false }
func (this *tracerOption[T]) Tracer() trace.Tracer          { return this.tracer }
func (this *tracerOption[T]) OutputMapper() func(out T) any { return nil }

func WithTracer[T any](tracer trace.Tracer) func() CommandOption[T] {
	return func() CommandOption[T] {
		return &tracerOption[T]{
			tracer: tracer,
		}
	}
}

type outputMapperOption[T any] struct {
	mapper func(out T) any
}

func (this *outputMapperOption[T]) AvoidBinding() bool            { return false }
func (this *outputMapperOption[T]) Tracer() trace.Tracer          { return nil }
func (this *outputMapperOption[T]) OutputMapper() func(out T) any { return this.mapper }

func WithMapper[T any](mapper func(out T) any) func() CommandOption[T] {
	return func() CommandOption[T] {
		return &outputMapperOption[T]{
			mapper: mapper,
		}
	}
}

func WithUuidMapper() CommandOption[uuid.UUID] {
	return WithMapper(func(out uuid.UUID) any {
		return struct {
			Id uuid.UUID `json:"id"`
		}{
			Id: out,
		}
	})()
}
