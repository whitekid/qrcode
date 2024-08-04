package config

import (
	"github.com/spf13/viper"
)

const (
	KeyBindAddr  = "bindAddr"
	KeyRateLimit = "rateLimit"

	keyGrpcBind = "grpc.bindAddr"
)

func init() {
	viper.SetEnvPrefix("qr")
	viper.AutomaticEnv()
}

func BindAddr() string { return viper.GetString(KeyBindAddr) }
func RateLimit() int   { return viper.GetInt(KeyRateLimit) }
