// Copyright (C) 2020 The go-mongo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

const (
	errorFileListNotFound     = "File (%s) not found"
	errorFileListBadExtension = "Invalid Extension (%s) != *.%s"
)

func TestFileListFiles(t *testing.T) {
	file := NewFileWithPath("./")

	// All files

	files, err := file.ListFiles()
	if err != nil {
		t.Error(err)
	}

	if 0 < len(files) {
		for _, file := range files {
			_, err := os.Stat(file.GetPath())
			if os.IsNotExist(err) {
				t.Error(err)
			}
		}
	} else {
		t.Errorf(errorFileListNotFound, "")
	}

	// *.go files

	ext := "go"
	files, err = file.ListFilesWithExtention(ext)
	if err != nil {
		t.Error(err)
	}

	if 0 < len(files) {
		for _, file := range files {
			path := file.Path
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				t.Error(err)
			}
			if !strings.HasSuffix(path, ext) {
				t.Errorf(errorFileListBadExtension, file.Path, ext)
			}
		}
	} else {
		t.Errorf(errorFileListNotFound, ext)
	}

	// *.go files (Regexp)

	// nolint
	// linter says regexpMust: for const patterns like ".*\\.go", use regexp.MustCompile (gocritic)
	// but MustCompile doesn't work.
	re, err := regexp.Compile(".*\\.go")
	if err != nil {
		t.Error(err)
	}

	files, err = file.ListFilesWithRegexp(re)
	if err != nil {
		t.Error(err)
	}

	if 0 < len(files) {
		for _, file := range files {
			path := file.Path
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				t.Error(err)
			}
			if !strings.HasSuffix(path, ext) {
				t.Errorf(errorFileListBadExtension, file.Path, ext)
			}
		}
	} else {
		t.Errorf(errorFileListNotFound, ext)
	}
}
