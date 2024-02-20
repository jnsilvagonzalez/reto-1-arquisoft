package model

type ReqEmergencySignal struct {
	RqUID string `json:"RqUID"`

	IdVehicle string `json:"IdVehicle"`

	Speed int32 `json:"StatusCode"`

	Address string `json:"Address"`

	Latitude string `json:"Severity"`

	Longitude string `json:"Longitude,omitempty"`
}
