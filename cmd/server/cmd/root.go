package cmd

import (
	"go-server-base/server"

	"github.com/spf13/cobra"
)

func init() {}

var RootCmd = &cobra.Command{
	Use:   "Server",
	Short: "go-server-base",
	RunE: func(cmd *cobra.Command, args []string) error {
		server.Start()
		return nil
	},
}
