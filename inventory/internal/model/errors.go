package model

import (
	"errors"
)

var (
	ErrPartNotFound           = errors.New("part not found")
	ErrInvalidPart            = errors.New("invalid part data")
	ErrRepositoryOperation    = errors.New("repository operation failed")
	ErrConcurrentModification = errors.New("concurrent modification detected")
)
