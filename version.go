// Copyright (C) 2024  Morgan Stewart Hein
//
// This Source Code Form is subject to the terms
// of the Mozilla Public License, v. 2.0. If a copy
// of the MPL was not distributed with this file, You
// can obtain one at https://mozilla.org/MPL/2.0/.

package gogo

// These attributes should only be changed by the build script, do not change it manually.
var (
	BuildSha   string = "..."
	BuildTime  string = "..."
	Who        string = "..."
	State      string = "..."
	VersionTag        = "0.1.0"
)

func Version() version {
	return version{
		BuildSha:  BuildSha,
		BuildTime: BuildTime,
		Who:       Who,
		State:     State,
		Version:   VersionTag,
	}
}

type version struct {
	BuildSha  string
	BuildTime string
	Who       string
	State     string
	Version   string
}
