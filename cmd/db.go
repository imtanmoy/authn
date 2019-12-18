package cmd

import (
	"github.com/spf13/cobra"

	"github.com/imtanmoy/authy/db"
	"github.com/imtanmoy/logx"
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
			logx.Fatalf("%s : %s", "Database Could not be initiated", err)
		}
		logx.Info("Database Initiated...")
	},
}
