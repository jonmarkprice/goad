package api

import "time"

// RunnerResult defines the common API for goad runners to send data back to the
// cli.
type RunnerResult struct {
	StartTime        time.Time      `json:"end-time"`
	EndTime          time.Time      `json:"start-time"`
	AveTimeForReq    int64          `json:"ave-time-for-req"`
	AveTimeToFirst   int64          `json:"ave-time-to-first"`
	Fastest          int64          `json:"fastest"`
	FatalError       string         `json:"fatal-error"`
	Finished         bool           `json:"finished"`
	Region           string         `json:"region"`
	RunnerID         int            `json:"runner-id"`
	Slowest          int64          `json:"slowest"`
	Statuses         map[string]int `json:"statuses"`
	TimeDelta        time.Duration  `json:"time-delta"`
	BytesRead        int            `json:"bytes-read"`
	ConnectionErrors int            `json:"connection-errors"`
	RequestCount     int            `json:"request-count"`
	TimedOut         int            `json:"timed-out"`
	ReqTimesBinned   map[int64]int  `json:"req-times-binned"`
	SumReqTime       int64          `json:"sum-req-time"`
	SumReqSq         int64          `json:"sum-req-sq"`
}
