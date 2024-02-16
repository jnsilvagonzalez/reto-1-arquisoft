package model

type ErrorCodeModel struct {
	HTTPCode int `json:"HttpCode,omitempty"`

	HTTPMessage string `json:"HttpMessage,omitempty"`

	MoreInformation string `json:"MoreInformation,omitempty"`
}
