package model

type ResPostReceiveSignal struct {
	RqUID   string
	Message string
	Err     error `json:"Err,omitempty"`
}
