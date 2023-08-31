package repo

import "errors"

var (
	ErrSegmentAlreadyExists = errors.New("segment already exists")
	ErrSegmentNotFound      = errors.New("segment not found")
	ErrRelationNotFound     = errors.New("user is not in this segment")
	ErrNoSegments           = errors.New("users not found")
	ErrRelationExists       = errors.New("relation already exists")
)
