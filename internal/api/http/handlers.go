package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"user-segmentation/internal/service"
)

var ErrInvalidRequest = errors.New("invalid request")

func createSegment(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateSegmentRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(ErrInvalidRequest))
			return
		}
		err := svc.CreateSegment(c, req.Slug)
		handleError(c, err, errToSegmentProcessed(err))
	}
}

func deleteSegment(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req DeleteSegmentRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(ErrInvalidRequest))
			return
		}
		err := svc.DeleteSegment(c, req.Slug)
		handleError(c, err, errToSegmentProcessed(err))
	}
}

func changeUserSegments(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ChangeUserSegmentsRequest
		id, err := strconv.Atoi(c.Param("user_id"))
		if err == nil {
			err = c.BindJSON(&req)
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(ErrInvalidRequest))
			return
		}
		result := svc.ChangeUserSegments(c, int64(id), req.Add, req.Remove)
		response := changeResultToResponse(result)
		if !response.Done {
			err = ErrChanging
		}
		handleError(c, err, response)
	}
}

func getUserSegments(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(ErrInvalidRequest))
			return
		}
		seg, err := svc.GetUserSegments(c, int64(id))
		handleError(c, err, segmentsToResponse(seg))
	}
}
