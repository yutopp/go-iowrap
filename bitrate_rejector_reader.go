//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package iowrap

import (
	"github.com/pkg/errors"
	"io"
	"time"
)

type BitrateRejectorReader struct {
	reader         io.Reader
	maxBitrateKbps uint32

	bitrateKbps float64
	readSize    uint64
	now         func() time.Time // for mock
	last        time.Time
}

func NewBitrateRejectorReader(r io.Reader, maxBitrateKbps uint32) *BitrateRejectorReader {
	return &BitrateRejectorReader{
		reader:         r,
		maxBitrateKbps: maxBitrateKbps,

		now: time.Now,
	}
}

func (r *BitrateRejectorReader) Read(b []byte) (int, error) {
	// Check bitrate first
	cur := r.now()
	if r.last.IsZero() {
		r.last = cur
	}
	diff := cur.Sub(r.last)
	if diff >= 1*time.Second {
		r.bitrateKbps = (float64(r.readSize) / float64(diff/time.Second)) * 8 / 1024.0
		// reset
		r.readSize = 0
		r.last = cur

		if r.bitrateKbps > float64(r.maxBitrateKbps) {
			return 0, errors.Errorf(
				"Bitrate exceeded: Limit = %vkbps, Value = %vkbps",
				r.maxBitrateKbps,
				r.bitrateKbps,
			)
		}
	}

	n, err := r.reader.Read(b)
	if err != nil {
		return 0, err
	}

	r.readSize += uint64(n)

	return n, nil
}

func (r *BitrateRejectorReader) Close() error {
	if c, ok := r.reader.(io.Closer); ok {
		return c.Close()
	}

	return nil
}

func (r *BitrateRejectorReader) BitrateKbps() float64 {
	return r.bitrateKbps
}
