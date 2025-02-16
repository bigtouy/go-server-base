package cmd

import (
	"go-server-base/server"

	"github.com/spf13/cobra"
)

func init() {}

var RootCmd = &cobra.Command{
	Use:   "Rag",
	Short: "Rag，通义点金",
	RunE: func(cmd *cobra.Command, args []string) error {
		server.Start()
		return nil
	},
}
