// Copyright (C) 2022 Satoshi Konno All rights reserved.
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

package regexp

import (
	"regexp"
	"strings"
)

// MustCompile compiles the specified glob-style patterns.
func MustCompile(pattern string) *regexp.Regexp {
	// Syntax · google/re2 Wiki
	// https://github.com/google/re2/wiki/Syntax
	// glob (programming) - Wikipedia
	// https://en.wikipedia.org/wiki/Glob_(programming)
	repstrs := []struct {
		old string
		new string
	}{
		{old: "*", new: ".*"},
		{old: "?", new: "."},
	}
	re2Pattern := pattern
	for _, repstr := range repstrs {
		re2Pattern = strings.ReplaceAll(re2Pattern, repstr.old, repstr.new)
	}
	return regexp.MustCompile(re2Pattern)
}
