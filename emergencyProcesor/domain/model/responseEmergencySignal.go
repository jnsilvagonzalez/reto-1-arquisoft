package model

type ResEmergencySignal struct {
	RqUID   string
	Message string
	Err     error `json:"Err,omitempty"`
}
