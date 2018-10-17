//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package iowrap

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBitrateRejectorReaderRejected(t *testing.T) {
	br := bytes.NewReader(make([]byte, 4096))
	maxBitrateKbps := uint32(8) // 8Kbps

	r := NewBitrateRejectorReader(br, maxBitrateKbps)
	r.now = func() time.Time {
		return time.Unix(0, 0)
	}

	buf := make([]byte, 1024*2)

	// Read 2KB=16Kbits
	n, err := r.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, 2048, n)
	assert.InDelta(t, 0.0, r.BitrateKbps(), 0.01) // Not calculated yet

	// simulate 1 sec
	r.now = func() time.Time {
		return time.Unix(1, 0)
	}

	// First read after 1s
	// Read 8Kbits
	n, err = r.Read(buf[:1024])
	assert.EqualError(t, err, "Bitrate exceeded: Limit = 8kbps, Value = 16kbps")
	assert.InDelta(t, 16.0, r.BitrateKbps(), 0.01)
}

func TestBitrateRejectorReaderAccepted(t *testing.T) {
	br := bytes.NewReader(make([]byte, 4096))
	maxBitrateKbps := uint32(8) // 8Kbps

	r := NewBitrateRejectorReader(br, maxBitrateKbps)
	r.now = func() time.Time {
		return time.Unix(0, 0)
	}

	buf := make([]byte, 512)
	for i := 0; i < 4096/512; i++ {
		// Read 512Bytes=4Kbits
		n, err := r.Read(buf)
		assert.Nil(t, err)
		assert.Equal(t, 512, n)
		if i == 0 {
			assert.InDelta(t, 0.0, r.BitrateKbps(), 0.01)
		} else {
			assert.InDelta(t, 4.0, r.BitrateKbps(), 0.01)
		}

		// simulate 1 sec per loop
		r.now = func() time.Time {
			return time.Unix(int64(i), 0)
		}
	}
}
