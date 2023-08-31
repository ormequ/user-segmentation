package http

import (
	"encoding/csv"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"user-segmentation/internal/logger"
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
		result, err := svc.ChangeUserSegments(c, int64(id), req.Add, req.Remove)
		response := changeResultToResponse(result, err)
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

func getHistory(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		year, err := strconv.Atoi(c.Param("year"))
		var month int
		if err == nil {
			month, err = strconv.Atoi(c.Param("month"))
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(ErrInvalidRequest))
			return
		}
		ops, err := svc.GetHistory(c, year, month)
		if err != nil {
			code, err := hideError(err)
			c.JSON(code, errorResponse(err))
		}
		c.Writer.Header().Set("Content-Type", "text/csv")
		c.Writer.Header().Set("Content-Disposition", "attachment;filename=history.csv")
		c.Writer.WriteHeader(http.StatusOK)
		csvWriter := csv.NewWriter(c.Writer)
		defer csvWriter.Flush()
		for _, line := range ops {
			err := csvWriter.Write(line)
			if err != nil {
				logger.InternalErr(c, err, "api.http.getHistory")
			}
		}
	}
}
