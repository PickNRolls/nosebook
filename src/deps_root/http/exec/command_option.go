package exec

import (
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Binding[C any] func(ginctx *gin.Context) (*C, error)

func JsonBinding[C any](ginctx *gin.Context) (*C, error) {
	var c C
	err := ginctx.ShouldBindJSON(&c)
	return &c, err
}

type CommandOption[C any, T any] interface {
	AvoidBinding() bool
	Binding() Binding[C]
	Tracer() trace.Tracer
	OutputMapper() func(out T) any
}

type avoidBindingOption[C any, T any] struct{}

func (this *avoidBindingOption[C, T]) AvoidBinding() bool            { return true }
func (this *avoidBindingOption[C, T]) Binding() Binding[C]           { return nil }
func (this *avoidBindingOption[C, T]) Tracer() trace.Tracer          { return nil }
func (this *avoidBindingOption[C, T]) OutputMapper() func(out T) any { return nil }

func WithAvoidBinding[C any, T any]() CommandOption[C, T] {
	return &avoidBindingOption[C, T]{}
}

type tracerOption[C any, T any] struct {
	tracer trace.Tracer
}

func (this *tracerOption[C, T]) AvoidBinding() bool            { return false }
func (this *tracerOption[C, T]) Binding() Binding[C]           { return nil }
func (this *tracerOption[C, T]) Tracer() trace.Tracer          { return this.tracer }
func (this *tracerOption[C, T]) OutputMapper() func(out T) any { return nil }

func WithTracer[C any, T any](tracer trace.Tracer) func() CommandOption[C, T] {
	return func() CommandOption[C, T] {
		return &tracerOption[C, T]{
			tracer: tracer,
		}
	}
}

type outputMapperOption[C any, T any] struct {
	mapper func(out T) any
}

func (this *outputMapperOption[C, T]) AvoidBinding() bool            { return false }
func (this *outputMapperOption[C, T]) Binding() Binding[C]           { return nil }
func (this *outputMapperOption[C, T]) Tracer() trace.Tracer          { return nil }
func (this *outputMapperOption[C, T]) OutputMapper() func(out T) any { return this.mapper }

func WithMapper[C any, T any](mapper func(out T) any) func() CommandOption[C, T] {
	return func() CommandOption[C, T] {
		return &outputMapperOption[C, T]{
			mapper: mapper,
		}
	}
}

func WithUuidMapper[C any]() CommandOption[C, uuid.UUID] {
	return WithMapper[C](func(out uuid.UUID) any {
		return struct {
			Id uuid.UUID `json:"id"`
		}{
			Id: out,
		}
	})()
}

type Writer[Pointer any] interface {
	Write(filename string, bytes []byte) error
	*Pointer
}

type fileBindingOption[C any, T any] struct {
	binding Binding[C]
}

func fileBinding[C any, PC Writer[C]](ginctx *gin.Context) (*C, error) {
	file, err := ginctx.FormFile("file")
	if err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(src)

	var command C
	var pointer PC = &command
	err = pointer.Write(file.Filename, buf.Bytes())
	if err != nil {
		return nil, err
	}

	return pointer, nil
}

func (this *fileBindingOption[C, T]) AvoidBinding() bool            { return false }
func (this *fileBindingOption[C, T]) Binding() Binding[C]           { return this.binding }
func (this *fileBindingOption[C, T]) Tracer() trace.Tracer          { return nil }
func (this *fileBindingOption[C, T]) OutputMapper() func(out T) any { return nil }

func WithFileBinding[C any, V any, PC Writer[C]]() CommandOption[C, V] {
	return &fileBindingOption[C, V]{
		binding: fileBinding[C, PC],
	}
}
