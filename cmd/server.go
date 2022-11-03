package cmd

import (
	"github.com/msoft-dev/filetify/pkg/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs Filetify in Server mode.",
	Long:  `Use this command to run Filetify as a server. In this mode, Filetify will syncrhonize files with clients and store them in the specified directory on server.`,
	Run: func(cmd *cobra.Command, args []string) {
		server.StartServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
