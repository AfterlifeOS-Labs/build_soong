// Copyright 2024 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package android

import (
	"testing"
)

var testCases []struct {
	name          string
	fs            MockFS
	expectedError string
} = []struct {
	name          string
	fs            MockFS
	expectedError string
}{
	{
		name: "Can't reference variable before assignment",
		fs: map[string][]byte{
			"Android.bp": []byte(`
x = foo
foo = "hello"
`),
		},
		expectedError: "undefined variable foo",
	},
	{
		name: "Can't append to variable before assigned to",
		fs: map[string][]byte{
			"Android.bp": []byte(`
foo += "world"
foo = "hello"
`),
		},
		expectedError: "modified non-existent variable \"foo\" with \\+=",
	},
	{
		name: "Can't reassign variable",
		fs: map[string][]byte{
			"Android.bp": []byte(`
foo = "hello"
foo = "world"
`),
		},
		expectedError: "variable already set, previous assignment:",
	},
	{
		name: "Can't reassign variable in inherited scope",
		fs: map[string][]byte{
			"Android.bp": []byte(`
foo = "hello"
`),
			"foo/Android.bp": []byte(`
foo = "world"
`),
		},
		expectedError: "variable already set in inherited scope, previous assignment:",
	},
	{
		name: "Can't modify variable in inherited scope",
		fs: map[string][]byte{
			"Android.bp": []byte(`
foo = "hello"
`),
			"foo/Android.bp": []byte(`
foo += "world"
`),
		},
		expectedError: "modified non-local variable \"foo\" with \\+=",
	},
	{
		name: "Can't modify variable after referencing",
		fs: map[string][]byte{
			"Android.bp": []byte(`
foo = "hello"
x = foo
foo += "world"
`),
		},
		expectedError: "modified variable \"foo\" with \\+= after referencing",
	},
}

func TestBlueprintErrors(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fixtures := FixtureMergeMockFs(tc.fs)
			fixtures = fixtures.ExtendWithErrorHandler(FixtureExpectsOneErrorPattern(tc.expectedError))
			fixtures.RunTest(t)
		})
	}
}
