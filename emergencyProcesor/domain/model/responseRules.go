package model

type ResponseRules struct {
	VehiculoID string
	ReglaId    string
	Actions    []Action
}

type Action struct {
	ActionId  string `json:"actionId"`
	Parameter string `json:"parameter"`
	Value     string `json:"value"`
}
