package model

import (
	"time"
)

type StatusModel struct {
	StatusCode int `json:"StatusCode"`

	ServerStatusCode int `json:"ServerStatusCode,omitempty"`

	Severity string `json:"Severity,omitempty"`

	StatusDesc string `json:"StatusDesc,omitempty"`

	AdditionalStatus *AdditionalStatusModel `json:"AdditionalStatus,omitempty"`

	EndDt time.Time `json:"EndDt,omitempty"`
}
