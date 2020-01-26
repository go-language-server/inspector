// Copyright 2020 The go-language-server Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package inspector

import (
	"strconv"
	"time"
)

// LogFormat represents a Language Server Protocol Inspector log format.
type LogFormat uint8

const (
	// TextFormat represents a text format of inspector log.
	TextFormat LogFormat = 1 + iota

	// JSONFormat represents a JSON format of inspector log.
	JSONFormat
)

// String implements fmt.Stringer.
func (lf LogFormat) String() string {
	switch lf {
	case TextFormat:
		return "text"
	case JSONFormat:
		return "json"
	default:
		return strconv.FormatUint(uint64(lf), 10)
	}
}

// Payload represents a Language Server Protocol Inspector specification payload.
type Payload struct {
	Time       time.Time     `json:"time"`
	Msg        string        `json:"msg"`
	MsgKind    MessageKind   `json:"msgKind"`
	MsgType    string        `json:"msgType"`
	MsgID      string        `json:"msgId,omitempty"`
	MsgLatency time.Duration `json:"msgLatency,omitempty"`
	Arg        []interface{} `json:"arg"`
}

// MessageKind represents a message kind.
type MessageKind uint8

const (
	// SendNotification represents a send notification message kind.
	SendNotification MessageKind = 1 + iota

	// ReceiveNotification represents a receive notification message kind.
	ReceiveNotification

	// SendRequest represents a send request message kind.
	SendRequest

	// ReceiveRequest represents a receive request message kind.
	ReceiveRequest

	// SendResponse represents a send response message kind.
	SendResponse

	// ReceiveResponse represents a receive response message kind.
	ReceiveResponse
)

// String implements fmt.Stringer.
func (ms MessageKind) String() string {
	switch ms {
	case SendNotification:
		return "send-notification"
	case ReceiveNotification:
		return "recv-notification"
	case SendRequest:
		return "send-request"
	case ReceiveRequest:
		return "recv-request"
	case SendResponse:
		return "send-response"
	case ReceiveResponse:
		return "recv-response"
	default:
		return strconv.FormatUint(uint64(ms), 10)
	}
}
