package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/imtanmoy/authn/db"
	"github.com/imtanmoy/authn/server"
	"github.com/imtanmoy/logx"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server with configured api",
	Run: func(cmd *cobra.Command, args []string) {
		// initializing database
		err := db.InitDB()
		if err != nil {
			logx.Fatalf("%s : %s", "Database Could not be initiated", err)
		}
		logx.Info("Database Initiated...")

		// initializing server
		server, err := server.NewServer()
		if err != nil {
			logx.Fatalf("%s : %s", "Server could not be started", err)
		}
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGSTOP)

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			oscall := <-c
			logx.Infof("system call:%+v", oscall)
			cancel()
		}()

		if err := server.Run(ctx); err != nil {
			logx.Infof("failed to serve:+%v\n", err)
		}
		close(c)
	},
}
