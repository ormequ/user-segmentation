package service

import (
	"context"
	"maps"
	"user-segmentation/internal/entities/operations"
	"user-segmentation/internal/entities/segments"
)

type ChangeErrors map[string]string

type SegmentsRepo interface {
	Store(ctx context.Context, seg segments.Segment) error
	Delete(ctx context.Context, seg segments.Segment) error
	ChangeUserSegments(ctx context.Context, userID int64, add []segments.Segment, remove []segments.Segment) ChangeErrors
	GetUserSegments(ctx context.Context, userID int64) ([]segments.Segment, error)
}

type HistoryRepo interface {
	Get(ctx context.Context, userID int64) ([]operations.Operation, error)
	Put(ctx context.Context, ops []operations.Operation) error
}

type Service struct {
	repo SegmentsRepo
}

func (s Service) CreateSegment(ctx context.Context, slug string) error {
	seg, err := segments.New(slug)
	if err != nil {
		return err
	}
	return s.repo.Store(ctx, seg)
}

func (s Service) DeleteSegment(ctx context.Context, slug string) error {
	seg, err := segments.New(slug)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, seg)
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

func (s Service) ChangeUserSegments(ctx context.Context, userID int64, add []string, remove []string) ChangeErrors {
	addToRepo, errs := createSegments(add)
	rmToRepo, errsRm := createSegments(remove)
	maps.Copy(errs, errsRm)
	if len(errs) != 0 {
		return errs
	}
	errs = s.repo.ChangeUserSegments(ctx, userID, addToRepo, rmToRepo)
	return errs
}

func (s Service) GetUserSegments(ctx context.Context, userID int64) ([]segments.Segment, error) {
	return s.repo.GetUserSegments(ctx, userID)
}

func New(repo SegmentsRepo) Service {
	return Service{repo: repo}
}
