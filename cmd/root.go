package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/imtanmoy/authn/config"
)

func init() {
	cobra.OnInitialize(config.InitConfig)
}

var rootCmd = &cobra.Command{
	Use:   "Root",
	Short: "authn",
	Long:  "authn service for authentication and identity",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
