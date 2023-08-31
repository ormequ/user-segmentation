package service

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"user-segmentation/internal/entities/operations"
	"user-segmentation/internal/entities/segments"
)

var ErrInvalidDates = errors.New("invalid dates")

type ChangeErrors map[string]string

//go:generate go run github.com/vektra/mockery/v2@v2.33.1 --name=SegmentsRepo
type SegmentsRepo interface {
	Store(ctx context.Context, seg segments.Segment) error
	Delete(ctx context.Context, seg segments.Segment) error
	ChangeUserSegments(ctx context.Context, userID int64, add []segments.Segment, remove []segments.Segment) ChangeErrors
	GetUserSegments(ctx context.Context, userID int64) ([]segments.Segment, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.33.1 --name=HistoryRepo
type HistoryRepo interface {
	Get(ctx context.Context, year int, month int) ([]operations.Operation, error)
	Put(ctx context.Context, ops []operations.Operation) error
}

type Storage interface {
	Put(ctx context.Context)
}

type Service struct {
	Segments SegmentsRepo
	History  HistoryRepo
}

func (s Service) CreateSegment(ctx context.Context, slug string) error {
	seg, err := segments.New(slug)
	if err != nil {
		return err
	}
	return s.Segments.Store(ctx, seg)
}

func (s Service) DeleteSegment(ctx context.Context, slug string) error {
	seg, err := segments.New(slug)
	if err != nil {
		return err
	}
	return s.Segments.Delete(ctx, seg)
}

func createSegments(slugs []string) ([]segments.Segment, ChangeErrors) {
	res := make([]segments.Segment, len(slugs))
	errs := make(ChangeErrors)
	var err error
	for i := range res {
		res[i], err = segments.New(slugs[i])
		if err != nil {
			errs[slugs[i]] = err.Error()
		}
	}
	return res, errs
}

func (s Service) ChangeUserSegments(ctx context.Context, userID int64, add []string, remove []string) (ChangeErrors, error) {
	addSeg, errs := createSegments(add)
	rmSeg, errsRm := createSegments(remove)
	maps.Copy(errs, errsRm)
	if len(errs) != 0 {
		return errs, nil
	}
	errs = s.Segments.ChangeUserSegments(ctx, userID, addSeg, rmSeg)
	if len(errs) != 0 {
		return errs, nil
	}
	ops := make([]operations.Operation, 0, len(addSeg)+len(rmSeg))
	for i := range addSeg {
		op, _ := operations.New(userID, addSeg[i], operations.Add)
		ops = append(ops, op)
	}
	for i := range rmSeg {
		op, _ := operations.New(userID, rmSeg[i], operations.Remove)
		ops = append(ops, op)
	}
	return nil, s.History.Put(ctx, ops)
}

func (s Service) GetUserSegments(ctx context.Context, userID int64) ([]segments.Segment, error) {
	return s.Segments.GetUserSegments(ctx, userID)
}

func (s Service) GetHistory(ctx context.Context, year int, month int) ([][]string, error) {
	if month < 1 || month > 12 || year < 1970 {
		return nil, ErrInvalidDates
	}
	ops, err := s.History.Get(ctx, year, month)
	if err != nil {
		return nil, err
	}
	res := make([][]string, 0, len(ops)+1)
	res = append(res, []string{"User ID", "Segment", "Operation", "Timestamp UTC"})
	for i := range ops {
		opType := "add"
		if ops[i].Type == operations.Remove {
			opType = "remove"
		}
		res = append(res, []string{fmt.Sprint(ops[i].UserID), ops[i].Segment.Slug, opType, ops[i].Time.String()})
	}
	return res, nil
}

func New(seg SegmentsRepo, his HistoryRepo) Service {
	return Service{Segments: seg, History: his}
}
