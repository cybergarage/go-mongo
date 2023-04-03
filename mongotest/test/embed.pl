#!/usr/bin/perl
# Copyright (C) 2020 Snm                                                     atoshi Konno. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


use strict;
use warnings;
use File::Find;

my $test_dir = "./";

if (1 <= @ARGV){
  $test_dir = $ARGV[0];
}

my @test_files = ();

find(
  sub {
    if (-f && /\.qst$/) {
      push(@test_files, $File::Find::name);
    }
  },
, $test_dir);

my @embed_test_files = ();
my @embed_test_names = ();

foreach my $test_file(@test_files){
  my @test_file_paths = split(/\//, $test_file);
  my $test_file_name = $test_file_paths[-1];
  push(@embed_test_files, $test_file_name);
  my @test_file_names = split(/\./, $test_file_name);
  my $snake_case_test_name = $test_file_names[0];
  my $camel_case_test_name = join('', map{ucfirst($_)} split(/_/, $snake_case_test_name));
  push(@embed_test_names, $camel_case_test_name);
}

print<<HEADER;
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

package test

import (
	_ "embed"
)

var EmbedTests = map[string][]byte {
HEADER
foreach my $name(@embed_test_names){
  printf("\t\"%s\": %s,\n", $name, lcfirst($name));
}
print<<FOTTER;
}

FOTTER
print<<HEADER;
HEADER
my $n;
foreach my $name(@embed_test_names){
  printf("//go:embed %s\n", $embed_test_files[$n++]);
  printf("var %s []byte\n\n", lcfirst($name));
}
