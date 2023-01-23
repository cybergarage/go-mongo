// Copyright (C) 2017 The go-mongo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mongo

import (
	"errors"
	"fmt"
)

var ErrQuery = errors.New("query error")
var ErrQueryNotSupported = errors.New("query not supported")

const (
	errorLostConnection                    = "lost connection to %s:%d"
	errorCollectionNotFound                = "collection (%s:%s) not found"
	errorMessageHanderNotImplemented       = "MessageHandler does not implemented"
	errorMessageHandeUnknownOpCode         = "MessageHandler does not support OpCode (%d)"
	errorMessageHanderNotSupported         = "MessageHandler does not support (%d)"
	errorQueryHanderNotImplemented         = "QueryHandler does not support (%s)"
	errorOpMsgDocumentSequenceNotSupported = "document Sequence does not supported"
)

func NewQueryError(q *Query) error {
	return fmt.Errorf("%w (%v)", ErrQuery, q)
}

func NewNotSupported(q *Query) error {
	return fmt.Errorf("%w (%v)", ErrQueryNotSupported, q)
}
