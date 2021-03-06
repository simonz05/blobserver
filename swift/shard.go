// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package swift

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
)

const shardCount = 8 << 7

type sharder [shardCount]string

func (s sharder) shard(v string) string {
	src := md5.Sum([]byte(v))
	str := hex.EncodeToString(src[:])
	return s[int(s.num(str))]
}

func (s sharder) num(digest string) uint32 {
	return s.sum32(digest) % uint32(shardCount)
}

func (s sharder) sum32(digest string) uint32 {
	vv, _ := strconv.ParseUint(digest[:4], 16, 32)
	return uint32(vv)
}

func newSharder() (s sharder) {
	for i := range s {
		s[i] = fmt.Sprintf("%0.2x", i)
	}
	return s
}
