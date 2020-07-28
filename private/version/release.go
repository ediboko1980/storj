// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package version

import _ "unsafe" // needed for go:linkname

//go:linkname buildTimestamp storj.io/private/version.buildTimestamp
var buildTimestamp string = "1595928394"

//go:linkname buildCommitHash storj.io/private/version.buildCommitHash
var buildCommitHash string = "a7bfee5f2c3bf76136c4c1c260af89b91b826f02"

//go:linkname buildVersion storj.io/private/version.buildVersion
var buildVersion string = "v1.9.3"

//go:linkname buildRelease storj.io/private/version.buildRelease
var buildRelease string = "true"

// ensure that linter understands that the variables are being used.
func init() { use(buildTimestamp, buildCommitHash, buildVersion, buildRelease) }

func use(...interface{}) {}
