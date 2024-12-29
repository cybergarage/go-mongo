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
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
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
	body        bson.Document // Body
	documentIDs []string
	documents   []bson.Document // Document Sequence
	Checksum    uint32          // optional CRC-32C checksum
}

// NewMsg returns a new msg instance.
func NewMsg() *Msg {
	op := &Msg{
		Header:      NewHeaderWithOpCode(OpMsg),
		FlagBits:    0,
		body:        nil,
		documentIDs: make([]string, 0),
		documents:   make([]bson.Document, 0),
		Checksum:    0,
	}
	op.SetMessageLength(op.Size())
	return op
}

// NewMsgWithBody returns a new msg instance with the specified body document.
func NewMsgWithBody(body bson.Document) *Msg {
	op := NewMsg()
	op.body = body
	op.SetMessageLength(op.Size())
	return op
}

// NewMsgWithHeaderAndBody returns a new update instance with the specified bytes.
func NewMsgWithHeaderAndBody(header *Header, body []byte) (*Msg, error) {
	flagBits, offsetBody, ok := ReadUint32(body)
	if !ok {
		return nil, newErrMessageRequest(OpMsg, body)
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
				return nil, newErrMessageRequest(OpMsg, body)
			}
		case sectionTypeDocumentSequence:
			var docID string
			var docs []bson.Document
			docID, docs, offsetBody, ok = ReadDocumentSequence(offsetBody)
			if !ok {
				return nil, newErrMessageRequest(OpMsg, body)
			}
			documentIDs = append(documentIDs, docID)
			documents = append(documents, docs...)
		}
		offsetBodyLen = len(offsetBody)
	}

	checksum := uint32(0)
	if (MsgFlag(flagBits) & checksumPresent) != 0 {
		checksum, _, ok = ReadUint32(offsetBody)
		if !ok {
			return nil, newErrMessageRequest(OpMsg, body)
		}
		// TODO : Check the CRC-32C checksum
	}

	op := &Msg{
		Header:      header,
		FlagBits:    MsgFlag(flagBits),
		body:        docBody,
		documentIDs: documentIDs,
		documents:   documents,
		Checksum:    checksum,
	}

	return op, nil
}

// SetBody sets the specified body.
func (op *Msg) SetBody(doc bson.Document) {
	op.body = doc
	op.SetMessageLength(op.Size())
}

// AddDocument adds a document with the ID.
func (op *Msg) AddDocument(docID string, doc bson.Document) {
	op.documentIDs = append(op.documentIDs, docID)
	op.documents = append(op.documents, doc)
}

// AddDocuments adds documents with the IDs.
func (op *Msg) AddDocuments(docIDs []string, docs []bson.Document) {
	op.documentIDs = append(op.documentIDs, docIDs...)
	op.documents = append(op.documents, docs...)
}

// Body returns the body document.
func (op *Msg) Body() bson.Document {
	return op.body
}

// Documents returns the sequence documents.
func (op *Msg) Documents() []bson.Document {
	return op.documents
}

// Size returns the message size including the header.
func (op *Msg) Size() int32 {
	bodySize := 4
	if op.body != nil {
		bodySize += 1 + len(op.body)
	}
	for n, doc := range op.documents {
		bodySize += 1 + (len(op.documentIDs[n]) + 1) + len(doc)
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
	if op.body != nil {
		dst = AppendByte(dst, byte(sectionTypeBody))
		dst = AppendDocument(dst, op.body)
	}
	for n, doc := range op.documents {
		dst = AppendByte(dst, byte(sectionTypeDocumentSequence))
		dst = AppendCString(dst, op.documentIDs[n])
		dst = AppendDocument(dst, doc)
	}
	// if (op.FlagBits & checksumPresent) != 0 {
	// 	// TODO : Add a CRC-32C checksum.
	// }
	return dst
}

// String returns the string description.
func (op *Msg) String() string {
	str := fmt.Sprintf("%s %X ",
		op.Header.String(),
		op.FlagBits,
	)

	if op.body != nil {
		str += fmt.Sprintf("%s ", op.body.String())
	}
	for n, doc := range op.documents {
		str += fmt.Sprintf("%s : %s ", op.documentIDs[n], doc.String())
	}

	if (op.FlagBits & checksumPresent) != 0 {
		str += fmt.Sprintf("%X", op.Checksum)
	}

	return str
}
