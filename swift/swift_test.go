// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package swift

import (
	"fmt"
	"os"
	"testing"

	"github.com/simonz05/blobserver"
	"github.com/simonz05/blobserver/config"
	"github.com/simonz05/blobserver/storagetest"
)

func storageFromConf(t *testing.T) blobserver.Storage {
	configFile := os.Getenv("BLOBSERVER_SWIFT_TEST_CONFIG")
	if configFile == "" {
		t.Skip("Skipping manual test. To enable, set the environment variable BLOBSERVER_SWIFT_TEST_CONFIG to the path of a JSON configuration for the s3 storage type.")
	}
	conf, err := config.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Error reading swift configuration file %s: %v", configFile, err)
	}

	sto, err := newFromConfig(conf)
	if err != nil {
		t.Fatalf("newFromConfig error: %v", err)
	}

	return sto
}

func TestSwift(t *testing.T) {
	storagetest.Test(t, func(t *testing.T) (sto blobserver.Storage, cleanup func()) {
		return storageFromConf(t), func() {}
	})
}

func creator(ech chan error, in, out chan string, sto *swiftStorage) {
	for cont := range in {
		err := sto.createContainer(cont)
		if err != nil {
			ech <- err
			return
		}

		out <- cont
	}
}

func statter(ech chan error, in, out chan string, sto *swiftStorage) {
	for cont := range in {
		_, headers, err := sto.conn.Container(cont)
		if err != nil {
			ech <- err
			return
		}
		r := headers["X-Container-Read"]
		exp := ".r:*,.rlistings"

		if r != exp {
			ech <- fmt.Errorf("exp %s, got %s", exp, r)
			return
		}
		out <- cont
	}
}

func deleter(ech chan error, in, out chan string, sto *swiftStorage) {
	for cont := range in {
		err := sto.conn.ContainerDelete(cont)
		if err != nil {
			ech <- err
			return
		}
		out <- cont
	}
}

func TestSwiftContainerACL(t *testing.T) {
	sto := storageFromConf(t)
	sw := sto.(*swiftStorage)
	w := 64
	if testing.Short() {
		w = 4
	}
	cch := make(chan string)
	sch := make(chan string)
	dch := make(chan string)
	ech := make(chan error)
	qch := make(chan string, 32)

	for i := 0; i < w; i++ {
		go creator(ech, cch, sch, sw)
		go statter(ech, sch, dch, sw)
		go deleter(ech, dch, qch, sw)
	}

	var shards []string
	ss := newSharder()

	if testing.Short() {
		shards = ss[:5]
	} else {
		shards = ss[:100]
	}

	for _, shard := range shards {
		go func(ch chan string, shard string) {
			ch <- "ctest-" + shard
		}(cch, shard)
	}

	for i := len(shards); i > 0; i-- {
		select {
		case <-qch:
		case err := <-ech:
			t.Fatal(err)
		}
	}

	close(cch)
	close(sch)
	close(dch)
	close(qch)
}
