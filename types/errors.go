package types

import "errors"

var (
	ErrDuplicatedModule  = errors.New("duplicated node")
	ErrInvalidVersion    = errors.New("invalid version")
	ErrInvalidName       = errors.New("invalid name")
	ErrNotFindDependency = errors.New("not find dependency")
)
