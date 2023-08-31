package test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
	"user-segmentation/internal/entities/operations"
	"user-segmentation/internal/entities/segments"
	"user-segmentation/internal/service"
	"user-segmentation/internal/service/mocks"
)

const longSlug = "111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111"

func storeRepo(t *testing.T) service.SegmentsRepo {
	r := mocks.NewSegmentsRepo(t)
	r.
		On("Store", mock.Anything, mock.AnythingOfType("segments.Segment")).
		Return(nil)
	return r
}

func deleteRepo(t *testing.T) service.SegmentsRepo {
	r := mocks.NewSegmentsRepo(t)
	r.
		On("Delete", mock.Anything, mock.AnythingOfType("segments.Segment")).
		Return(nil)
	return r
}

func getHistoryRepo(t *testing.T, rows int64, opType operations.Type) service.HistoryRepo {
	res := make([]operations.Operation, rows)
	for i := range res {
		res[i].UserID = int64(i)
		res[i].Segment = segments.Segment{Slug: fmt.Sprintf("slug-%d", i)}
		res[i].Time = time.Now().UTC()
		res[i].Type = opType
	}
	r := mocks.NewHistoryRepo(t)
	r.
		On("Get", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int")).
		Return(res, nil)
	return r
}

func TestService_CreateSegment(t *testing.T) {
	type fields struct {
		segments service.SegmentsRepo
		history  service.HistoryRepo
	}
	type args struct {
		ctx  context.Context
		slug string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "correct adding",
			fields: fields{
				segments: storeRepo(t),
				history:  nil,
			},
			args: args{
				ctx:  context.Background(),
				slug: "slug",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "empty slug",
			fields: fields{
				segments: nil,
				history:  nil,
			},
			args: args{
				ctx:  context.Background(),
				slug: "",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, segments.ErrEmptySlug)
			},
		},
		{
			name: "slug too long",
			fields: fields{
				segments: nil,
				history:  nil,
			},
			args: args{
				ctx:  context.Background(),
				slug: longSlug,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, segments.ErrSlugToLong)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.Service{
				Segments: tt.fields.segments,
				History:  tt.fields.history,
			}
			tt.wantErr(t, s.CreateSegment(tt.args.ctx, tt.args.slug), fmt.Sprintf("CreateSegment(%v, %v)", tt.args.ctx, tt.args.slug))
		})
	}
}

func TestService_DeleteSegment(t *testing.T) {
	type fields struct {
		segments service.SegmentsRepo
		history  service.HistoryRepo
	}
	type args struct {
		ctx  context.Context
		slug string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "correct deleting",
			fields: fields{
				segments: deleteRepo(t),
				history:  nil,
			},
			args: args{
				ctx:  context.Background(),
				slug: "slug",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "empty slug",
			fields: fields{
				segments: nil,
				history:  nil,
			},
			args: args{
				ctx:  context.Background(),
				slug: "",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, segments.ErrEmptySlug)
			},
		},
		{
			name: "slug too long",
			fields: fields{
				segments: nil,
				history:  nil,
			},
			args: args{
				ctx:  context.Background(),
				slug: "111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, segments.ErrSlugToLong)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.Service{
				Segments: tt.fields.segments,
				History:  tt.fields.history,
			}
			tt.wantErr(t, s.DeleteSegment(tt.args.ctx, tt.args.slug), fmt.Sprintf("DeleteSegment(%v, %v)", tt.args.ctx, tt.args.slug))
		})
	}
}

func TestService_GetHistory(t *testing.T) {
	type fields struct {
		segments service.SegmentsRepo
		history  service.HistoryRepo
	}
	type args struct {
		ctx   context.Context
		year  int
		month int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    [][]string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "basic getting add",
			fields: fields{
				segments: nil,
				history:  getHistoryRepo(t, 2, operations.Add),
			},
			args: args{
				ctx:   context.Background(),
				year:  2023,
				month: 9,
			},
			want: [][]string{
				{"User ID", "Segment", "Operation", "Timestamp UTC"},
				{"0", "slug-0", "add", time.Now().UTC().String()},
				{"1", "slug-1", "add", time.Now().UTC().String()},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "basic getting remove",
			fields: fields{
				segments: nil,
				history:  getHistoryRepo(t, 2, operations.Remove),
			},
			args: args{
				ctx:   context.Background(),
				year:  2023,
				month: 9,
			},
			want: [][]string{
				{"User ID", "Segment", "Operation", "Timestamp UTC"},
				{"0", "slug-0", "remove", time.Now().UTC().String()},
				{"1", "slug-1", "remove", time.Now().UTC().String()},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "empty result",
			fields: fields{
				segments: nil,
				history:  getHistoryRepo(t, 0, operations.Remove),
			},
			args: args{
				ctx:   context.Background(),
				year:  2023,
				month: 9,
			},
			want: [][]string{
				{"User ID", "Segment", "Operation", "Timestamp UTC"},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "invalid year",
			fields: fields{
				segments: nil,
				history:  nil,
			},
			args: args{
				ctx:   context.Background(),
				year:  0,
				month: 9,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, service.ErrInvalidDates)
			},
		},
		{
			name: "invalid month",
			fields: fields{
				segments: nil,
				history:  nil,
			},
			args: args{
				ctx:   context.Background(),
				year:  2023,
				month: 0,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, service.ErrInvalidDates)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.Service{
				Segments: tt.fields.segments,
				History:  tt.fields.history,
			}
			got, err := s.GetHistory(tt.args.ctx, tt.args.year, tt.args.month)
			if !tt.wantErr(t, err, fmt.Sprintf("GetHistory(%v, %v, %v)", tt.args.ctx, tt.args.year, tt.args.month)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetHistory(%v, %v, %v)", tt.args.ctx, tt.args.year, tt.args.month)
		})
	}
}

func putHistoryRepo(t *testing.T) service.HistoryRepo {
	r := mocks.NewHistoryRepo(t)
	r.On("Put", mock.Anything, mock.Anything).Return(nil)
	return r
}

func changeUserSegmentsRepo(t *testing.T) service.SegmentsRepo {
	r := mocks.NewSegmentsRepo(t)
	r.
		On("ChangeUserSegments", mock.Anything, mock.AnythingOfType("int64"), mock.Anything, mock.Anything).
		Return(service.ChangeErrors{})
	return r
}

func TestService_ChangeUserSegments(t *testing.T) {
	type fields struct {
		segments service.SegmentsRepo
		history  service.HistoryRepo
	}
	type args struct {
		ctx    context.Context
		userID int64
		add    []string
		remove []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    service.ChangeErrors
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "correct changing",
			fields: fields{
				segments: changeUserSegmentsRepo(t),
				history:  putHistoryRepo(t),
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
				add:    []string{"slug-1", "slug-2"},
				remove: []string{"slug-3", "slug-4"},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "correct changing with empty args",
			fields: fields{
				segments: changeUserSegmentsRepo(t),
				history:  putHistoryRepo(t),
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
				add:    []string{},
				remove: []string{},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "incorrect segments",
			fields: fields{
				segments: nil,
				history:  nil,
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
				add:    []string{"slug-1", ""},
				remove: []string{longSlug, "slug-4"},
			},
			want: service.ChangeErrors{
				"":       segments.ErrEmptySlug.Error(),
				longSlug: segments.ErrSlugToLong.Error(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.Service{
				Segments: tt.fields.segments,
				History:  tt.fields.history,
			}
			got, err := s.ChangeUserSegments(tt.args.ctx, tt.args.userID, tt.args.add, tt.args.remove)
			if !tt.wantErr(t, err, fmt.Sprintf("ChangeUserSegments(%v, %v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.add, tt.args.remove)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ChangeUserSegments(%v, %v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.add, tt.args.remove)
		})
	}
}
