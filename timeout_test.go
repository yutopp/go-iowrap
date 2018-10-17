//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package iowrap

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	rwc := &TimeoutConn{
		now: func() time.Time {
			return time.Time{}
		},
	}

	assert.Equal(t, time.Time{}, rwc.calcDeadline(0))
	assert.Equal(t, time.Time{}.Add(1*time.Second), rwc.calcDeadline(1*time.Second))
}

func TestTimeoutRead(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	rwc := NewTimeoutConn(server, 500*time.Millisecond, 0)

	ch := make(chan struct{})
	go func() {
		buf := make([]byte, 1024)
		_, err := rwc.Read(buf)
		assert.NotNil(t, err)
		assert.Equal(t, true, err.(*net.OpError).Timeout())
		ch <- struct{}{}
	}()

	<-ch
}

func TestTimeoutWrite(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	rwc := NewTimeoutConn(server, 0, 500*time.Millisecond)

	ch := make(chan struct{})
	go func() {
		buf := make([]byte, 1024)
		_, err := rwc.Write(buf)
		assert.NotNil(t, err)
		assert.Equal(t, true, err.(*net.OpError).Timeout())
		ch <- struct{}{}
	}()

	<-ch
}
