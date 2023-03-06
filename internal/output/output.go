// Copyright 2020-2021 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package output

import (
	"fmt"
	"os"
)

// ExitOnError prints the error and exits.
func ExitOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
