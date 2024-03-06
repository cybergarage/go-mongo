// Copyright (C) 2019 The go-mongo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protocol

import (
	"encoding/hex"
	"fmt"
)

const (
	// HeaderSize is the static header size of MongoDB wire protocol.
	HeaderSize = (4 + 4 + 4 + 4)
)

// Header represents a standard header of MongoDB wire protocol.
type Header struct {
	MessageLength int32  // total message size, including this
	RequestID     int32  // identifier for this message
	ResponseTo    int32  // requestID from the original request
	OpCode        OpCode // request type
}

// NewHeader returns a new header instance.
func NewHeader() *Header {
	header := &Header{
		MessageLength: 0,
		RequestID:     0,
		ResponseTo:    0,
		OpCode:        0,
	}
	return header
}

// NewHeaderWithOpCode returns a new header instance with the specified opcode.
func NewHeaderWithOpCode(opcode OpCode) *Header {
	header := &Header{
		MessageLength: 0,
		RequestID:     0,
		ResponseTo:    0,
		OpCode:        opcode,
	}
	return header
}

// NewHeaderWithBytes returns a new header instance of the specified bytes.
func NewHeaderWithBytes(msg []byte) (*Header, error) {
	header := NewHeader()
	return header, header.ParseBytes(msg)
}

// SetMessageLength sets a message length.
func (header *Header) SetMessageLength(l int32) {
	header.MessageLength = l
}

// GetMessageLength gets the message length.
func (header *Header) GetMessageLength() int32 {
	return header.MessageLength
}

// SetRequestID sets a message identifier.
func (header *Header) SetRequestID(id int32) {
	header.RequestID = id
}

// GetRequestID gets the message identifier.
func (header *Header) GetRequestID() int32 {
	return header.RequestID
}

// SetResponseTo sets a response message identifier.
func (header *Header) SetResponseTo(id int32) {
	header.ResponseTo = id
}

// GetResponseTo gets the response message identifier.
func (header *Header) GetResponseTo() int32 {
	return header.ResponseTo
}

// GetOpCode returns a response identifier.
func (header *Header) GetOpCode() OpCode {
	return header.OpCode
}

// ParseBytes parses the specified bytes.
func (header *Header) ParseBytes(msg []byte) error {
	if len(msg) < HeaderSize {
		return fmt.Errorf(errorInvalidMessageHeader, hex.EncodeToString(msg))
	}

	var ok bool
	header.MessageLength, header.RequestID, header.ResponseTo, header.OpCode, _, ok = ReadHeader(msg)
	if !ok {
		return fmt.Errorf(errorInvalidMessageHeader, hex.EncodeToString(msg[:4]))
	}

	return nil
}

// GetBodySize returns body size excluding header.
func (header *Header) GetBodySize() int32 {
	return header.MessageLength - HeaderSize
}

// Bytes returns the binary description.
func (header *Header) Bytes() []byte {
	dst := make([]byte, 0)
	dst = AppendHeader(dst, header.MessageLength, header.RequestID, header.ResponseTo, header.OpCode)
	return dst
}

// String returns the string description.
func (header *Header) String() string {
	return fmt.Sprintf("%d",
		header.OpCode)
}

// ReadHeader reads a wire message header from src.
func ReadHeader(src []byte) (int32, int32, int32, OpCode, []byte, bool) {
	if len(src) < 16 {
		return 0, 0, 0, 0, src, false
	}
	length := (int32(src[0]) | int32(src[1])<<8 | int32(src[2])<<16 | int32(src[3])<<24)
	requestID := (int32(src[4]) | int32(src[5])<<8 | int32(src[6])<<16 | int32(src[7])<<24)
	responseTo := (int32(src[8]) | int32(src[9])<<8 | int32(src[10])<<16 | int32(src[11])<<24)
	opcode := OpCode(int32(src[12]) | int32(src[13])<<8 | int32(src[14])<<16 | int32(src[15])<<24)
	return length, requestID, responseTo, opcode, src[16:], true
}

// AppendHeader appends a header to dst.
func AppendHeader(dst []byte, length, reqid, respto int32, opcode OpCode) []byte {
	dst = AppendInt32(dst, length)
	dst = AppendInt32(dst, reqid)
	dst = AppendInt32(dst, respto)
	dst = AppendInt32(dst, int32(opcode))
	return dst
}
