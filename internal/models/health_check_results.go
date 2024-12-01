package models

import (
	"time"
)

type HealthCheckResult struct {
	ID             string    `json:"id"`
	TargetEndpoint string    `json:"target_endpoint"`
	HTTPMethod     string    `json:"http_method"`
	CalloutTime    time.Time `json:"callout_timestamp"`
	ResponseCode   int       `json:"response_code"`
	Response       string    `json:"response"`
	TimeTaken      float64   `json:"time_taken"`
}