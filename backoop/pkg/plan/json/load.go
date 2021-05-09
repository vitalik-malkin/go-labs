package json

import (
	"context"
	go_json "encoding/json"
	"fmt"
	"io"
)

func Load(ctx context.Context, r io.Reader) (Plan, error) {
	decoder := go_json.NewDecoder(r)
	p := &plan{}
	err := decoder.Decode(p)
	if err != nil {
		return nil, fmt.Errorf("error while decoding json; err: %w", err)
	}

	return p, nil
}
