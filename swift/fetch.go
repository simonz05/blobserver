// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package swift

import (
	"fmt"
	"io"
	"strconv"

	"github.com/simonz05/blobserver/blob"
	"github.com/simonz05/util/log"
)

// getInt64FromHeader is a helper function to decode int64 from header.
func getInt64FromHeader(headers map[string]string, header string) (result int64, err error) {
	value := headers[header]
	result, err = strconv.ParseInt(value, 10, 64)
	if err != nil {
		err = fmt.Errorf("Bad Header '%s': '%s': %s", header, value, err)
	}
	return
}

func (sto *swiftStorage) Fetch(br blob.Ref) (file io.ReadCloser, size uint32, err error) {
	ref, cont := sto.refContainer(br)
	log.Println("Fetch: ", ref, cont)
	f, h, err := sto.conn.ObjectOpen(cont, ref, true, nil)
	if err != nil {
		return
	}
	n, err := getInt64FromHeader(h, "Content-Length")
	if err != nil {
		return
	}
	return f, uint32(n), err
}
