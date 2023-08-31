package history

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
	"user-segmentation/internal/entities/operations"
	"user-segmentation/internal/logger"
	"user-segmentation/internal/repo"
)

const constrSegmentID = "operations_segment_id_key"

type Repo struct {
	db *pgx.Conn
}

func (r Repo) Get(ctx context.Context, year int, month int) ([]operations.Operation, error) {
	const fn = "repo.history.Get"
	const query = `SELECT user_id, segments.slug, type, time FROM operations 
                   JOIN segments ON operations.segment_id = segments.id
                   WHERE time BETWEEN $1 AND $2`
	m := time.Month(month)
	minTime := time.Date(year, m, 1, 0, 0, 0, 0, time.UTC)
	maxTime := time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC)
	rows, err := r.db.Query(ctx, query, minTime, maxTime)
	if errors.Is(err, pgx.ErrNoRows) {
		return []operations.Operation{}, nil
	}
	defer rows.Close()
	var res []operations.Operation
	for rows.Next() {
		op := operations.Operation{}
		err := rows.Scan(&op.UserID, &op.Segment.Slug, &op.Type, &op.Time)
		if err != nil {
			logger.InternalErr(ctx, err, fn)
			return nil, err
		}
		res = append(res, op)
	}
	return res, nil
}

func (r Repo) Put(ctx context.Context, ops []operations.Operation) error {
	const fn = "repo.history.Put"
	const query = `INSERT INTO operations (user_id, segment_id, type, time) VALUES 
				   ($1, (SELECT id FROM segments WHERE slug=$2), $3, $4)`
	batch := &pgx.Batch{}
	for _, op := range ops {
		batch.Queue(query, op.UserID, op.Segment.Slug, op.Type, op.Time)
	}
	br := r.db.SendBatch(ctx, batch)
	defer func(br pgx.BatchResults) {
		_ = br.Close()
	}(br)
	for range ops {
		_, err := br.Exec()
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.ConstraintName == constrSegmentID {
				err = repo.ErrSegmentNotFound
			} else {
				logger.InternalErr(ctx, err, fn)
			}
			return err
		}
	}
	return nil
}

func New(db *pgx.Conn) Repo {
	return Repo{db}
}
