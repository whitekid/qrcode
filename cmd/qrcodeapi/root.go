package main

import (
	"github.com/spf13/cobra"
	"github.com/whitekid/goxp/log"

	"qrcodeapi"
	"qrcodeapi/config"
)

var rootCmd = &cobra.Command{
	Use: "qrcodeapi",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := qrcodeapi.Run(cmd.Context()); err != nil {
			log.Errorf("%+v", err)
			return err
		}
		return nil
	},
}

func init() {
	config.InitFlagSet(rootCmd.Use, rootCmd.Flags())
}
