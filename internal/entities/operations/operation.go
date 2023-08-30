package operations

import (
	"errors"
	"time"
	"user-segmentation/internal/entities/segments"
)

var (
	ErrIncorrectType = errors.New("incorrect type")
)

type Type int64

const (
	Add Type = iota
	Remove
)

type Operation struct {
	UserID  int64
	Segment segments.Segment
	Type    Type
	Time    time.Time
}

func New(userID int64, seg segments.Segment, opType Type) (Operation, error) {
	if opType != Add && opType != Remove {
		return Operation{}, ErrIncorrectType
	}
	return Operation{
		UserID:  userID,
		Segment: seg,
		Type:    opType,
		Time:    time.Now().UTC(),
	}, nil
}
