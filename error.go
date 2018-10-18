//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package iowrap

import (
	"fmt"
)

type BitrateExceededError struct {
	MaxKbps     float64
	CurrentKbps float64
}

func (err *BitrateExceededError) Error() string {
	return fmt.Sprintf(
		"Bitrate exceeded: Limit = %vkbps, Value = %vkbps",
		err.MaxKbps,
		err.CurrentKbps,
	)
}
