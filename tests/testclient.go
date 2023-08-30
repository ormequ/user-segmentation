package tests

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	httpserver "user-segmentation/internal/api/http"
	"user-segmentation/internal/repo/segments"
	"user-segmentation/internal/service"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
	ErrConflict   = errors.New("conflict")
)

var db *pgx.Conn

func setupClient() *testClient {
	a := service.New(
		segments.New(db),
	)
	srv := httpserver.New(slog.Default(), ":8888", gin.ReleaseMode, a)
	testSrv := httptest.NewServer(srv.Handler)

	return &testClient{
		client:  testSrv.Client(),
		baseURL: testSrv.URL,
	}
}

type testClient struct {
	client  *http.Client
	baseURL string
}

func (tc *testClient) request(body map[string]any, method string, endpoint string, out any) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(method, tc.baseURL+"/api/"+endpoint, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := tc.client.Do(req)
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}
	var code error
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			code = ErrNotFound
		} else if resp.StatusCode == http.StatusBadRequest {
			code = ErrBadRequest
		} else if resp.StatusCode == http.StatusConflict {
			code = ErrConflict
		} else {
			return fmt.Errorf("unexpected status code: %s", resp.Status)
		}
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read response: %w", err)
	}
	err = json.Unmarshal(respBody, out)
	if err != nil {
		return fmt.Errorf("unable to unmarshal: %w", err)
	}

	return code
}

type segmentProcessed httpserver.SegmentProcessedResponse
type segmentProcessedResponse struct {
	Data  segmentProcessed `json:"data"`
	Error string           `json:"error"`
}

func (tc *testClient) createSegment(slug string) (segmentProcessedResponse, error) {
	body := map[string]any{
		"slug": slug,
	}
	var response segmentProcessedResponse
	err := tc.request(body, http.MethodPost, "segments", &response)
	return response, err
}

func (tc *testClient) deleteSegment(slug string) (segmentProcessedResponse, error) {
	body := map[string]any{
		"slug": slug,
	}
	var response segmentProcessedResponse
	err := tc.request(body, http.MethodDelete, "segments", &response)
	return response, err
}

type changeResult httpserver.ChangeResultResponse
type changeResultResponse struct {
	Data  changeResult `json:"data"`
	Error string       `json:"error"`
}

func (tc *testClient) changeUserSegments(userID int64, add []string, remove []string) (changeResultResponse, error) {
	body := map[string]any{
		"add":    add,
		"remove": remove,
	}
	var response changeResultResponse
	err := tc.request(body, http.MethodPost, fmt.Sprintf("users/%d", userID), &response)
	return response, err
}

type segment httpserver.SegmentResponse
type segmentsResponse struct {
	Data  []segment `json:"data"`
	Error string    `json:"error"`
}

func (tc *testClient) getUserSegments(userID int64) (segmentsResponse, error) {
	body := map[string]any{}
	var response segmentsResponse
	err := tc.request(body, http.MethodGet, fmt.Sprintf("users/%d", userID), &response)
	return response, err
}
