// Copyright 2020 The go-language-server Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package inspector

import (
	"fmt"
	"time"
	"unsafe"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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

// Logger represents a Language Server Protocol Inspector logger.
type Logger struct {
	*zap.Logger
}

// NewZapLogger returns a new zap.Logger which implemented
// Language Server Protocol Inspector specification log format.
func NewZapLogger(opts ...zap.Option) *Logger {
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "",
		NameKey:        "",
		CallerKey:      "",
		MessageKey:     "",
		StacktraceKey:  "",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    nil,
		EncodeTime:     zapcore.TimeEncoder(ISO8601TimeEncoder),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   nil,
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encCfg,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build(opts...)
	if err != nil {
		panic(fmt.Errorf("inspector.NewZapLogger: %v", err))
	}

	return &Logger{logger}
}

// ISO8601TimeEncoder serializes a time.Time to an ISO8601-formatted string
// with millisecond precision.
// Same as t.Format("2006-01-02T15:04:05.000Z").
//
// It optimized at byte slice level, faster than time.Format.
func ISO8601TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	t = t.UTC()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	msec := t.Nanosecond() / 1e6

	b := make([]byte, 24)

	b[0] = byte((year/1000)%10) + '0'
	b[1] = byte((year/100)%10) + '0'
	b[2] = byte((year/10)%10) + '0'
	b[3] = byte(year%10) + '0'
	b[4] = '-'
	b[5] = byte((month)/10) + '0'
	b[6] = byte((month)%10) + '0'
	b[7] = '-'
	b[8] = byte((day)/10) + '0'
	b[9] = byte((day)%10) + '0'
	b[10] = 'T'
	b[11] = byte((hour)/10) + '0'
	b[12] = byte((hour)%10) + '0'
	b[13] = ':'
	b[14] = byte((minute)/10) + '0'
	b[15] = byte((minute)%10) + '0'
	b[16] = ':'
	b[17] = byte((second)/10) + '0'
	b[18] = byte((second)%10) + '0'
	b[19] = '.'
	b[20] = byte((msec/100)%10) + '0'
	b[21] = byte((msec/10)%10) + '0'
	b[22] = byte((msec)%10) + '0'
	b[23] = 'Z'

	enc.AppendString(*(*string)(unsafe.Pointer(&b)))
}
