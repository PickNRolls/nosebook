package generics

import (
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type batchQueryEntry[T any] struct {
	Id     uuid.UUID
	Result SingleQueryResult[T]
}

type BatchQueryResult[T any] struct {
	Err     *errors.Error
	Results []*batchQueryEntry[T]
}

func (result *BatchQueryResult[T]) HasEntry(id uuid.UUID) bool {
	for _, entry := range result.Results {
		if entry.Id == id {
			return true
		}
	}

	return false
}

func (result *BatchQueryResult[T]) EntryById(id uuid.UUID) *batchQueryEntry[T] {
	for _, entry := range result.Results {
		if entry.Id == id {
			return entry
		}
	}

	return nil
}

func (result *BatchQueryResult[T]) Entries() []*batchQueryEntry[T] {
	return result.Results
}

func (result *BatchQueryResult[T]) AddEntryOnce(id uuid.UUID) {
	singleResult := result.SingleResultOf(id)

	if singleResult == nil {
		sr := &batchQueryEntry[T]{
			Id:     id,
			Result: SingleQueryResult[T]{},
		}
		result.Results = append(result.Results, sr)
	}
}

func (result *BatchQueryResult[T]) SingleResultOf(id uuid.UUID) *SingleQueryResult[T] {
	var singleResult *SingleQueryResult[T]
	for _, res := range result.Results {
		if res.Id == id {
			singleResult = &res.Result
		}
	}

	return singleResult
}
