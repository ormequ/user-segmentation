package http

import (
	"net/http"
	"user-segmentation/internal/entities/segments"
	"user-segmentation/internal/service"
)

type segment struct {
	Slug string `json:"slug" binding:"required"`
}

type CreateSegmentRequest segment

type DeleteSegmentRequest segment

type SegmentProcessedResponse struct {
	Done bool `json:"done"`
}

func errToSegmentProcessed(err error) SegmentProcessedResponse {
	code, _ := hideError(err)
	return SegmentProcessedResponse{
		Done: code == http.StatusOK,
	}
}

type ChangeUserSegmentsRequest struct {
	Remove []string `json:"remove"`
	Add    []string `json:"add"`
}

type ChangeResultResponse struct {
	Done   bool              `json:"done"`
	Errors map[string]string `json:"errors"`
}

func changeResultToResponse(res service.ChangeErrors, err error) ChangeResultResponse {
	if err != nil {
		return ChangeResultResponse{}
	}
	if len(res) == 0 {
		res = nil
	}
	return ChangeResultResponse{
		Done:   len(res) == 0,
		Errors: res,
	}
}

type SegmentResponse segment

func segmentsToResponse(seg []segments.Segment) []SegmentResponse {
	res := make([]SegmentResponse, len(seg))
	for i := range res {
		res[i].Slug = seg[i].Slug
	}
	return res
}
