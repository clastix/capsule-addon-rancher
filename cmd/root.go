// Copyright 2020-2021 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/clastix/capsule-addon-rancher/cmd/webhook"
	"github.com/clastix/capsule-addon-rancher/internal/output"
)

func New() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "capsule-addon-rancher",
		Short: "Allowing Rancher shell and cluster-agent to interact with the Capsule Proxy",
	}
	rootCmd.AddCommand(webhook.New())

	return rootCmd
}

func Execute() {
	rootCmd := New()
	output.ExitOnError(rootCmd.Execute())
}
