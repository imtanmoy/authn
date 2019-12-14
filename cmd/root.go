package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/imtanmoy/authy/config"
	"github.com/imtanmoy/authy/logger"
)

func init() {
	cobra.OnInitialize(config.InitConfig)
	err := logger.InitLogger()
	if err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "Root",
	Short: "Authy",
	Long:  "Authy service for authentication and identity",
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