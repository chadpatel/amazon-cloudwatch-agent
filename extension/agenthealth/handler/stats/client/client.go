// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/amazon-contributing/opentelemetry-collector-contrib/extension/awsmiddleware"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jellydator/ttlcache/v3"

	"github.com/aws/amazon-cloudwatch-agent/extension/agenthealth/handler/stats/agent"
)

const (
	handlerID   = "cloudwatchagent.ClientStats"
	ttlDuration = 10 * time.Second
	cacheSize   = 1000
)

var (
	rejectedEntityInfo = []byte("\"rejectedEntityInfo\"")
)

type Stats interface {
	awsmiddleware.RequestHandler
	awsmiddleware.ResponseHandler
	agent.StatsProvider
}

type requestRecorder struct {
	start        time.Time
	payloadBytes int64
}

type clientStatsHandler struct {
	filter           agent.OperationsFilter
	getOperationName func(ctx context.Context) string
	getRequestID     func(ctx context.Context) string

	statsByOperation sync.Map
	requestCache     *ttlcache.Cache[string, *requestRecorder]
}

var _ Stats = (*clientStatsHandler)(nil)

func NewHandler(filter agent.OperationsFilter) Stats {
	requestCache := ttlcache.New[string, *requestRecorder](
		ttlcache.WithTTL[string, *requestRecorder](ttlDuration),
		ttlcache.WithCapacity[string, *requestRecorder](cacheSize),
		ttlcache.WithDisableTouchOnHit[string, *requestRecorder](),
	)
	go requestCache.Start()
	return &clientStatsHandler{
		filter:           filter,
		getOperationName: awsmiddleware.GetOperationName,
		getRequestID:     awsmiddleware.GetRequestID,
		requestCache:     requestCache,
	}
}

func (csh *clientStatsHandler) ID() string {
	return handlerID
}

func (csh *clientStatsHandler) Position() awsmiddleware.HandlerPosition {
	return awsmiddleware.After
}

func (csh *clientStatsHandler) HandleRequest(ctx context.Context, r *http.Request) {
	operation := csh.getOperationName(ctx)
	if !csh.filter.IsAllowed(operation) {
		return
	}
	requestID := csh.getRequestID(ctx)
	recorder := &requestRecorder{start: time.Now()}

	// Check if ContentLength is already provided
	if r.ContentLength > 0 {
		recorder.payloadBytes = r.ContentLength
	} else if r.Body != nil {
		// Create a ReadSeeker from the body if it's not already one
		rsc, err := toReadSeeker(r.Body)
		if err != nil {
			// Handle error if the body cannot be read
			return
		}

		// Try to determine the length of the request body
		if length, err := getSeekerLength(rsc); err == nil && length > 0 {
			recorder.payloadBytes = length
		} else if body, err := r.GetBody(); err == nil {
			// If the length cannot be determined, use GetBody() to copy the body to io.Discard
			recorder.payloadBytes, _ = io.Copy(io.Discard, body)
		}
	}

	// Store the recorder in cache with the request ID as the key
	csh.requestCache.Set(requestID, recorder, ttlcache.DefaultTTL)
}

// Helper function to convert io.ReadCloser to io.ReadSeeker by reading it into memory
func toReadSeeker(rc io.ReadCloser) (io.ReadSeeker, error) {
	defer rc.Close()
	bodyBytes, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(bodyBytes), nil
}

// Helper function to get the length of an io.ReadSeeker
func getSeekerLength(seeker io.ReadSeeker) (int64, error) {
	// Save the current position
	currentPos, err := seeker.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	// Seek to the end to find the length
	length, err := seeker.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	// Return to the original position
	_, err = seeker.Seek(currentPos, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return length, nil
}

func (csh *clientStatsHandler) HandleResponse(ctx context.Context, r *http.Response) {
	operation := csh.getOperationName(ctx)
	if !csh.filter.IsAllowed(operation) {
		return
	}
	requestID := csh.getRequestID(ctx)
	item, ok := csh.requestCache.GetAndDelete(requestID)
	if !ok {
		return
	}
	recorder := item.Value()
	stats := agent.Stats{
		PayloadBytes: aws.Int(int(recorder.payloadBytes)),
		StatusCode:   aws.Int(r.StatusCode),
	}
	latency := time.Since(recorder.start)
	stats.LatencyMillis = aws.Int64(latency.Milliseconds())
	if rejectedEntityInfoExists(r) {
		stats.EntityRejected = aws.Int(1)
	}
	csh.statsByOperation.Store(operation, stats)
}

func (csh *clientStatsHandler) Stats(operation string) agent.Stats {
	value, ok := csh.statsByOperation.Load(operation)
	if !ok {
		return agent.Stats{}
	}
	stats, ok := value.(agent.Stats)
	if !ok {
		return agent.Stats{}
	}
	return stats
}

// rejectedEntityInfoExists checks if the response body
// contains element rejectedEntityInfo
func rejectedEntityInfoExists(r *http.Response) bool {
	// Example body for rejectedEntityInfo would be:
	// {"rejectedEntityInfo":{"errorType":"InvalidAttributes"}}
	if r == nil || r.Body == nil {
		return false
	}
	defer r.Body.Close()
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return false
	}
	return bytes.Contains(bodyBytes, rejectedEntityInfo)
}
