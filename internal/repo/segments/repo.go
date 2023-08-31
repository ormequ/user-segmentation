package segments

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
	"user-segmentation/internal/entities/segments"
	"user-segmentation/internal/logger"
	"user-segmentation/internal/repo"
	"user-segmentation/internal/service"
)

const (
	constrSegmentID      = "user_segments_segment_id_key"
	constrSegmentExists  = "segments_slug_key"
	constrRelationExists = "user_segments_pkey"
)

type Repo struct {
	db *pgx.Conn
}

func (r Repo) Store(ctx context.Context, seg segments.Segment) error {
	const fn = "repo.segments.Store"
	const query = "INSERT INTO segments (slug) VALUES ($1)"
	_, err := r.db.Exec(ctx, query, seg.Slug)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == constrSegmentExists {
			err = repo.ErrSegmentAlreadyExists
		} else {
			logger.InternalErr(ctx, err, fn)
		}
	}
	return err
}

func (r Repo) Delete(ctx context.Context, seg segments.Segment) error {
	const fn = "repo.segments.Remove"
	const query = "DELETE FROM segments WHERE slug=$1"
	cmd, err := r.db.Exec(ctx, query, seg.Slug)
	if errors.Is(err, pgx.ErrNoRows) || cmd.RowsAffected() == 0 {
		return repo.ErrSegmentNotFound
	}
	if err != nil {
		logger.InternalErr(ctx, err, fn)
	}
	return err
}

var ErrChangingInternal = errors.New("internal error")

func (r Repo) ChangeUserSegments(ctx context.Context, userID int64, add []segments.Segment, remove []segments.Segment) service.ChangeErrors {
	const fn = "repo.segments.ChangeUserSegments"
	const addQuery = `INSERT INTO user_segments (user_id, segment_id) VALUES ($1, (SELECT id FROM segments WHERE slug=$2))`
	const rmQuery = `DELETE FROM user_segments WHERE user_id=$1 AND segment_id=(SELECT id FROM segments WHERE slug=$2)`
	batch := &pgx.Batch{}
	for _, seg := range remove {
		batch.Queue(rmQuery, userID, seg.Slug)
	}
	for _, seg := range add {
		batch.Queue(addQuery, userID, seg.Slug)
	}
	br := r.db.SendBatch(ctx, batch)
	defer func(br pgx.BatchResults) {
		_ = br.Close()
	}(br)
	errs := make(service.ChangeErrors)
	for _, seg := range remove {
		_, err := br.Exec()
		logger.Log(ctx).Debug("exec rm", slog.Any("error", err), slog.String("segment", seg.Slug))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = repo.ErrRelationNotFound
			} else {
				logger.InternalErr(ctx, err, fn)
				err = ErrChangingInternal
			}
			errs[seg.Slug] = err.Error()
		}
	}
	for _, seg := range add {
		_, err := br.Exec()
		logger.Log(ctx).Debug("exec add", slog.Any("error", err), slog.String("segment", seg.Slug))
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.ConstraintName == constrSegmentID {
				err = repo.ErrSegmentNotFound
			} else if pgErr.ConstraintName == constrRelationExists {
				err = repo.ErrRelationExists
			} else {
				logger.InternalErr(ctx, err, fn)
				err = ErrChangingInternal
			}
			errs[seg.Slug] = err.Error()
		}
	}
	return errs
}

func (r Repo) GetUserSegments(ctx context.Context, userID int64) ([]segments.Segment, error) {
	const fn = "repo.segments.GetUserSegments"
	const query = "SELECT slug FROM segments WHERE id=ANY (SELECT segment_id FROM user_segments WHERE user_id=$1)"
	rows, err := r.db.Query(ctx, query, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return []segments.Segment{}, repo.ErrNoSegments
	}
	defer rows.Close()
	var res []segments.Segment
	for rows.Next() {
		var slug string
		err := rows.Scan(&slug)
		if err != nil {
			logger.InternalErr(ctx, err, fn)
			return nil, err
		}
		res = append(res, segments.Segment{Slug: slug})
	}
	return res, nil
}

func New(db *pgx.Conn) Repo {
	return Repo{db: db}
}
