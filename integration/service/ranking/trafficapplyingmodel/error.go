package trafficapplyingmodel

import "errors"

// common errors
var (
	ErrNoValidTrafficQuerier   = errors.New("at least one of live traffic and historical speed querier should be valid")
	ErrEmptyRoute              = errors.New("route is nil")
	ErrEmptyLeg                = errors.New("leg is nil")
	ErrEmptyAnnotation         = errors.New("annotation is nil")
	ErrEmptyAnnotationMetadata = errors.New("annotation/metadata is nil")
)
