package main

import (
	"github.com/spf13/cobra"
	"github.com/whitekid/cobrax"
	"github.com/whitekid/goxp/flags"
	"github.com/whitekid/goxp/log"

	"qrcodeapi/apiserver/grpcserver"
)

func init() {
	cobrax.Add(rootCmd, &cobra.Command{
		Use: "grpc-server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO viper와 공존하려면?...
			bindAddr := cobrax.Apply(cmd.Flags().GetString, "bind_addr")

			if err := grpcserver.Run(cmd.Context(), bindAddr); err != nil {
				log.Errorf("%+v", err)
				return err
			}
			return nil
		},
	}, []flags.Flag{
		{"bind_addr", "B", "127.0.0.1:9000", "bind addr"},
	}, nil)
}
