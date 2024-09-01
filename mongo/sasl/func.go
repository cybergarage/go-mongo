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

package sasl

import (
	"crypto/md5"
	"fmt"
	"io"
)

// Authentication - MongoDB Specifications
// https://github.com/mongodb/specifications/blob/master/source/auth/auth.md

// SCRAM-SHA-1
// Since: 3.0
// SCRAM-SHA-1 is defined in RFC 5802.
// Page 11 of the RFC specifies that user names be prepared with SASLprep, but drivers MUST NOT do so.
// Page 8 of the RFC identifies the "SaltedPassword" as := Hi(Normalize(password), salt, i). The password variable MUST be the mongodb hashed variant.
// The mongo hashed variant is computed as hash = HEX( MD5( UTF8( username + ':mongo:' + plain_text_password ))), where plain_text_password is actually plain text.
// The username and password MUST NOT be prepared with SASLprep before hashing.
func MongoPasswordDigest(username, password string) (string, error) {
	h := md5.New()
	if _, err := io.WriteString(h, username); err != nil {
		return "", err
	}
	if _, err := io.WriteString(h, ":mongo:"); err != nil {
		return "", err
	}
	if _, err := io.WriteString(h, password); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
