// Copyright 2020 The go-language-server Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package inspector_zap

import (
	"encoding/json"
	"fmt"
	"time"
	"unsafe"

	"github.com/go-language-server/protocol"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/go-language-server/inspector"
)

// logger represents an implemented Language Server Protocol Inspector zap logger.
type logger struct {
	l          *zap.Logger
	traceLevel protocol.TraceMode
}

// NewLogger returns a new Logger which implemented
// Language Server Protocol Inspector specification log format.
func NewLogger(traceLevel protocol.TraceMode, opts ...zap.Option) inspector.Logger {
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

	l, err := cfg.Build(opts...)
	if err != nil {
		panic(fmt.Errorf("inspector_zap.NewLogger: %v", err))
	}

	return &logger{
		l:          l,
		traceLevel: traceLevel,
	}
}

// Client-initiated request message
// interface Request extends Message {
// 	type: "request";
// 	/**
// 	 * The command to execute
// 	 */
// 	command: string;
// 	/**
// 	 * Object containing arguments for the command
// 	 */
// 	arguments?: any;
// }

// Response by server to client request message.
// interface Response extends Message {
// 	/**
// 	 * Sequence number of the message
// 	 */
// 	seq: number;
//
// 	type: "response";
// 	/**
// 	 * Sequence number of the request message.
// 	 */
// 	request_seq: number;
// 	/**
// 	 * Outcome of the request.
// 	 */
// 	success: boolean;
// 	/**
// 	 * The command requested.
// 	 */
// 	command: string;
// 	/**
// 	 * If success === false, this should always be provided.
// 	 * Otherwise, may (or may not) contain a success message.
// 	 */
// 	message?: string;
// 	/**
// 	 * Contains message body if success === true.
// 	 */
// 	body?: any;
// 	/**
// 	 * Contains extra information that plugin can include to be passed on
// 	 */
// 	metadata?: unknown;
// }

// public logTrace(serverId: string, message: string, data?: any): void {
// 	if (this.trace !== Trace.Off) {
// 		this.logger.logLevel('Trace', `<${serverId}> ${message}`, data);
// 	}
// }
func (l *logger) logTrace(serverID string, message string, data interface{}) {
	l.l.Info("Trace", zap.String("<serverID>", serverID), zap.String("message", message), zap.Any("data", data))
}

// public traceRequest(serverId: string, request: Proto.Request, responseExpected: boolean, queueLength: number): void {
// 	if (this.trace === Trace.Off) {
// 		return;
// 	}
// 	let data: string | undefined = undefined;
// 	if (this.trace === Trace.Verbose && request.arguments) {
// 		data = `Arguments: ${JSON.stringify(request.arguments, null, 4)}`;
// 	}
// 	this.logTrace(serverId, `Sending request: ${request.command} (${request.seq}). Response expected: ${responseExpected ? 'yes' : 'no'}. Current queue length: ${queueLength}`, data);
// }
func (l *logger) TraceRequest(serverID string, req *inspector.Payload, expected bool, queueLength int) error {
	if l.traceLevel == protocol.TraceOff {
		return nil
	}

	var data []byte
	if l.traceLevel == protocol.TraceVerbose && req.Arg != nil {
		data, err := json.MarshalIndent(req.Arg, "", "    ")
		if err != nil {
			return err
		}
		data = append([]byte("Arguments:\n"), data...)
	}

	l.logTrace(serverID, fmt.Sprintf("Sending request: %s (%s). Response expected: %t. Current queue length: %d", req.MsgKind, req.MsgLatency, expected, queueLength), data)

	return nil
}

// public traceResponse(serverId: string, response: Proto.Response, meta: RequestExecutionMetadata): void {
// 	if (this.trace === Trace.Off) {
// 		return;
// 	}
// 	let data: string | undefined = undefined;
// 	if (this.trace === Trace.Verbose && response.body) {
// 		data = `Result: ${JSON.stringify(response.body, null, 4)}`;
// 	}
// 	this.logTrace(serverId, `Response received: ${response.command} (${response.request_seq}). Request took ${Date.now() - meta.queuingStartTime} ms. Success: ${response.success} ${!response.success ? '. Message: ' + response.message : ''}`, data);
// }
func (l *logger) TraceResponse(serverID string, req *inspector.Payload, meta inspector.Metadata) error {
	if l.traceLevel == protocol.TraceOff {
		return nil
	}

	var data []byte
	if l.traceLevel == protocol.TraceVerbose && req.Arg != nil {
		data, err := json.MarshalIndent(req.Arg, "", "    ")
		if err != nil {
			return err
		}
		data = append([]byte("Result:\n"), data...)
	}

	took := time.Now().Sub(meta.QueuingStartTime)
	msg := ". Message: "
	l.logTrace(serverID, fmt.Sprintf("Response received: %s (%s). Request took %d ms. Success: %s", req.MsgKind, req.MsgLatency, took, msg), data)

	return nil
}

// public traceRequestCompleted(serverId: string, command: string, request_seq: number, meta: RequestExecutionMetadata): any {
// 	if (this.trace === Trace.Off) {
// 		return;
// 	}
// 	this.logTrace(serverId, `Async response received: ${command} (${request_seq}). Request took ${Date.now() - meta.queuingStartTime} ms.`);
// }
func (l *logger) TraceRequestCompleted(serverID string, command string, requestSeq int, meta inspector.Metadata) error {
	if l.traceLevel == protocol.TraceOff {
		return nil
	}

	took := time.Now().Sub(meta.QueuingStartTime)
	l.logTrace(serverID, fmt.Sprintf("Async response received: %s (%d). Request took %s ms.", command, requestSeq, took), nil)

	return nil
}

// public traceEvent(serverId: string, event: Proto.Event): void {
// 	if (this.trace === Trace.Off) {
// 		return;
// 	}
// 	let data: string | undefined = undefined;
// 	if (this.trace === Trace.Verbose && event.body) {
// 		data = `Data: ${JSON.stringify(event.body, null, 4)}`;
// 	}
// 	this.logTrace(serverId, `Event received: ${event.event} (${event.seq}).`, data);
// }
func (l *logger) TraceEvent(serverID string, event *inspector.Payload) error {
	if l.traceLevel == protocol.TraceOff {
		return nil
	}

	var data []byte
	if l.traceLevel == protocol.TraceVerbose && event.Arg != nil {
		data, err := json.MarshalIndent(event.Arg, "", "    ")
		if err != nil {
			return err
		}
		data = append([]byte("Data:\n"), data...)
	}

	l.logTrace(serverID, fmt.Sprintf("Event received: %s (%s).", event.Msg, event.MsgLatency), data)

	return nil
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
