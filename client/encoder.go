// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package client

import (
	"encoding/json"
	"io"
)

type ResourceEncoder interface {
	Encode(res Resources) error
}

type JSONEncoder struct {
	w *json.Encoder
}

func NewJSONEncoder(w io.Writer) *JSONEncoder {
	return &JSONEncoder{w: json.NewEncoder(w)}
}

func (e *JSONEncoder) Encode(res Resources) error {
	out := make(map[string]string, len(res))

	for _, v := range res {
		out[v.Path] = v.URL
	}

	return e.w.Encode(out)
}
