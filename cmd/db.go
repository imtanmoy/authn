package cmd

import (
	"github.com/spf13/cobra"

	"github.com/imtanmoy/authy/db"
	"github.com/imtanmoy/authy/logger"
)

func init() {
	rootCmd.AddCommand(dbCmd)
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "database command",
	Run: func(cmd *cobra.Command, args []string) {
		err := db.InitDB()
		if err != nil {
			logger.Fatalf("%s : %s", "Database Could not be initiated", err)
		}
		logger.Info("Database Initiated...")
	},
}
