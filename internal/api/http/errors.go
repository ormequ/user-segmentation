package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"user-segmentation/internal/entities/segments"
	"user-segmentation/internal/repo"
)

var (
	ErrInternal = errors.New("internal error")
	ErrChanging = errors.New("changing error")
)

func hideError(err error) (int, error) {
	if err == nil || errors.Is(err, repo.ErrNoSegments) {
		return http.StatusOK, nil
	}
	if errors.Is(err, repo.ErrSegmentAlreadyExists) {
		return http.StatusConflict, err
	}
	if errors.Is(err, repo.ErrRelationNotFound) || errors.Is(err, repo.ErrSegmentNotFound) {
		return http.StatusNotFound, err
	}
	if errors.Is(err, segments.ErrEmptySlug) || errors.Is(err, segments.ErrSlugToLong) || errors.Is(err, ErrChanging) {
		return http.StatusBadRequest, err
	}
	return http.StatusInternalServerError, ErrInternal
}

func errorResponse(err error) gin.H {
	return gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}

func handleError(c *gin.Context, err error, data any) {
	code, hidden := hideError(err)
	var resErr any
	if errors.Is(hidden, nil) {
		resErr = nil
	} else {
		resErr = hidden.Error()
	}
	c.JSON(code, gin.H{
		"data":  data,
		"error": resErr,
	})
}
