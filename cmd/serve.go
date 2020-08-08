package cmd

import (
	"github.com/imtanmoy/authn/config"
	"github.com/imtanmoy/authn/registry"
	"github.com/spf13/cobra"

	"github.com/imtanmoy/logx"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server with configured api",
	Run: func(cmd *cobra.Command, args []string) {

		r := registry.NewRegistry(config.Conf)
		err := r.Init()
		if err != nil {
			logx.Fatalf("%s : %s", "could not init registry", err)
		}

		ServeAll(r)

		r.Close()
	},
}
