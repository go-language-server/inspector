// Copyright 2020 The go-language-server Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package inspector

import "go.uber.org/zap/zapcore"

// MarshalLogObject implements zapcore.ObjectMarshaler.
func (p *Payload) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddTime("time", p.Time)
	enc.AddString("msg", p.Msg)
	enc.AddString("msgKind", p.MsgKind.String())
	enc.AddString("msgType", p.MsgType)
	enc.AddString("msgId", p.MsgID)
	enc.AddDuration("msgLatency", p.MsgLatency)
	return enc.AddReflected("arg", p.Arg)
}
