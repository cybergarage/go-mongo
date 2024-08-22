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
	messageLength int32  // total message size, including this
	requestID     int32  // identifier for this message
	responseTo    int32  // requestID from the original request
	opCode        OpCode // request type
}

// NewHeader returns a new header instance.
func NewHeader() *Header {
	header := &Header{
		messageLength: 0,
		requestID:     0,
		responseTo:    0,
		opCode:        0,
	}
	return header
}

// NewHeaderWithOpCode returns a new header instance with the specified opcode.
func NewHeaderWithOpCode(opcode OpCode) *Header {
	header := &Header{
		messageLength: 0,
		requestID:     0,
		responseTo:    0,
		opCode:        opcode,
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
	header.messageLength = l
}

// MessageLength gets the message length.
func (header *Header) MessageLength() int32 {
	return header.messageLength
}

// SetRequestID sets a message identifier.
func (header *Header) SetRequestID(id int32) {
	header.requestID = id
}

// RequestID gets the message identifier.
func (header *Header) RequestID() int32 {
	return header.requestID
}

// SetResponseTo sets a response message identifier.
func (header *Header) SetResponseTo(id int32) {
	header.responseTo = id
}

// ResponseTo gets the response message identifier.
func (header *Header) ResponseTo() int32 {
	return header.responseTo
}

// OpCode returns a response identifier.
func (header *Header) OpCode() OpCode {
	return header.opCode
}

// ParseBytes parses the specified bytes.
func (header *Header) ParseBytes(msg []byte) error {
	if len(msg) < HeaderSize {
		return fmt.Errorf(errorInvalidMessageHeader, hex.EncodeToString(msg))
	}

	var ok bool
	header.messageLength, header.requestID, header.responseTo, header.opCode, _, ok = ReadHeader(msg)
	if !ok {
		return fmt.Errorf(errorInvalidMessageHeader, hex.EncodeToString(msg[:4]))
	}

	return nil
}

// BodySize returns body size excluding header.
func (header *Header) BodySize() int32 {
	return header.messageLength - HeaderSize
}

// Bytes returns the binary description.
func (header *Header) Bytes() []byte {
	dst := make([]byte, 0)
	dst = AppendHeader(dst, header.messageLength, header.requestID, header.responseTo, header.opCode)
	return dst
}

// String returns the string description.
func (header *Header) String() string {
	return fmt.Sprintf("%d",
		header.opCode)
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
