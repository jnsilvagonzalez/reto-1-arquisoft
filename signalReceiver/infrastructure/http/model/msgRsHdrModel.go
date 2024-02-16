package model

type MsgRsHdrModel struct {
	RqUID string `json:"RqUID,omitempty"`

	Status *StatusModel `json:"Status,omitempty"`
}
