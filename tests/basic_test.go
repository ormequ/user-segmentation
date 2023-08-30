package tests

import (
	"crypto/rand"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
	"math/big"
	"testing"
	httpserver "user-segmentation/internal/api/http"
	"user-segmentation/internal/entities/segments"
	"user-segmentation/internal/repo"
)

func randString(sz int) string {
	res := make([]byte, sz)
	if _, err := rand.Read(res); err != nil {
		return ""
	}
	return string(res)
}

func TestCreateDeleteSegment(t *testing.T) {
	client := setupClient()
	res, err := client.createSegment("")
	assert.Equal(t, res.Error, httpserver.ErrInvalidRequest.Error())
	require.ErrorIs(t, err, ErrBadRequest)
	require.False(t, res.Data.Done)

	res, err = client.createSegment(randString(300))
	assert.Equal(t, res.Error, segments.ErrSlugToLong.Error())
	require.ErrorIs(t, err, ErrBadRequest)
	require.False(t, res.Data.Done)

	res, err = client.createSegment("test-segment")
	require.NoError(t, err)
	require.True(t, res.Data.Done)

	res, err = client.deleteSegment("")
	assert.Equal(t, res.Error, httpserver.ErrInvalidRequest.Error())
	require.ErrorIs(t, err, ErrBadRequest)
	require.False(t, res.Data.Done)

	res, err = client.deleteSegment(randString(300))
	assert.Equal(t, res.Error, segments.ErrSlugToLong.Error())
	assert.ErrorIs(t, err, ErrBadRequest)
	require.False(t, res.Data.Done)

	res, err = client.deleteSegment("no-segment")
	assert.Equal(t, res.Error, repo.ErrSegmentNotFound.Error())
	assert.ErrorIs(t, err, ErrNotFound)
	require.False(t, res.Data.Done)

	res, err = client.deleteSegment("test-segment")
	require.NoError(t, err)
	require.True(t, res.Data.Done)
}

func randInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

func TestUserSegments(t *testing.T) {
	const segLen = 100
	client := setupClient()
	seg := make([]string, segLen)
	was := make(map[string]struct{}, segLen)
	for i := range seg {
		s := gofakeit.UUID()
		for _, ok := was[s]; ok; {
			s = gofakeit.UUID()
		}
		seg[i] = s
		was[s] = struct{}{}
		res, err := client.createSegment(seg[i])
		require.NoError(t, err)
		require.Empty(t, res.Error)
		require.True(t, res.Data.Done)
	}

	resGet, err := client.getUserSegments(123)
	require.NoError(t, err)
	require.Empty(t, resGet.Error)
	require.Empty(t, resGet.Data)

	res, err := client.changeUserSegments(123, []string{}, []string{})
	require.NoError(t, err)
	require.True(t, res.Data.Done)
	require.Empty(t, res.Data.Errors)

	n := randInt(segLen)
	addSet := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		addSet[seg[randInt(segLen)]] = struct{}{}
	}
	add := make([]string, 0, len(addSet))
	for k := range addSet {
		add = append(add, k)
	}

	res, err = client.changeUserSegments(1, add, []string{})
	require.NoError(t, err)
	require.True(t, res.Data.Done)
	require.Empty(t, res.Data.Errors)

	res, err = client.changeUserSegments(1, []string{}, []string{})
	require.NoError(t, err)
	require.True(t, res.Data.Done)
	require.Empty(t, res.Data.Errors)

	res, err = client.changeUserSegments(1, []string{}, []string{"i am not exists!"})
	require.NoError(t, err)
	require.True(t, res.Data.Done)
	require.Empty(t, res.Data.Errors)

	res, err = client.changeUserSegments(1, []string{add[0]}, []string{})
	require.ErrorIs(t, err, ErrBadRequest)
	require.False(t, res.Data.Done)
	require.Len(t, res.Data.Errors, 1)
	require.Equal(t, res.Data.Errors[add[0]], repo.ErrRelationExists.Error())

	resGet, err = client.getUserSegments(1)
	require.NoError(t, err)
	require.Len(t, resGet.Data, len(add))

	res, err = client.changeUserSegments(1, []string{}, add[1:])
	require.NoError(t, err)
	require.True(t, res.Data.Done)
	require.Empty(t, res.Data.Errors)

	resGet, err = client.getUserSegments(1)
	require.NoError(t, err)
	require.Len(t, resGet.Data, 1)
	require.Equal(t, resGet.Data[0].Slug, add[0])
}
