package cmd

import (
	"github.com/spf13/cobra"

	"github.com/clastix/capsule-addon-rancher/cmd/webhook"
)

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "capsule-addon-rancher",
	Short: "Allowing Rancher shell and cluster-agent to interact with the Capsule Proxy",
}

func Execute() error {
	rootCmd.AddCommand(webhook.New())

	return rootCmd.Execute()
}
