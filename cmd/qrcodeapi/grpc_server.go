package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/whitekid/goxp/flags"
	"github.com/whitekid/goxp/log"

	"qrcodeapi/apiserver/grpcserver"
	"qrcodeapi/config"
)

func init() {
	cmd := &cobra.Command{
		Use: "grpc-server",
		RunE: func(cmd *cobra.Command, args []string) error {
			bindAddr := viper.GetString("bindAddr")

			if err := grpcserver.Run(cmd.Context(), bindAddr); err != nil {
				log.Errorf("%+v", err)
				return err
			}
			return nil
		},
	}

	fs := cmd.PersistentFlags()
	flags.String(fs, config.KeyBindAddr, "bind_addr", "B", "127.0.0.1:9000", "bind addr")

	rootCmd.AddCommand(cmd)
}
