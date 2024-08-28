package errors

import (
	"nosebook/src/errors"
)

type RouteNotFoundError = errors.Error

func NewRouteNotFoundError() *RouteNotFoundError {
	return errors.New("NotFound", "Route not found")
}
