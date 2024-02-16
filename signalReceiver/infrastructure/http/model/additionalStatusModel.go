package model

type AdditionalStatusModel struct {
	StatusCode int32 `json:"StatusCode"`

	ServerStatusCode string `json:"ServerStatusCode,omitempty"`

	Severity string `json:"Severity"`

	StatusDesc string `json:"StatusDesc,omitempty"`
}
