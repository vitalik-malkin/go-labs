package models

import (
	"context"
	"encoding/json"
)

type ParamVal struct {
	Value   json.RawMessage `json:"value,omitempty"`
	Decoder string          `json:"decoder,omitempty"`
}

func (v *ParamVal) String(ctx context.Context) (string, error) {
	panic("not implemented yet")
}
