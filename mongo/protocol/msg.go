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
	"fmt"

	"github.com/cybergarage/go-mongo/mongo/bson"

	"go.mongodb.org/mongo-driver/x/network/wiremessage"
)

const (
	checksumPresent             = MsgFlag(0x01)
	moreToCome                  = MsgFlag(0x02)
	exhaustAllowed              = MsgFlag(0xF0)
	sectionTypeBody             = SectionType(0)
	sectionTypeDocumentSequence = SectionType(1)
)

// MsgFlag represents a MsgFlag of MongoDB wire protocol.
type MsgFlag = wiremessage.MsgFlag

// SectionType represents a SectionType of MongoDB wire protocol.
type SectionType = wiremessage.SectionType

// Msg represents a OP_MSG of MongoDB wire protocol.
// See : MongoDB Wire Protocol
// https://docs.mongodb.com/manual/reference/mongodb-wire-protocol/
type Msg struct {
	*Header                   // A standard wire protocol header
	FlagBits    MsgFlag       // message flags
	Body        bson.Document // Body
	DocumentIDs []string
	Documents   []bson.Document // Document Sequence
	Checksum    uint32          // optional CRC-32C checksum
}

// NewMsg returns a new msg instance.
func NewMsg() *Msg {
	op := &Msg{
		Header:      NewHeaderWithOpCode(OpMsg),
		FlagBits:    0,
		Body:        nil,
		DocumentIDs: make([]string, 0),
		Documents:   make([]bson.Document, 0),
		Checksum:    0,
	}
	op.SetMessageLength(op.Size())
	return op
}

// NewMsgWithBody returns a new msg instance with the specified body document.
func NewMsgWithBody(body bson.Document) *Msg {
	op := NewMsg()
	op.Body = body
	op.SetMessageLength(op.Size())
	return op
}

// NewMsgWithHeaderAndBody returns a new update instance with the specified bytes.
func NewMsgWithHeaderAndBody(header *Header, body []byte) (*Msg, error) {
	flagBits, offsetBody, ok := ReadUint32(body)
	if !ok {
		return nil, newMessageRequestError(OpMsg, body)
	}

	var docBody bson.Document
	documentIDs := make([]string, 0)
	documents := make([]bson.Document, 0)

	var sectionType SectionType
	offsetBodyLen := len(offsetBody)
	for 0 < offsetBodyLen {
		sectionType, offsetBody, ok = ReadSectionType(offsetBody)
		if !ok {
			break
		}
		switch sectionType {
		case sectionTypeBody:
			docBody, offsetBody, ok = ReadDocument(offsetBody)
			if !ok {
				return nil, newMessageRequestError(OpMsg, body)
			}
		case sectionTypeDocumentSequence:
			var docID string
			var docs []bson.Document
			docID, docs, offsetBody, ok = ReadDocumentSequence(offsetBody)
			if !ok {
				return nil, newMessageRequestError(OpMsg, body)
			}
			documentIDs = append(documentIDs, docID)
			documents = append(documents, docs...)
		}
		offsetBodyLen = len(offsetBody)
	}

	checksum := uint32(0)
	if (MsgFlag(flagBits) & checksumPresent) != 0 {
		checksum, offsetBody, ok = ReadUint32(offsetBody)
		if !ok {
			return nil, newMessageRequestError(OpMsg, body)
		}
		// TODO : Check the CRC-32C checksum
	}

	op := &Msg{
		Header:      header,
		FlagBits:    MsgFlag(flagBits),
		Body:        docBody,
		DocumentIDs: documentIDs,
		Documents:   documents,
		Checksum:    checksum,
	}

	return op, nil
}

// SetBody sets the specified body.
func (op *Msg) SetBody(doc bson.Document) {
	op.Body = doc
	op.SetMessageLength(op.Size())
}

// AddDocument adds a document with the ID.
func (op *Msg) AddDocument(docID string, doc bson.Document) {
	op.DocumentIDs = append(op.DocumentIDs, docID)
	op.Documents = append(op.Documents, doc)
}

// AddDocuments adds documents with the IDs.
func (op *Msg) AddDocuments(docIDs []string, docs []bson.Document) {
	op.DocumentIDs = append(op.DocumentIDs, docIDs...)
	op.Documents = append(op.Documents, docs...)
}

// GetBody returns the body document.
func (op *Msg) GetBody() bson.Document {
	return op.Body
}

// GetDocuments returns the sequence documents.
func (op *Msg) GetDocuments() []bson.Document {
	return op.Documents
}

// Size returns the message size including the header.
func (op *Msg) Size() int32 {
	bodySize := 4
	if op.Body != nil {
		bodySize += 1 + len(op.Body)
	}
	for n, doc := range op.Documents {
		bodySize += 1 + (len(op.DocumentIDs[n]) + 1) + len(doc)
	}
	if (op.FlagBits & checksumPresent) != 0 {
		bodySize += 4
	}
	return int32(HeaderSize + bodySize)
}

// Bytes returns the binary description of BSON format.
func (op *Msg) Bytes() []byte {
	dst := op.Header.Bytes()
	dst = AppendInt32(dst, int32(op.FlagBits))
	if op.Body != nil {
		dst = AppendByte(dst, byte(sectionTypeBody))
		dst = AppendDocument(dst, op.Body)
	}
	for n, doc := range op.Documents {
		dst = AppendByte(dst, byte(sectionTypeDocumentSequence))
		dst = AppendCString(dst, op.DocumentIDs[n])
		dst = AppendDocument(dst, doc)
	}
	if (op.FlagBits & checksumPresent) != 0 {
		// TODO : Add a CRC-32C checksum
	}
	return dst
}

// String returns the string description.
func (op *Msg) String() string {
	str := fmt.Sprintf("%s %X ",
		op.Header.String(),
		op.FlagBits,
	)

	if op.Body != nil {
		str += fmt.Sprintf("%s ", op.Body.String())
	}
	for n, doc := range op.Documents {
		str += fmt.Sprintf("%s : %s ", op.DocumentIDs[n], doc.String())
	}

	if (op.FlagBits & checksumPresent) != 0 {
		str += fmt.Sprintf("%X", op.Checksum)
	}

	return str
}
