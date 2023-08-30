package tests

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"user-segmentation/internal/entities/operations"
	"user-segmentation/internal/entities/segments"
	"user-segmentation/internal/repo/history"
)

func TestHistoryBasic(t *testing.T) {
	_, _ = db.Exec(context.Background(), "TRUNCATE operations")
	client := setupClient()
	_, err := client.getHistory(2030, 100)
	require.ErrorIs(t, err, ErrBadRequest)
	for i := 0; i < 10; i++ {
		_, err = client.createSegment(fmt.Sprintf("seg-%d", i))
		require.NoError(t, err)
	}

	res, err := client.getHistory(2030, 1)
	require.NoError(t, err)
	require.Len(t, res, 1)

	_, err = client.changeUserSegments(1, []string{"seg-1", "seg-2"}, []string{})
	require.NoError(t, err)
	_, err = client.changeUserSegments(2, []string{"seg-1", "seg-2"}, []string{})
	require.NoError(t, err)

	now := time.Now().UTC()
	res, err = client.getHistory(now.Year(), int(now.Month()))
	require.NoError(t, err)
	require.Len(t, res, 5)
	csvSeg := make(map[string]struct{}, len(res))
	csvUsr := make(map[string]struct{}, len(res))
	csvType := make(map[string]struct{}, len(res))
	for i := 1; i < len(res); i++ {
		csvUsr[res[i][0]] = struct{}{}
		csvSeg[res[i][1]] = struct{}{}
		csvType[res[i][2]] = struct{}{}
	}
	require.Len(t, csvUsr, 2)
	require.Len(t, csvSeg, 2)
	require.Len(t, csvType, 1)
}

func TestHistoryDates(t *testing.T) {
	_, _ = db.Exec(context.Background(), "TRUNCATE operations")
	client := setupClient()
	_, err := client.createSegment("slug")
	require.NoError(t, err)
	repo := history.New(db)
	now := time.Now().UTC()
	err = repo.Put(context.Background(), []operations.Operation{
		{
			UserID:  0,
			Segment: segments.Segment{Slug: "slug"},
			Type:    0,
			Time:    now.Add(-time.Hour * 24 * 60),
		},
		{
			UserID:  0,
			Segment: segments.Segment{Slug: "slug"},
			Type:    0,
			Time:    now.Add(-time.Hour * 24 * 60),
		},
	})
	require.NoError(t, err)
	res, err := client.getHistory(now.Year(), int(now.Month()))
	require.NoError(t, err)
	require.Len(t, res, 1)
}
