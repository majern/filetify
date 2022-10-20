package cmd

import (
	"github.com/msoft-dev/filetify/pkg/client"
	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Runs Filetify in Client mode.",
	Long: `Use Client mode on your personal computer, this will connect to the server
and synchronize all files from paths specified in the configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		client.StartClient()
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
