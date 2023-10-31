package dto

import "encoding/json"

type ResultDTO struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data *json.RawMessage `json:"data"`
}
